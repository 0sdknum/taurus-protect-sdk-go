package model

// ContractArgumentValue represents a primitive or composite smart-contract argument value.
type ContractArgumentValue struct {
	Primitive string                  `json:"primitive,omitempty"`
	Composite []ContractArgumentValue `json:"composite,omitempty"`
}

// ContractArgument represents a named, typed smart-contract argument.
type ContractArgument struct {
	Name  string                 `json:"name,omitempty"`
	Type  string                 `json:"type,omitempty"`
	Value *ContractArgumentValue `json:"value,omitempty"`
}

// ContractCall describes an EVM-compatible smart-contract method call.
type ContractCall struct {
	FunctionSignature string             `json:"function_signature,omitempty"`
	Arguments         []ContractArgument `json:"arguments,omitempty"`
}

// XTZContractArgumentSource references the source address for a Tezos argument.
type XTZContractArgumentSource struct {
	FromAddressID string `json:"from_address_id,omitempty"`
}

// XTZContractArgumentDestination references an internal or whitelisted destination.
type XTZContractArgumentDestination struct {
	ToAddressID            string `json:"to_address_id,omitempty"`
	ToWhitelistedAddressID string `json:"to_whitelisted_address_id,omitempty"`
}

// XTZContractArgument represents a recursive Tezos Micheline argument.
type XTZContractArgument struct {
	Kind        string                          `json:"kind,omitempty"`
	Primitive   string                          `json:"primitive,omitempty"`
	Arguments   []XTZContractArgument           `json:"arguments,omitempty"`
	Source      *XTZContractArgumentSource      `json:"source,omitempty"`
	Destination *XTZContractArgumentDestination `json:"destination,omitempty"`
	String      string                          `json:"string,omitempty"`
	Integer     string                          `json:"integer,omitempty"`
	Bytes       string                          `json:"bytes,omitempty"`
	Annotations []string                        `json:"annotations,omitempty"`
}

// XTZContractCall describes a Tezos contract entrypoint call.
type XTZContractCall struct {
	Entrypoint string               `json:"entrypoint,omitempty"`
	Argument   *XTZContractArgument `json:"argument,omitempty"`
}

// GenericContractCall contains the chain-specific representation of a contract call.
type GenericContractCall struct {
	Blockchain string           `json:"blockchain,omitempty"`
	ETH        *ContractCall    `json:"eth,omitempty"`
	XTZ        *XTZContractCall `json:"xtz,omitempty"`
	EVM        *ContractCall    `json:"evm,omitempty"`
}

// CreateContractEVM contains EVM contract deployment data.
type CreateContractEVM struct {
	Bytecode    string        `json:"bytecode,omitempty"`
	Constructor *ContractCall `json:"constructor,omitempty"`
}

// XTZContractDelegate references the delegate for a Tezos contract deployment.
type XTZContractDelegate struct {
	ToAddressID            string `json:"to_address_id,omitempty"`
	ToWhitelistedAddressID string `json:"to_whitelisted_address_id,omitempty"`
}

// CreateContractXTZ contains Tezos contract deployment data.
type CreateContractXTZ struct {
	Code     string               `json:"code,omitempty"`
	Storage  *XTZContractArgument `json:"storage,omitempty"`
	Delegate *XTZContractDelegate `json:"delegate,omitempty"`
}

// GenericCreateContract contains chain-specific contract deployment data.
type GenericCreateContract struct {
	Blockchain string             `json:"blockchain"`
	ETH        *CreateContractEVM `json:"eth,omitempty"`
	XTZ        *CreateContractXTZ `json:"xtz,omitempty"`
	EVM        *CreateContractEVM `json:"evm,omitempty"`
}

// CreateOutgoingCallContractRequest contains parameters for creating a contract call request.
type CreateOutgoingCallContractRequest struct {
	FromAddressID          string               `json:"from_address_id"`
	ToWhitelistedAddressID string               `json:"to_whitelisted_address_id"`
	Method                 ContractCall         `json:"method"`
	GasLimit               string               `json:"gas_limit,omitempty"`
	GasPriceLimit          string               `json:"gas_price_limit,omitempty"`
	Comment                string               `json:"comment,omitempty"`
	ContractType           string               `json:"contract_type,omitempty"`
	Amount                 string               `json:"amount,omitempty"`
	FeePayerID             string               `json:"fee_payer_id,omitempty"`
	FeeLimit               string               `json:"fee_limit,omitempty"`
	Call                   *GenericContractCall `json:"call,omitempty"`
	TransactionReference   string               `json:"transaction_reference,omitempty"`
	ExternalRequestID      string               `json:"external_request_id,omitempty"`
}

// CreateOutgoingDeployContractRequest contains parameters for creating a contract deployment request.
type CreateOutgoingDeployContractRequest struct {
	FromAddressID              string                `json:"from_address_id"`
	Bytecode                   string                `json:"bytecode,omitempty"`
	Constructor                *ContractCall         `json:"constructor,omitempty"`
	GenerateWhitelistedAddress bool                  `json:"generate_whitelisted_address,omitempty"`
	GasLimit                   string                `json:"gas_limit,omitempty"`
	GasPriceLimit              string                `json:"gas_price_limit,omitempty"`
	Comment                    string                `json:"comment,omitempty"`
	ContractType               string                `json:"contract_type,omitempty"`
	FeePayerID                 string                `json:"fee_payer_id,omitempty"`
	FeeLimit                   string                `json:"fee_limit,omitempty"`
	Contract                   GenericCreateContract `json:"contract"`
	TransactionReference       string                `json:"transaction_reference,omitempty"`
	ExternalRequestID          string                `json:"external_request_id,omitempty"`
}
