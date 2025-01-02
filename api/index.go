package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ioki-mobility/summaraizer-slack/slack"
)

var slackSigningSecretIndex = os.Getenv("SLACK_SIGNING_SECRET")

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if valid := slack.VerifySignature(r.Header, body, slackSigningSecretIndex); valid != true {
		log.Printf("Slack signature doesn't match!")
		return
	}

	var event map[string]interface{}
	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	switch event["type"] {
	case "url_verification":
		handleURLVerification(w, event)
	case "event_callback":
		acknowledgeSlackRequest(w)
		handleEventCallback(r, body)
	default:
		http.Error(w, "Unknown event type", http.StatusBadRequest)
	}
}

func handleURLVerification(w http.ResponseWriter, event map[string]interface{}) {
	challenge, ok := event["challenge"].(string)
	if !ok {
		http.Error(w, "Invalid challenge", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(challenge))
}

// For Slack to acknowledge the request, we need to respond with a 200 status code.
// See https://api.slack.com/apis/events-api#responding
func acknowledgeSlackRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func handleEventCallback(r *http.Request, body []byte) {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	req := r.Clone(r.Context())
	req.RequestURI = ""
	req.URL, _ = url.Parse(fmt.Sprintf("%s://%s/event", scheme, r.Host))
	req.Method = "POST"
	req.Body = io.NopCloser(strings.NewReader(string(body)))

	client := &http.Client{
		Timeout: 150 * time.Millisecond,
	}
	client.Do(req)
}
