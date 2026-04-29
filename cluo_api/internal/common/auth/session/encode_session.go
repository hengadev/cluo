package session

import (
	"encoding/json"
)

// EncodeSession encodes a Session to JSON for storage
func EncodeSession(sess *Session) ([]byte, error) {
	return json.Marshal(sess)
}
