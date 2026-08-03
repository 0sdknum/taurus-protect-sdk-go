package openapi

import (
	"encoding/json"
	"testing"
)

func TestTgvalidatordMetadataUnmarshalPayloadCompatibility(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		body        string
		wantPayload map[string]interface{}
		wantErr     bool
	}{
		{
			name:        "object payload remains available",
			body:        `{"hash":"hash","payload":{"key":"value"},"payloadAsString":"{\"key\":\"value\"}"}`,
			wantPayload: map[string]interface{}{"key": "value"},
		},
		{
			name: "array payload uses authoritative payload string",
			body: `{"hash":"hash","payload":[{"key":"currency","value":"ETH"}],"payloadAsString":"[{\"key\":\"currency\",\"value\":\"ETH\"}]"}`,
		},
		{
			name:    "array payload without payload string is rejected",
			body:    `{"hash":"hash","payload":[{"key":"currency","value":"ETH"}]}`,
			wantErr: true,
		},
		{
			name:    "scalar payload is rejected",
			body:    `{"hash":"hash","payload":"unexpected","payloadAsString":"unexpected"}`,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var metadata TgvalidatordMetadata
			err := json.Unmarshal([]byte(test.body), &metadata)
			if test.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unmarshal metadata: %v", err)
			}
			if metadata.PayloadAsString == nil {
				t.Fatal("payloadAsString was not preserved")
			}
			if len(test.wantPayload) == 0 {
				if metadata.Payload != nil {
					t.Fatalf("payload = %#v, want nil", metadata.Payload)
				}
				return
			}
			if metadata.Payload["key"] != test.wantPayload["key"] {
				t.Fatalf("payload = %#v, want %#v", metadata.Payload, test.wantPayload)
			}
		})
	}
}
