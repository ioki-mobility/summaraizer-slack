package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ioki-mobility/summaraizer"
	"github.com/ioki-mobility/summaraizer-slack/slack"
)

const summarizeAiPrompt = `
I give you a discussion and you give me a summary.
Each comment of the discussion is wrapped in a <comment> tag.
Your summary should not be longer than 1200 chars.
Here is the discussion:
{{ range $comment := . }}
<comment>{{ $comment.Body }}</comment>
{{end}}
`

const chatbotAiPrompt = `
You are an advanced chatbot tasked with responding to a user message. Your responses should be informed, contextually aware, and relevant. Below are your instructions:

1. If a history of previous discussions is provided, use it to understand the context of the current message and craft your response accordingly.
2. If no history is provided, focus solely on the current message and generate a standalone response.
3. Use the previous discussion solely as background knowledge to inform your reply. Do not refer to it explicitly.

Example (What NOT to do):
- "Considering the previous discussion, I believe..."
- "Given the context, my response is..."
- "It looks like you're still curious about our conversation..."

Example (What to do):
- Respond directly: "Here is the information you need..." or "Yes, we can proceed with..."

The input will have the following structure:

- Previous Discussion (if available):
<comment>First comment in the history</comment>
<comment>Second comment in the history</comment>
...

- Current Message (to respond to):
<message>
This is the new user message you should respond to.
</message>

**Your task**: Respond to the message thoughtfully and contextually, considering any relevant details in the previous discussion (if present).

Below is the input:

Previous Discussion:
{{ range $comment := . }}
<comment>{{ $comment.Body }}</comment>
{{end}}

Message to Respond To:
`

var slackBotToken = os.Getenv("SLACK_BOT_TOKEN")
var slackSigningSecretEvent = os.Getenv("SLACK_SIGNING_SECRET")

var openAiToken = os.Getenv("OPENAI_API_TOKEN")
var ollamUrl = os.Getenv("OLLAMA_URL")

func EventHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if valid := slack.VerifySignature(r.Header, body, slackSigningSecretEvent); valid != true {
		http.Error(w, "Slack signature doesn't match!", http.StatusBadRequest)
		return
	}

	var event map[string]interface{}
	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	innerEvent, ok := event["event"].(map[string]interface{})
	if !ok {
		http.Error(w, "Invalid challenge", http.StatusBadRequest)
		return
	}

	if innerEvent["type"] == "app_mention" {
		user, _ := innerEvent["user"].(string)
		text, _ := innerEvent["text"].(string)
		text = strings.Join(strings.Fields(text)[1:], " ") // Remove the bot name from the text
		threadTs, _ := innerEvent["thread_ts"].(string)
		ts, _ := innerEvent["ts"].(string)
		channel, _ := innerEvent["channel"].(string)

		if threadTs != "" && strings.HasPrefix(strings.ToLower(text), "summarize please") {
			log.Printf("Summarize request in channel %s, thread %s by %s", channel, threadTs, user)
			messageTs := slack.SendMessage(":brain: Summarizing...", channel, threadTs, slackBotToken)
			summarization := fetchAndSummarize(channel, threadTs)
			slack.UpdateMessage(summarizeResponseMessageTemplate(summarization), channel, messageTs, slackBotToken)
			return
		}

		log.Printf("Chat request in channel %s, ts %s by %s", channel, threadTs, user)
		messageTs := slack.SendMessage(":brain: Responding...", channel, ts, slackBotToken)
		chatbotResponse := fetchAndResponse(channel, threadTs, text)
		slack.UpdateMessage(chatbotResponse, channel, messageTs, slackBotToken)
	}
}

func fetchAndSummarize(channel, threadTs string) string {
	threadDiscussion := fetchSlackThread(channel, threadTs)
	return summarize(summarizeAiPrompt, threadDiscussion)
}

func fetchAndResponse(channel, threadTs, message string) string {
	var threadDiscussion = `[{ "author": "", "body": "" }]`
	if threadTs != "" {
		threadDiscussion = fetchSlackThread(channel, threadTs)
	}
	prompt := fmt.Sprintf("%s <message>%s</message>", chatbotAiPrompt, message)
	return summarize(prompt, threadDiscussion)
}

func fetchSlackThread(channel, threadTs string) string {
	buffer := bytes.Buffer{}
	source := summaraizer.Slack{
		Token:   slackBotToken,
		Channel: channel,
		TS:      threadTs,
	}
	source.Fetch(&buffer)
	return buffer.String()
}

func summarize(prompt, threadDiscussion string) string {
	var summarization string
	var err error
	var summarizer summaraizer.Summarizer
	switch {
	case openAiToken != "":
		summarizer = &summaraizer.OpenAi{
			Model:    "gpt-4o-mini",
			Prompt:   prompt,
			ApiToken: openAiToken,
		}
		break
	case ollamUrl != "":
		summarizer = &summaraizer.Ollama{
			Model:  "llama3.1:latest",
			Prompt: prompt,
			Url:    ollamUrl,
		}
		break
	}
	if summarizer == nil {
		log.Fatal("OpenAiToken AND OllamaUrl are nil. Please set one of them.")
	}

	summarization, err = summarizer.Summarize(strings.NewReader(threadDiscussion))

	if err != nil {
		log.Fatal(err)
	}

	return summarization
}

func summarizeResponseMessageTemplate(message string) string {
	summaraizerLink := "<https://github.com/ioki-mobility/summaraizer|summaraizer>"
	summaraizerSlackLink := "<https://github.com/ioki-mobility/summaraizer-slack|summaraizer-slack>"
	messageTmpl := "> This is a AI generated summarization of this thread. Powered by %s via %s:\n\n%s"
	return fmt.Sprintf(messageTmpl, summaraizerLink, summaraizerSlackLink, message)
}
