package mapper

import (
	"github.com/0sdknum/taurus-protect-sdk-go/internal/openapi"
	"github.com/0sdknum/taurus-protect-sdk-go/pkg/protect/model"
)

// CreateOutgoingCallContractRequestToDTO maps a public contract call request to OpenAPI.
func CreateOutgoingCallContractRequestToDTO(request *model.CreateOutgoingCallContractRequest) openapi.TgvalidatordCreateOutgoingCallContractRequestRequest {
	result := openapi.TgvalidatordCreateOutgoingCallContractRequestRequest{
		FromAddressId:          request.FromAddressID,
		ToWhitelistedAddressId: request.ToWhitelistedAddressID,
		Method:                 ContractCallToDTO(request.Method),
	}
	if request.GasLimit != "" {
		result.GasLimit = new(request.GasLimit)
	}
	if request.GasPriceLimit != "" {
		result.GasPriceLimit = new(request.GasPriceLimit)
	}
	if request.Comment != "" {
		result.Comment = new(request.Comment)
	}
	if request.ContractType != "" {
		result.ContractType = new(request.ContractType)
	}
	if request.Amount != "" {
		result.Amount = new(request.Amount)
	}
	if request.FeePayerID != "" {
		result.FeePayerId = new(request.FeePayerID)
	}
	if request.FeeLimit != "" {
		result.FeeLimit = new(request.FeeLimit)
	}
	if request.Call != nil {
		result.Call = new(GenericContractCallToDTO(*request.Call))
	}
	if request.TransactionReference != "" {
		result.TransactionReference = new(request.TransactionReference)
	}
	if request.ExternalRequestID != "" {
		result.ExternalRequestId = new(request.ExternalRequestID)
	}
	return result
}

// CreateOutgoingDeployContractRequestToDTO maps a public deployment request to OpenAPI.
func CreateOutgoingDeployContractRequestToDTO(request *model.CreateOutgoingDeployContractRequest) openapi.TgvalidatordCreateOutgoingDeployContractRequestRequest {
	result := openapi.TgvalidatordCreateOutgoingDeployContractRequestRequest{
		FromAddressId:              request.FromAddressID,
		GenerateWhitelistedAddress: new(request.GenerateWhitelistedAddress),
		Contract:                   GenericCreateContractToDTO(request.Contract),
	}
	if request.Bytecode != "" {
		result.Bytecode = new(request.Bytecode)
	}
	if request.Constructor != nil {
		result.Constructor = new(ContractCallToDTO(*request.Constructor))
	}

	if request.GasLimit != "" {
		result.GasLimit = new(request.GasLimit)
	}
	if request.GasPriceLimit != "" {
		result.GasPriceLimit = new(request.GasPriceLimit)
	}
	if request.Comment != "" {
		result.Comment = new(request.Comment)
	}
	if request.ContractType != "" {
		result.ContractType = new(request.ContractType)
	}
	if request.FeePayerID != "" {
		result.FeePayerId = new(request.FeePayerID)
	}
	if request.FeeLimit != "" {
		result.FeeLimit = new(request.FeeLimit)
	}
	if request.TransactionReference != "" {
		result.TransactionReference = new(request.TransactionReference)
	}
	if request.ExternalRequestID != "" {
		result.ExternalRequestId = new(request.ExternalRequestID)
	}
	return result
}

// ContractCallToDTO maps an EVM-compatible contract call to OpenAPI.
func ContractCallToDTO(call model.ContractCall) openapi.TgvalidatordContractCall {
	result := openapi.TgvalidatordContractCall{}
	if call.FunctionSignature != "" {
		result.FunctionSignature = new(call.FunctionSignature)
	}
	if len(call.Arguments) > 0 {
		result.Args = make([]openapi.TgvalidatordContractArg, len(call.Arguments))
		for index := range call.Arguments {
			result.Args[index] = contractArgumentToDTO(call.Arguments[index])
		}
	}
	return result
}

// GenericContractCallToDTO maps chain-specific call data to OpenAPI.
func GenericContractCallToDTO(call model.GenericContractCall) openapi.TgvalidatordGenericContractCall {
	result := openapi.TgvalidatordGenericContractCall{}
	if call.Blockchain != "" {
		result.Blockchain = new(call.Blockchain)
	}
	if call.ETH != nil {
		result.Eth = new(ContractCallToDTO(*call.ETH))
	}
	if call.XTZ != nil {
		result.Xtz = new(xtzContractCallToDTO(*call.XTZ))
	}
	if call.EVM != nil {
		result.Evm = new(ContractCallToDTO(*call.EVM))
	}
	return result
}

