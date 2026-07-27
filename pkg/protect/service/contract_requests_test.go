package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/0sdknum/taurus-protect-sdk-go/internal/openapi"
	"github.com/0sdknum/taurus-protect-sdk-go/pkg/protect/model"
)

func TestRequestService_CreateOutgoingCallContractRequest_HTTPContract(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", request.Method)
		}
		if request.URL.Path != "/api/rest/v1/requests/outgoing/contracts/call" {
			t.Errorf("path = %s", request.URL.Path)
		}
		if request.Header.Get("Content-Type") != "application/json" {
			t.Errorf("content type = %q", request.Header.Get("Content-Type"))
		}

		var body map[string]any
		if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["fromAddressId"] != "address-1" || body["toWhitelistedAddressId"] != "whitelist-1" {
			t.Errorf("required IDs = %v", body)
		}
		if body["externalRequestId"] != "external-1" || body["amount"] != "0" {
			t.Errorf("idempotency or amount = %v", body)
		}
		if _, ok := body["call"]; ok {
			t.Errorf("legacy request contains generic call = %v", body["call"])
		}
		method := body["method"].(map[string]any)
		if method["functionSignature"] != "mint(address,uint256)" {
			t.Errorf("method = %v", method)
		}

		response.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(response).Encode(map[string]any{
			"result": map[string]any{"id": "request-1", "status": "CREATED"},
		})
	}))
	t.Cleanup(server.Close)

	service := NewRequestService(newServiceTestAPIClient(server))
	result, err := service.CreateOutgoingCallContractRequest(context.Background(), &model.CreateOutgoingCallContractRequest{
		FromAddressID:          "address-1",
		ToWhitelistedAddressID: "whitelist-1",
		Method:                 model.ContractCall{FunctionSignature: "mint(address,uint256)"},
		Amount:                 "0",
		ExternalRequestID:      "external-1",
	})
	if err != nil {
		t.Fatalf("CreateOutgoingCallContractRequest() error = %v", err)
	}
	if result.ID != "request-1" || result.Status != "CREATED" {
		t.Fatalf("result = %+v", result)
	}
}

func TestRequestService_CreateOutgoingDeployContractRequest_HTTPContract(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost || request.URL.Path != "/api/rest/v1/requests/outgoing/contracts/deploy" {
			t.Errorf("request = %s %s", request.Method, request.URL.Path)
		}

		var body map[string]any
		if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		contract := body["contract"].(map[string]any)
		if contract["blockchain"] != "ETH" {
			t.Errorf("contract = %v", contract)
		}
		if body["externalRequestId"] != "external-1" {
			t.Errorf("externalRequestId = %v", body["externalRequestId"])
		}

		response.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(response).Encode(map[string]any{
			"result": map[string]any{"id": "request-2", "status": "CREATED"},
		})
	}))
	t.Cleanup(server.Close)

	service := NewRequestService(newServiceTestAPIClient(server))
	result, err := service.CreateOutgoingDeployContractRequest(context.Background(), &model.CreateOutgoingDeployContractRequest{
		FromAddressID:     "address-1",
		ExternalRequestID: "external-1",
		Contract: model.GenericCreateContract{
			Blockchain: "ETH",
			ETH:        &model.CreateContractEVM{Bytecode: "Ynl0ZWNvZGU="},
		},
	})
	if err != nil {
		t.Fatalf("CreateOutgoingDeployContractRequest() error = %v", err)
	}
	if result.ID != "request-2" {
		t.Fatalf("result ID = %q", result.ID)
	}
}

func TestRequestService_ContractRequests_Validation(t *testing.T) {
	tests := []struct {
		name string
		call func(*RequestService) error
		want string
	}{
		{
			name: "nil call request",
			call: func(service *RequestService) error {
				_, err := service.CreateOutgoingCallContractRequest(context.Background(), nil)
				return err
			},
			want: "request cannot be nil",
		},
		{
			name: "call missing source",
			call: func(service *RequestService) error {
				_, err := service.CreateOutgoingCallContractRequest(context.Background(), &model.CreateOutgoingCallContractRequest{})
				return err
			},
			want: "fromAddressID is required",
		},
		{
			name: "call missing destination",
			call: func(service *RequestService) error {
				_, err := service.CreateOutgoingCallContractRequest(context.Background(), &model.CreateOutgoingCallContractRequest{FromAddressID: "address-1"})
				return err
			},
			want: "toWhitelistedAddressID is required",
		},
		{
			name: "legacy method and generic call",
			call: func(service *RequestService) error {
				_, err := service.CreateOutgoingCallContractRequest(context.Background(), &model.CreateOutgoingCallContractRequest{
					FromAddressID:          "address-1",
					ToWhitelistedAddressID: "whitelist-1",
					Method:                 model.ContractCall{FunctionSignature: "mint(address,uint256)"},
					Call:                   &model.GenericContractCall{Blockchain: "ETH"},
				})
				return err
			},
			want: "method and call are mutually exclusive",
		},
		{
			name: "nil deploy request",
			call: func(service *RequestService) error {
				_, err := service.CreateOutgoingDeployContractRequest(context.Background(), nil)
				return err
			},
			want: "request cannot be nil",
		},
		{
			name: "deploy missing source",
			call: func(service *RequestService) error {
				_, err := service.CreateOutgoingDeployContractRequest(context.Background(), &model.CreateOutgoingDeployContractRequest{})
				return err
			},
			want: "fromAddressID is required",
		},
		{
			name: "deploy missing blockchain",
			call: func(service *RequestService) error {
				_, err := service.CreateOutgoingDeployContractRequest(context.Background(), &model.CreateOutgoingDeployContractRequest{FromAddressID: "address-1"})
				return err
			},
			want: "contract blockchain is required",
		},
	}

	service := &RequestService{errMapper: NewErrorMapper()}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.call(service)
			if err == nil || err.Error() != test.want {
				t.Fatalf("error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestRequestService_CreateOutgoingCallContractRequest_RateLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		response.Header().Set("Retry-After", "4")
		response.WriteHeader(http.StatusTooManyRequests)
		_ = json.NewEncoder(response).Encode(map[string]string{"message": "rate limited"})
	}))
	t.Cleanup(server.Close)

	service := NewRequestService(newServiceTestAPIClient(server))
	_, err := service.CreateOutgoingCallContractRequest(context.Background(), &model.CreateOutgoingCallContractRequest{
		FromAddressID:          "address-1",
		ToWhitelistedAddressID: "whitelist-1",
	})
	if err == nil {
		t.Fatal("CreateOutgoingCallContractRequest() error = nil")
	}

	var apiError *APIError
	if !errors.As(err, &apiError) {
		t.Fatalf("error type = %T, want *APIError", err)
	}
	if apiError.Code != http.StatusTooManyRequests || apiError.RetryAfter != 4*time.Second {
		t.Fatalf("APIError = %+v", apiError)
	}
	if apiError.Message != "rate limited" || apiError.ResponseBody == "" {
		t.Fatalf("APIError details = %+v", apiError)
	}
}

