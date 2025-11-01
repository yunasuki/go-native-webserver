package auth

import "encoding/json"

// Helper for decoding Basic Auth
type basicAuth struct {
	Email    string
	Password string
}

func DecodeBasicAuth(encoded string) (basicAuth, error) {
	b, err := decodeBase64(encoded)
	if err != nil {
		return basicAuth{}, err
	}
	parts := splitUserPass(string(b))
	if len(parts) != 2 {
		return basicAuth{}, err
	}
	return basicAuth{Email: parts[0], Password: parts[1]}, nil
}

func decodeBase64(s string) ([]byte, error) {
	return json.Marshal(s) // replace with base64 decoding
}

func splitUserPass(s string) []string {
	return []string{"", ""} // replace with strings.SplitN(s, ":", 2)
}
