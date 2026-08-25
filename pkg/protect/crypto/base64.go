package crypto

import "encoding/base64"

func isBase64Alphabet(character byte) bool {
	return character >= 'A' && character <= 'Z' ||
		character >= 'a' && character <= 'z' ||
		character >= '0' && character <= '9' ||
		character == '+' || character == '/'
}

// DecodeBase64 matches Apache Commons Codec Base64.decodeBase64 behavior used by the Java SDK.
func DecodeBase64(encoded string) ([]byte, error) {
	normalized := make([]byte, 0, len(encoded))

decode:
	for index := 0; index < len(encoded); index++ {
		character := encoded[index]
		switch {
		case isBase64Alphabet(character):
			normalized = append(normalized, character)
		case character == '-':
			normalized = append(normalized, '+')
		case character == '_':
			normalized = append(normalized, '/')
		case character == '=':
			break decode
		}
	}

	if len(normalized)%4 == 1 {
		normalized = normalized[:len(normalized)-1]
	}

	return base64.RawStdEncoding.DecodeString(string(normalized))
}
