package openapi

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// UnmarshalJSON accepts the array-shaped request metadata returned by some
// Taurus-PROTECT deployments while preserving object payloads used elsewhere.
// Request consumers must use and verify payloadAsString; array payload is not
// exposed through the object-only Payload field.
func (o *TgvalidatordMetadata) UnmarshalJSON(data []byte) error {
	var wire struct {
		Hash            *string         `json:"hash,omitempty"`
		Payload         json.RawMessage `json:"payload,omitempty"`
		PayloadAsString *string         `json:"payloadAsString,omitempty"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}

	var payload map[string]interface{}
	rawPayload := bytes.TrimSpace(wire.Payload)
	if len(rawPayload) > 0 && !bytes.Equal(rawPayload, []byte("null")) {
		switch rawPayload[0] {
		case '{':
			if err := json.Unmarshal(rawPayload, &payload); err != nil {
				return fmt.Errorf("decode metadata payload object: %w", err)
			}
		case '[':
			if wire.PayloadAsString == nil || *wire.PayloadAsString == "" {
				return fmt.Errorf("decode metadata payload array: payloadAsString is required")
			}
		default:
			return fmt.Errorf("decode metadata payload: expected object or array")
		}
	}

	o.Hash = wire.Hash
	o.Payload = payload
	o.PayloadAsString = wire.PayloadAsString
	return nil
}
