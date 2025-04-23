package security

import "encoding/base64"

func EncodeTokenURLSafe(token string) string {
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(token))
}

func DecodeTokenURLSafe(encoded string) (string, error) {
	decoded, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
