package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

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
		handleEventCallback(w, r, body)
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

func handleEventCallback(w http.ResponseWriter, r *http.Request, body []byte) {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	url := fmt.Sprintf("%s://%s/event", scheme, r.Host)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", r.Header.Get("Content-Type"))

	client := &http.Client{
		Timeout: 150 * time.Millisecond,
	}
	_, err = client.Do(req)
}
