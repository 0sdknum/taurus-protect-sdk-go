package mapper

import (
	"encoding/json"
	"testing"

	"github.com/0sdknum/taurus-protect-sdk-go/pkg/protect/model"
)

func TestCreateOutgoingCallContractRequestToDTO_AllFields(t *testing.T) {
	request := &model.CreateOutgoingCallContractRequest{
		FromAddressID:          "address-1",
		ToWhitelistedAddressID: "whitelist-1",
		Method: model.ContractCall{
			FunctionSignature: "transfer(address,uint256)",
			Arguments: []model.ContractArgument{{
				Name: "recipients",
				Type: "address[]",
				Value: &model.ContractArgumentValue{Composite: []model.ContractArgumentValue{
					{Primitive: "0xabc"},
					{Primitive: "0xdef"},
				}},
			}},
		},
		GasLimit:             "100000",
		GasPriceLimit:        "200",
		Comment:              "mint",
		ContractType:         "GENERIC",
		Amount:               "0",
		FeePayerID:           "fee-payer-1",
		FeeLimit:             "300",
		TransactionReference: "legacy-reference",
		ExternalRequestID:    "external-1",
		Call: &model.GenericContractCall{
			Blockchain: "ETH",
			ETH: &model.ContractCall{
				FunctionSignature: "mint(address,uint256)",
				Arguments: []model.ContractArgument{{
					Name:  "amount",
					Type:  "uint256",
					Value: &model.ContractArgumentValue{Primitive: "100"},
				}},
			},
		},
	}

	data, err := json.Marshal(CreateOutgoingCallContractRequestToDTO(request))
	if err != nil {
		t.Fatalf("marshal DTO: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal DTO: %v", err)
	}
	if got["fromAddressId"] != "address-1" || got["toWhitelistedAddressId"] != "whitelist-1" {
		t.Fatalf("required IDs not mapped: %v", got)
	}
	if got["externalRequestId"] != "external-1" || got["amount"] != "0" {
		t.Fatalf("idempotency or amount not mapped: %v", got)
	}
	call := got["call"].(map[string]any)
	eth := call["eth"].(map[string]any)
	if eth["functionSignature"] != "mint(address,uint256)" {
		t.Fatalf("call function signature = %v", eth["functionSignature"])
	}
	method := got["method"].(map[string]any)
	arguments := method["args"].([]any)
	value := arguments[0].(map[string]any)["value"].(map[string]any)
	if len(value["composite"].([]any)) != 2 {
		t.Fatalf("composite arguments = %v", value["composite"])
	}
}

func TestCreateOutgoingDeployContractRequestToDTO_XTZ(t *testing.T) {
	request := &model.CreateOutgoingDeployContractRequest{
		FromAddressID:              "address-1",
		GenerateWhitelistedAddress: true,
		ExternalRequestID:          "external-1",
		Contract: model.GenericCreateContract{
			Blockchain: "XTZ",
			XTZ: &model.CreateContractXTZ{
				Code: "contract-code",
				Storage: &model.XTZContractArgument{
					Primitive: "pair",
					Arguments: []model.XTZContractArgument{{String: "value"}},
				},
				Delegate: &model.XTZContractDelegate{ToWhitelistedAddressID: "whitelist-1"},
			},
		},
	}

	data, err := json.Marshal(CreateOutgoingDeployContractRequestToDTO(request))
	if err != nil {
		t.Fatalf("marshal DTO: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal DTO: %v", err)
	}
	if got["generateWhitelistedAddress"] != true {
		t.Fatalf("generateWhitelistedAddress = %v", got["generateWhitelistedAddress"])
	}
	contract := got["contract"].(map[string]any)
	xtz := contract["xtz"].(map[string]any)
	delegate := xtz["delegate"].(map[string]any)
	if delegate["toWhitelistedAddressId"] != "whitelist-1" {
		t.Fatalf("delegate = %v", delegate)
	}
	storage := xtz["storage"].(map[string]any)
	if storage["prim"] != "pair" {
		t.Fatalf("storage = %v", storage)
	}
}