// GenericCreateContractToDTO maps chain-specific deployment data to OpenAPI.
func GenericCreateContractToDTO(contract model.GenericCreateContract) openapi.TgvalidatordGenericCreateContract {
	result := openapi.TgvalidatordGenericCreateContract{Blockchain: contract.Blockchain}
	if contract.ETH != nil {
		result.Eth = new(createContractEVMToDTO(*contract.ETH))
	}
	if contract.XTZ != nil {
		result.Xtz = new(createContractXTZToDTO(*contract.XTZ))
	}
	if contract.EVM != nil {
		result.Evm = new(createContractEVMToDTO(*contract.EVM))
	}
	return result
}

func contractArgumentToDTO(argument model.ContractArgument) openapi.TgvalidatordContractArg {
	result := openapi.TgvalidatordContractArg{}
	if argument.Name != "" {
		result.Name = new(argument.Name)
	}
	if argument.Type != "" {
		result.Type = new(argument.Type)
	}
	if argument.Value != nil {
		result.Value = new(contractArgumentValueToDTO(*argument.Value))
	}
	return result
}

func contractArgumentValueToDTO(value model.ContractArgumentValue) openapi.TgvalidatordContractArgValue {
	result := openapi.TgvalidatordContractArgValue{}
	if value.Primitive != "" {
		result.Primitive = new(value.Primitive)
	}
	if len(value.Composite) > 0 {
		result.Composite = make([]openapi.TgvalidatordContractArgValue, len(value.Composite))
		for index := range value.Composite {
			result.Composite[index] = contractArgumentValueToDTO(value.Composite[index])
		}
	}
	return result
}

func xtzContractCallToDTO(call model.XTZContractCall) openapi.TgvalidatordXTZContractCall {
	result := openapi.TgvalidatordXTZContractCall{}
	if call.Entrypoint != "" {
		result.Entrypoint = new(call.Entrypoint)
	}
	if call.Argument != nil {
		result.Arg = new(xtzContractArgumentToDTO(*call.Argument))
	}
	return result
}

func xtzContractArgumentToDTO(argument model.XTZContractArgument) openapi.TgvalidatordXTZContractArg {
	result := openapi.TgvalidatordXTZContractArg{}
	if argument.Kind != "" {
		result.Kind = new(argument.Kind)
	}
	if argument.Primitive != "" {
		result.Prim = new(argument.Primitive)
	}
	if len(argument.Arguments) > 0 {
		result.Args = make([]openapi.TgvalidatordXTZContractArg, len(argument.Arguments))
		for index := range argument.Arguments {
			result.Args[index] = xtzContractArgumentToDTO(argument.Arguments[index])
		}
	}
	if argument.Source != nil {
		result.Source = &openapi.TgvalidatordXTZContractArgSource{}
		if argument.Source.FromAddressID != "" {
			result.Source.FromAddressId = new(argument.Source.FromAddressID)
		}
	}
	if argument.Destination != nil {
		result.Destination = &openapi.TgvalidatordXTZContractArgDestination{}
		if argument.Destination.ToAddressID != "" {
			result.Destination.ToAddressId = new(argument.Destination.ToAddressID)
		}
		if argument.Destination.ToWhitelistedAddressID != "" {
			result.Destination.ToWhitelistedAddressId = new(argument.Destination.ToWhitelistedAddressID)
		}
	}
	if argument.String != "" {
		result.String = new(argument.String)
	}
	if argument.Integer != "" {
		result.Int = new(argument.Integer)
	}
	if argument.Bytes != "" {
		result.Bytes = new(argument.Bytes)
	}
	if len(argument.Annotations) > 0 {
		result.Annotations = append([]string(nil), argument.Annotations...)
	}
	return result
}

func createContractEVMToDTO(contract model.CreateContractEVM) openapi.GenericCreateContractEVMContract {
	result := openapi.GenericCreateContractEVMContract{}
	if contract.Bytecode != "" {
		result.Bytecode = new(contract.Bytecode)
	}
	if contract.Constructor != nil {
		result.Constructor = new(ContractCallToDTO(*contract.Constructor))
	}
	return result
}

func createContractXTZToDTO(contract model.CreateContractXTZ) openapi.GenericCreateContractXTZContract {
	result := openapi.GenericCreateContractXTZContract{}
	if contract.Code != "" {
		result.Code = new(contract.Code)
	}
	if contract.Storage != nil {
		result.Storage = new(xtzContractArgumentToDTO(*contract.Storage))
	}
	if contract.Delegate != nil {
		result.Delegate = &openapi.XTZContractDelegate{}
		if contract.Delegate.ToAddressID != "" {
			result.Delegate.ToAddressId = new(contract.Delegate.ToAddressID)
		}
		if contract.Delegate.ToWhitelistedAddressID != "" {
			result.Delegate.ToWhitelistedAddressId = new(contract.Delegate.ToWhitelistedAddressID)
		}
	}
	return result
}
