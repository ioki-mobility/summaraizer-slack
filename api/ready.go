package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func ReadyHandler(w http.ResponseWriter, r *http.Request) {
	validSlackToken := checkSlackBotToken()
	if !validSlackToken {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Invalid Slack Bot Token"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func checkSlackBotToken() bool {
	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	if slackBotToken == "" {
		return false
	}

	req, err := http.NewRequest("POST", "https://slack.com/api/auth.test", nil)
	if err != nil {
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+slackBotToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false
	}

	return response["ok"].(bool)
}
