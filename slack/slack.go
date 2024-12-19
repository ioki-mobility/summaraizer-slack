package slack

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const slackApiUrl = "https://slack.com/api/"

// VerifySignature will verify request comes from slack.
// See also https://api.slack.com/authentication/verifying-requests-from-slack
func VerifySignature(headers http.Header, body []byte, signingSecret string) bool {
	versionNumber := "v0"

	timestamp := headers.Get("X-Slack-Request-Timestamp")

	expectedSignature := headers.Get("X-Slack-Signature")

	textToEncrypt := fmt.Sprintf(
		"%s:%s:%s",
		versionNumber,
		timestamp,
		string(body),
	)

	hash := hmac.New(sha256.New, []byte(signingSecret))
	hash.Write([]byte(textToEncrypt))
	encryptedResult := hex.EncodeToString(hash.Sum(nil))
	encryptedResultWithVersionNumber := "v0=" + encryptedResult

	return hmac.Equal(
		[]byte(expectedSignature),
		[]byte(encryptedResultWithVersionNumber),
	)
}

func SendMessage(
	message string,
	channel string,
	threadTs string,
	slackBotToken string,
) {
	payload := map[string]interface{}{
		"channel":   channel,
		"text":      message,
		"thread_ts": threadTs,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return
	}

	req, err := http.NewRequest("POST", slackApiUrl+"chat.postMessage", strings.NewReader(string(jsonPayload)))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+slackBotToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Slack API response: %s", string(body))
}
