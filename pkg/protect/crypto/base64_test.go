package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"strings"
	"testing"
)

func TestDecodeBase64JavaCompatible(t *testing.T) {
	t.Parallel()

	want := []byte{0xfb, 0xff, 0x00, 0x01}
	standard := base64.StdEncoding.EncodeToString(want)
	tests := []struct {
		name      string
		encoded   string
		wantEmpty bool
	}{
		{name: "standard padded", encoded: standard},
		{name: "standard unpadded", encoded: strings.TrimRight(standard, "=")},
		{name: "URL safe unpadded", encoded: base64.RawURLEncoding.EncodeToString(want)},
		{name: "ASCII whitespace", encoded: " \t" + standard[:4] + "\r\n" + standard[4:]},
		{name: "unsupported bytes ignored", encoded: "!" + standard[:4] + "@" + standard[4:] + "#"},
		{name: "incomplete quantum", encoded: "A", wantEmpty: true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := DecodeBase64(test.encoded)
			if err != nil {
				t.Fatalf("DecodeBase64() error = %v", err)
			}
			expected := want
			if test.wantEmpty {
				expected = nil
			}
			if !bytes.Equal(got, expected) {
				t.Fatalf("DecodeBase64() = %x, want %x", got, expected)
			}
		})
	}
}

func TestVerifySignatureJavaCompatibleBase64(t *testing.T) {
	t.Parallel()

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	data := []byte("governance rules")
	signature, err := SignData(privateKey, data)
	if err != nil {
		t.Fatalf("sign data: %v", err)
	}

	tests := []struct {
		name      string
		signature string
	}{
		{name: "standard padded", signature: signature},
		{name: "standard unpadded", signature: strings.TrimRight(signature, "=")},
		{name: "URL safe unpadded", signature: strings.TrimRight(strings.NewReplacer("+", "-", "/", "_").Replace(signature), "=")},
		{name: "ASCII whitespace", signature: " \t" + signature[:4] + "\r\n" + signature[4:]},
		{name: "unsupported bytes ignored", signature: "!" + signature[:4] + "@" + signature[4:] + "#"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			// Lenient Base64 decoding produces an empty signature; the fixed-width r||s gate must reject it.
			valid, err := VerifySignature(&privateKey.PublicKey, data, test.signature)
			if err != nil {
				t.Fatalf("VerifySignature() error = %v", err)
			}
			if !valid {
				t.Fatal("VerifySignature() = false, want true")
			}
		})
	}
}

func TestVerifySignatureRejectsIncorrectSignatureLength(t *testing.T) {
	t.Parallel()

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	valid, err := VerifySignature(&privateKey.PublicKey, []byte("governance rules"), "A!")
	if err == nil {
		t.Fatal("VerifySignature() error = nil")
	}
	if valid {
		t.Fatal("VerifySignature() = true, want false")
	}
}
