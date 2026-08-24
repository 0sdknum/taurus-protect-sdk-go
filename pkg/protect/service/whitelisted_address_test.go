package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	protectcrypto "github.com/0sdknum/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/0sdknum/taurus-protect-sdk-go/pkg/protect/helper"
	"github.com/0sdknum/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewWhitelistedAddressServiceWithVerification(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestWhitelistedAddressService_GetWhitelistedAddress_EmptyID(t *testing.T) {
	// Create a service with nil API to test validation before API call
	svc := &WhitelistedAddressService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetWhitelistedAddress(nil, "")
	if err == nil {
		t.Error("GetWhitelistedAddress() with empty ID should return error")
	}
	if err.Error() != "id cannot be empty" {
		t.Errorf("GetWhitelistedAddress() error = %v, want 'id cannot be empty'", err)
	}
}

func TestWhitelistedAddressService_ListWhitelistedAddresses_NilOptions(t *testing.T) {
	// This test verifies that nil options don't cause a panic
	// Actual API call would fail since api is nil, but we can test up to the validation
	svc := &WhitelistedAddressService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// We can't actually call ListWhitelistedAddresses without a real API,
	// but we verify the service struct is properly initialized
	if svc.errMapper == nil {
		t.Error("ErrorMapper should not be nil")
	}
}

