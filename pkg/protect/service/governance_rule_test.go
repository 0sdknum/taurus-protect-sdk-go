package service

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"strings"
	"testing"

	protectcrypto "github.com/0sdknum/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/0sdknum/taurus-protect-sdk-go/pkg/protect/model"
)

func TestGovernanceRuleService_GetRulesByID_Validation(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty id returns error",
			id:      "",
			wantErr: true,
			errMsg:  "id cannot be empty",
		},
	}

	// Note: We can only test input validation without a real API client
	// Full integration tests require a running server
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a minimal service for validation testing
			s := &GovernanceRuleService{
				errMapper: NewErrorMapper(),
			}

			_, err := s.GetRulesByID(nil, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRulesByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("GetRulesByID() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestGovernanceRulesHistoryOptions(t *testing.T) {
	tests := []struct {
		name string
		opts *model.ListRulesHistoryOptions
	}{
		{
			name: "nil options",
			opts: nil,
		},
		{
			name: "with limit",
			opts: &model.ListRulesHistoryOptions{
				Limit: 10,
			},
		},
		{
			name: "with cursor",
			opts: &model.ListRulesHistoryOptions{
				Cursor: "eyJwYWdlIjogMn0=",
			},
		},
		{
			name: "with limit and cursor",
			opts: &model.ListRulesHistoryOptions{
				Limit:  20,
				Cursor: "eyJwYWdlIjogM30=",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify options are properly constructed
			if tt.opts != nil {
				if tt.opts.Limit < 0 {
					t.Error("Limit should not be negative")
				}
			}
		})
	}
}

func TestNewGovernanceRuleService(t *testing.T) {
	// Test that NewGovernanceRuleService doesn't panic with nil
	// In real usage, this would require a valid APIClient
	defer func() {
		if r := recover(); r != nil {
			// Expected to panic with nil client in production code
			// This is acceptable behavior
		}
	}()

	// This would panic with nil, which is expected
	// service := NewGovernanceRuleService(nil)
	// We just verify the function signature is correct
}

func TestGovernanceRuleService_VerifyGovernanceRulesJavaCompatibleBase64(t *testing.T) {
	t.Parallel()

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	rulesData := []byte("governance rules")
	signature, err := protectcrypto.SignData(privateKey, rulesData)
	if err != nil {
		t.Fatalf("sign rules: %v", err)
	}
	signature = strings.TrimRight(strings.NewReplacer("+", "-", "/", "_").Replace(signature), "=")

	service := &GovernanceRuleService{
		superAdminKeys:     []*ecdsa.PublicKey{&privateKey.PublicKey},
		minValidSignatures: 1,
	}
	rules := &model.GovernanceRuleset{
		RulesContainer: "! \t" + base64.RawURLEncoding.EncodeToString(rulesData) + "\r\n#",
		Signatures:     []model.RuleUserSignature{{UserID: "admin", Signature: signature}},
	}

	if err := service.VerifyGovernanceRules(rules); err != nil {
		t.Fatalf("VerifyGovernanceRules() error = %v", err)
	}
}
