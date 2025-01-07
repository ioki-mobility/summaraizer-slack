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

const aiPrompt = `
I give you a discussion and you give me a summary.
Each comment of the discussion is wrapped in a <comment> tag.
Your summary should not be longer than 1200 chars.
Here is the discussion:
{{ range $comment := . }}
<comment>{{ $comment.Body }}</comment>
{{end}}
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
		threadTs, _ := innerEvent["thread_ts"].(string)
		channel, _ := innerEvent["channel"].(string)

		if threadTs != "" && strings.Contains(strings.ToLower(text), "summarize") {
			log.Printf("Summarize request in channel %s, thread %s by %s", channel, threadTs, user)
			messageTs := slack.SendMessage(":brain: Thinking...", channel, threadTs, slackBotToken)
			summarization := fetchAndSummarize(channel, threadTs)
			slack.UpdateMessage(messageTemplate(summarization), channel, messageTs, slackBotToken)
		}
	}
}

func fetchAndSummarize(channel, threadTs string) string {
	buffer := bytes.Buffer{}
	slack := summaraizer.Slack{
		Token:   slackBotToken,
		Channel: channel,
		TS:      threadTs,
	}
	slack.Fetch(&buffer)

	var summarization string
	var err error
	var summarizer summaraizer.Summarizer
	switch {
	case openAiToken != "":
		summarizer = &summaraizer.OpenAi{
			Model:    "gpt-4o-mini",
			Prompt:   aiPrompt,
			ApiToken: openAiToken,
		}
		break
	case ollamUrl != "":
		summarizer = &summaraizer.Ollama{
			Model:  "llama3.1:latest",
			Prompt: aiPrompt,
			Url:    ollamUrl,
		}
		break
	}
	if summarizer == nil {
		log.Fatal("OpenAiToken AND OllamaUrl are nil. Please set one of them.")
	}

	summarization, err = summarizer.Summarize(&buffer)

	if err != nil {
		log.Fatal(err)
	}

	return summarization
}

func messageTemplate(message string) string {
	summaraizerLink := "<https://github.com/ioki-mobility/summaraizer|summaraizer>"
	summaraizerSlackLink := "<https://github.com/ioki-mobility/summaraizer-slack|summaraizer-slack>"
	messageTmpl := "> This is a AI generated summarization of this thread. Powered by %s via %s:\n\n%s"
	return fmt.Sprintf(messageTmpl, summaraizerLink, summaraizerSlackLink, message)
}