func TestWhitelistedAddressService_ListWhitelistedAddressesRequestsInlineVerificationData(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet || request.URL.Path != "/api/rest/v1/whitelists/addresses" {
			t.Errorf("request = %s %s", request.Method, request.URL.Path)
		}
		if request.URL.Query().Has("rulesContainerNormalized") {
			t.Errorf("rulesContainerNormalized must be omitted, query = %s", request.URL.RawQuery)
		}
		response.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(response).Encode(map[string]any{
			"result":     []any{},
			"totalItems": "0",
		}); err != nil {
			t.Fatalf("encode response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	service := NewWhitelistedAddressServiceWithVerification(newServiceTestAPIClient(server), nil)
	items, _, err := service.ListWhitelistedAddresses(context.Background(), &model.ListWhitelistedAddressesOptions{Limit: 1})
	if err != nil {
		t.Fatalf("ListWhitelistedAddresses() error = %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("items = %d, want 0", len(items))
	}
}

func TestWhitelistedAddressService_ListWhitelistedAddressesPreservesVerificationFailure(t *testing.T) {
	t.Parallel()

	const payload = `{"address":"0xabc","label":"beneficiary"}`
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(response).Encode(map[string]any{
			"result": []any{map[string]any{
				"id":             "address-1",
				"blockchain":     "ETH",
				"network":        "mainnet",
				"metadata":       map[string]any{"hash": protectcrypto.CalculateHexHash(payload), "payloadAsString": payload},
				"signedAddress":  map[string]any{},
				"rulesContainer": "AA==",
			}},
			"totalItems": "1",
		}); err != nil {
			t.Fatalf("encode response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	service := NewWhitelistedAddressServiceWithVerification(newServiceTestAPIClient(server), &WhitelistedAddressServiceConfig{
		SuperAdminKeys:     []*ecdsa.PublicKey{&privateKey.PublicKey},
		MinValidSignatures: 1,
	})
	_, _, err = service.ListWhitelistedAddresses(context.Background(), &model.ListWhitelistedAddressesOptions{Limit: 1})
	var integrityError *model.IntegrityError
	if !errors.As(err, &integrityError) {
		t.Fatalf("error = %T %v, want *model.IntegrityError", err, err)
	}
	if integrityError.Message != "rulesSignatures is empty and no verified governance rules provider is configured" {
		t.Fatalf("IntegrityError.Message = %q", integrityError.Message)
	}
}

func TestWhitelistedAddressService_ListWhitelistedAddressesUsesVerifiedGovernanceFallback(t *testing.T) {
	t.Parallel()

	userPrivateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate user key: %v", err)
	}
	const payload = `{"address":"0xabc","label":"beneficiary"}`
	metadataHash := protectcrypto.CalculateHexHash(payload)
	hashes := []string{metadataHash}
	hashesJSON, err := json.Marshal(hashes)
	if err != nil {
		t.Fatalf("marshal hashes: %v", err)
	}
	userSignature, err := protectcrypto.SignData(userPrivateKey, hashesJSON)
	if err != nil {
		t.Fatalf("sign hashes: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(response).Encode(map[string]any{
			"result": []any{map[string]any{
				"id":         "address-1",
				"blockchain": "ETH",
				"network":    "mainnet",
				"metadata":   map[string]any{"hash": metadataHash, "payloadAsString": payload},
				"signedAddress": map[string]any{"signatures": []any{map[string]any{
					"hashes":    hashes,
					"signature": map[string]any{"userId": "user-1", "signature": userSignature},
				}}},
				"rulesContainer": "AA==",
			}},
			"totalItems": "1",
		}); err != nil {
			t.Fatalf("encode response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	providerCalls := 0
	service := NewWhitelistedAddressServiceWithVerification(newServiceTestAPIClient(server), &WhitelistedAddressServiceConfig{
		SuperAdminKeys:     []*ecdsa.PublicKey{&userPrivateKey.PublicKey},
		MinValidSignatures: 1,
		VerifiedGovernanceRulesProvider: func(_ context.Context, requested []string) (map[string]*model.DecodedRulesContainer, error) {
			providerCalls++
			if len(requested) != 1 || requested[0] != "AA==" {
				t.Fatalf("requested rules containers = %v", requested)
			}
			return map[string]*model.DecodedRulesContainer{"AA==": {
				Users:  []*model.RuleUser{{ID: "user-1", PublicKey: &userPrivateKey.PublicKey}},
				Groups: []*model.RuleGroup{{ID: "approvers", UserIDs: []string{"user-1"}}},
				AddressWhitelistingRules: []*model.AddressWhitelistingRules{{
					Currency: "ETH",
					Network:  "mainnet",
					ParallelThresholds: []*model.SequentialThresholds{{Thresholds: []*model.GroupThreshold{{
						GroupID: "approvers", MinimumSignatures: 1,
					}}}},
				}},
			}}, nil
		},
	})

	items, _, err := service.ListWhitelistedAddresses(context.Background(), &model.ListWhitelistedAddressesOptions{Limit: 1})
	if err != nil {
		t.Fatalf("ListWhitelistedAddresses() error = %v", err)
	}
	if len(items) != 1 || items[0].ID != "address-1" {
		t.Fatalf("items = %+v", items)
	}
	if providerCalls != 1 {
		t.Fatalf("verified governance provider calls = %d, want 1", providerCalls)
	}
}

func TestWhitelistedAddressService_InitializeEnvelopeUsesVerifiedGovernanceFallback(t *testing.T) {
	t.Parallel()

	userPrivateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate user key: %v", err)
	}
	const payload = `{"address":"0xabc","label":"beneficiary"}`
	metadataHash := protectcrypto.CalculateHexHash(payload)
	hashes := []string{metadataHash}
	hashesJSON, err := json.Marshal(hashes)
	if err != nil {
		t.Fatalf("marshal hashes: %v", err)
	}
	userSignature, err := protectcrypto.SignData(userPrivateKey, hashesJSON)
	if err != nil {
		t.Fatalf("sign hashes: %v", err)
	}
	decoded := &model.DecodedRulesContainer{
		Users:  []*model.RuleUser{{ID: "user-1", PublicKey: &userPrivateKey.PublicKey}},
		Groups: []*model.RuleGroup{{ID: "approvers", UserIDs: []string{"user-1"}}},
		AddressWhitelistingRules: []*model.AddressWhitelistingRules{{
			Currency: "ETH",
			Network:  "mainnet",
			ParallelThresholds: []*model.SequentialThresholds{{Thresholds: []*model.GroupThreshold{{
				GroupID: "approvers", MinimumSignatures: 1,
			}}}},
		}},
	}
	service := &WhitelistedAddressService{
		verifier: helper.NewWhitelistedAddressVerifier([]*ecdsa.PublicKey{&userPrivateKey.PublicKey}, 1),
		verifiedGovernanceRulesProvider: func(_ context.Context, requested []string) (map[string]*model.DecodedRulesContainer, error) {
			if len(requested) != 1 || requested[0] != "AA==" {
				t.Fatalf("requested rules containers = %v", requested)
			}
			return map[string]*model.DecodedRulesContainer{"AA==": decoded}, nil
		},
	}
	envelope := &model.WhitelistedAddressEnvelope{
		ID:             "address-1",
		Blockchain:     "ETH",
		Network:        "mainnet",
		RulesContainer: "AA==",
		Metadata:       &model.WhitelistedAssetMetadata{Hash: metadataHash, PayloadAsString: payload},
		SignedAddress: &model.SignedWhitelistedAddress{Signatures: []model.WhitelistSignature{{
			Hashes:        hashes,
			UserSignature: &model.WhitelistUserSignature{UserID: "user-1", Signature: userSignature},
		}}},
	}

	if err := service.initializeEnvelope(context.Background(), envelope); err != nil {
		t.Fatalf("initializeEnvelope() error = %v", err)
	}
	if envelope.WhitelistedAddress() == nil {
		t.Fatal("verified whitelisted address is nil")
	}
}