func TestRequestService_CreateOutgoingCallContractRequest_BadRequestDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(response).Encode(map[string]any{
			"code":    3,
			"message": "invalid contract parameters",
			"details": []any{},
		})
	}))
	t.Cleanup(server.Close)

	service := NewRequestService(newServiceTestAPIClient(server))
	_, err := service.CreateOutgoingCallContractRequest(context.Background(), &model.CreateOutgoingCallContractRequest{
		FromAddressID:          "address-1",
		ToWhitelistedAddressID: "whitelist-1",
	})
	if err == nil {
		t.Fatal("CreateOutgoingCallContractRequest() error = nil")
	}

	var apiError *APIError
	if !errors.As(err, &apiError) {
		t.Fatalf("error type = %T, want *APIError", err)
	}
	if apiError.Message != "invalid contract parameters" || apiError.ErrorCode != "3" {
		t.Fatalf("APIError = %+v", apiError)
	}
	if apiError.ResponseBody == "" {
		t.Fatal("APIError response body is empty")
	}
}

func TestWhitelistedAddressService_CreateWhitelistedAddress_HTTPContract(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost || request.URL.Path != "/api/rest/v1/whitelists/addresses" {
			t.Errorf("request = %s %s", request.Method, request.URL.Path)
		}

		var body map[string]any
		if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["address"] != "0xabc" || body["label"] != "beneficiary" || body["blockchain"] != "ETH" {
			t.Errorf("required fields = %v", body)
		}
		if body["visibilityGroupID"] != "visibility-1" || body["customerId"] != "customer-1" {
			t.Errorf("optional fields = %v", body)
		}
		if len(body["linkedInternalAddressIds"].([]any)) != 2 || len(body["linkedWalletIds"].([]any)) != 1 {
			t.Errorf("linked IDs = %v", body)
		}

		response.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(response).Encode(map[string]any{"result": map[string]any{"id": "whitelist-1"}})
	}))
	t.Cleanup(server.Close)

	service := NewWhitelistedAddressServiceWithVerification(newServiceTestAPIClient(server), nil)
	id, err := service.CreateWhitelistedAddress(context.Background(), &model.CreateWhitelistedAddressRequest{
		Address:                  "0xabc",
		Memo:                     "memo",
		Label:                    "beneficiary",
		ExchangeAccountID:        "exchange-1",
		CustomerID:               "customer-1",
		LinkedInternalAddressIDs: []string{"address-1", "address-2"},
		AddressType:              "individual",
		ContractType:             "GENERIC",
		LinkedWalletIDs:          []string{"wallet-1"},
		Blockchain:               "ETH",
		Network:                  "mainnet",
		VisibilityGroupID:        "visibility-1",
	})
	if err != nil {
		t.Fatalf("CreateWhitelistedAddress() error = %v", err)
	}
	if id != "whitelist-1" {
		t.Fatalf("ID = %q", id)
	}
}

func TestWhitelistedAddressService_CreateWhitelistedAddress_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request *model.CreateWhitelistedAddressRequest
		want    string
	}{
		{name: "nil request", request: nil, want: "request cannot be nil"},
		{name: "missing address", request: &model.CreateWhitelistedAddressRequest{}, want: "address is required"},
		{name: "missing label", request: &model.CreateWhitelistedAddressRequest{Address: "0xabc"}, want: "label is required"},
		{name: "missing blockchain", request: &model.CreateWhitelistedAddressRequest{Address: "0xabc", Label: "label"}, want: "blockchain is required"},
	}

	service := &WhitelistedAddressService{errMapper: NewErrorMapper()}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := service.CreateWhitelistedAddress(context.Background(), test.request)
			if err == nil || err.Error() != test.want {
				t.Fatalf("error = %v, want %q", err, test.want)
			}
		})
	}
}

func newServiceTestAPIClient(server *httptest.Server) *openapi.APIClient {
	configuration := openapi.NewConfiguration()
	configuration.Servers = []openapi.ServerConfiguration{{URL: server.URL}}
	configuration.HTTPClient = server.Client()
	return openapi.NewAPIClient(configuration)
}
