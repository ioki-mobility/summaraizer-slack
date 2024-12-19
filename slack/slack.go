package slack

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
)

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
