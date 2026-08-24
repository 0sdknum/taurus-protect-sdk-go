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
	if integrityError.Message != "rulesSignatures is empty" {
		t.Fatalf("IntegrityError.Message = %q", integrityError.Message)
	}
}
