package service

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"strconv"

	"github.com/0sdknum/taurus-protect-sdk-go/internal/openapi"
	"github.com/0sdknum/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/0sdknum/taurus-protect-sdk-go/pkg/protect/helper"
	"github.com/0sdknum/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/0sdknum/taurus-protect-sdk-go/pkg/protect/model"
)

// WhitelistedAddressService provides whitelisted address management operations.
type WhitelistedAddressService struct {
	api       *openapi.AddressWhitelistingAPIService
	errMapper *ErrorMapper
	verifier  *helper.WhitelistedAddressVerifier
}

// WhitelistedAddressServiceConfig holds configuration for the WhitelistedAddressService.
type WhitelistedAddressServiceConfig struct {
	// SuperAdminKeys are the public keys used to verify governance rules signatures.
	SuperAdminKeys []*ecdsa.PublicKey
	// MinValidSignatures is the minimum number of valid SuperAdmin signatures required.
	MinValidSignatures int
}

// NewWhitelistedAddressServiceWithVerification creates a new WhitelistedAddressService with
// integrity verification enabled. When SuperAdmin keys are provided, all retrieved addresses
// will be cryptographically verified before being returned.
func NewWhitelistedAddressServiceWithVerification(
	client *openapi.APIClient,
	config *WhitelistedAddressServiceConfig,
) *WhitelistedAddressService {
	svc := &WhitelistedAddressService{
		api:       client.AddressWhitelistingAPI,
		errMapper: NewErrorMapper(),
	}

	// Create verifier if SuperAdmin keys are configured
	if config != nil && len(config.SuperAdminKeys) > 0 {
		svc.verifier = helper.NewWhitelistedAddressVerifier(
			config.SuperAdminKeys,
			config.MinValidSignatures,
		)
	}

	return svc
}

// GetWhitelistedAddress retrieves a whitelisted address by ID.
// If verification is enabled (SuperAdmin keys configured), the address integrity
// will be cryptographically verified before being returned.
func (s *WhitelistedAddressService) GetWhitelistedAddress(ctx context.Context, id string) (*model.WhitelistedAddress, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	resp, httpResp, err := s.api.WhitelistServiceGetWhitelistedAddress(ctx, id).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("whitelisted address not found")
	}

	// SECURITY: Verify hash BEFORE calling mapper to prevent extraction of unverified data.
	// The mapper extracts security-critical fields from PayloadAsString, so we must verify
	// that PayloadAsString has not been tampered with before the mapper runs.
	if err := s.verifyMetadataHashFromDTO(resp.Result); err != nil {
		return nil, err
	}

	// Now safe to call mapper (extracts from verified PayloadAsString)
	addr := mapper.WhitelistedAddressFromDTO(resp.Result)

	// Full verification (rules container signatures, whitelist signatures) — always enforced
	if addr != nil {
		if err := s.verifyAddress(addr); err != nil {
			return nil, err
		}
	}

	return addr, nil
}

// CreateWhitelistedAddress creates a whitelisted address and returns its ID.
func (s *WhitelistedAddressService) CreateWhitelistedAddress(ctx context.Context, request *model.CreateWhitelistedAddressRequest) (string, error) {
	if request == nil {
		return "", fmt.Errorf("request cannot be nil")
	}
	if request.Address == "" {
		return "", fmt.Errorf("address is required")
	}
	if request.Label == "" {
		return "", fmt.Errorf("label is required")
	}
	if request.Blockchain == "" {
		return "", fmt.Errorf("blockchain is required")
	}

	createRequest := openapi.TgvalidatordCreateWhitelistedAddressRequest{
		Address:    request.Address,
		Label:      request.Label,
		Blockchain: request.Blockchain,
	}
	if request.Memo != "" {
		createRequest.Memo = new(request.Memo)
	}
	if request.ExchangeAccountID != "" {
		createRequest.ExchangeAccountId = new(request.ExchangeAccountID)
	}
	if request.CustomerID != "" {
		createRequest.CustomerId = new(request.CustomerID)
	}
	if len(request.LinkedInternalAddressIDs) > 0 {
		createRequest.LinkedInternalAddressIds = append([]string(nil), request.LinkedInternalAddressIDs...)
	}
	if request.AddressType != "" {
		createRequest.AddressType = new(request.AddressType)
	}
	if request.ContractType != "" {
		createRequest.ContractType = new(request.ContractType)
	}
	if len(request.LinkedWalletIDs) > 0 {
		createRequest.LinkedWalletIds = append([]string(nil), request.LinkedWalletIDs...)
	}
	if request.Network != "" {
		createRequest.Network = new(request.Network)
	}
	if request.VisibilityGroupID != "" {
		createRequest.VisibilityGroupID = new(request.VisibilityGroupID)
	}
	if request.Currency != "" {
		createRequest.Currency = new(request.Currency)
	}

	response, httpResponse, err := s.api.WhitelistServiceCreateWhitelistedAddress(ctx).
		Body(createRequest).
		Execute()
	if err != nil {
		return "", s.errMapper.MapError(err, httpResponse)
	}
	if response == nil || response.Result == nil || response.Result.Id == nil || *response.Result.Id == "" {
		return "", fmt.Errorf("failed to create whitelisted address: no ID returned")
	}

	return *response.Result.Id, nil
}

// ListWhitelistedAddresses retrieves a list of whitelisted addresses.
// If verification is enabled (SuperAdmin keys configured), each address integrity
// will be cryptographically verified before being returned. If any address fails
// verification, an error is returned.
func (s *WhitelistedAddressService) ListWhitelistedAddresses(ctx context.Context, opts *model.ListWhitelistedAddressesOptions) ([]*model.WhitelistedAddress, *model.Pagination, error) {
	req := s.api.WhitelistServiceGetWhitelistedAddresses(ctx)

	if opts != nil {
		if opts.Limit > 0 {
			req = req.Limit(fmt.Sprintf("%d", opts.Limit))
		}
		if opts.Offset > 0 {
			req = req.Offset(fmt.Sprintf("%d", opts.Offset))
		}
		if opts.Blockchain != "" {
			req = req.Blockchain(opts.Blockchain)
		}
		if opts.Network != "" {
			req = req.Network(opts.Network)
		}
		if opts.Currency != "" {
			req = req.Currency(opts.Currency)
		}
		if opts.Query != "" {
			req = req.Query(opts.Query)
		}
		if opts.AddressType != "" {
			req = req.AddressType(opts.AddressType)
		}
		if len(opts.IDs) > 0 {
			req = req.Ids(opts.IDs)
		}
		if len(opts.Addresses) > 0 {
			req = req.Addresses(opts.Addresses)
		}
		if opts.IncludeForApproval {
			req = req.IncludeForApproval(true)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, nil, s.errMapper.MapError(err, httpResp)
	}

	// SECURITY: Verify hash BEFORE calling mapper for each address.
	// This ensures we don't extract fields from unverified PayloadAsString data.
	for i, dto := range resp.Result {
		if err := s.verifyMetadataHashFromDTO(&dto); err != nil {
			addrID := ""
			if dto.Id != nil {
				addrID = *dto.Id
			}
			return nil, nil, fmt.Errorf("hash verification failed for address %s (index %d): %w", addrID, i, err)
		}
	}

	// Now safe to call mapper (extracts from verified PayloadAsString)
	addresses := mapper.WhitelistedAddressesFromDTO(resp.Result)

	// Full verification (rules container signatures, whitelist signatures) — always enforced
	for _, addr := range addresses {
		if addr != nil {
			if err := s.verifyAddress(addr); err != nil {
				return nil, nil, fmt.Errorf("verification failed for address %s: %w", addr.ID, err)
			}
		}
	}

	var pagination *model.Pagination
	if resp.TotalItems != nil {
		pagination = &model.Pagination{}
		if total, parseErr := strconv.ParseInt(*resp.TotalItems, 10, 64); parseErr == nil {
			pagination.TotalItems = total
		}
		if opts != nil {
			pagination.Limit = opts.Limit
			pagination.Offset = opts.Offset
			pagination.HasMore = pagination.Offset+pagination.Limit < pagination.TotalItems
		}
	}

	return addresses, pagination, nil
}

// verifyMetadataHashFromDTO verifies the metadata hash before mapping fields from payloadAsString.
func (s *WhitelistedAddressService) verifyMetadataHashFromDTO(dto *openapi.TgvalidatordSignedWhitelistedAddressEnvelope) error {
	if dto == nil || dto.Metadata == nil {
		return nil // No metadata to verify
	}

	payloadAsString := dto.Metadata.PayloadAsString
	providedHash := dto.Metadata.Hash

	// Skip verification if no payload data (e.g., addresses pending approval)
	if payloadAsString == nil || *payloadAsString == "" {
		return nil
	}

	// If there's payloadAsString but no hash, that's an error
	if providedHash == nil || *providedHash == "" {
		return &model.IntegrityError{Message: "metadata hash is missing but payloadAsString is present"}
	}

	computedHash := crypto.CalculateHexHash(*payloadAsString)
	if !helper.ConstantTimeCompare(computedHash, *providedHash) {
		return &model.IntegrityError{Message: "metadata hash verification failed"}
	}

	return nil
}

// verifyAddress performs the 6-step integrity verification on a whitelisted address.
// Returns nil if verification passes, or an error describing the failure.
func (s *WhitelistedAddressService) verifyAddress(addr *model.WhitelistedAddress) error {
	if s.verifier == nil {
		return &model.IntegrityError{Message: "verification is required but no verifier is configured"}
	}

	// Verification is enabled but required data is missing — this is an error.
	// An attacker could strip verification data to bypass checks.
	if addr.Metadata == nil || addr.RulesContainer == "" || addr.SignedAddress == nil {
		return &model.IntegrityError{Message: "verification enabled but required data missing"}
	}

	_, err := s.verifier.VerifyWhitelistedAddress(
		addr,
		mapper.RulesContainerFromBase64,
		mapper.UserSignaturesFromBase64,
	)
	return err
}

// GetWhitelistedAddressEnvelope retrieves a whitelisted address envelope by ID and performs
// the complete 6-step verification flow. The envelope contains both the verified whitelisted
// address and the decoded rules container.
//
// This method requires verification to be enabled (SuperAdmin keys must be configured).
// If verification is not configured, an error is returned.
//
// The 6-step verification flow:
// 1. Verify metadata hash (SHA-256 of payloadAsString == metadata.hash)
// 2. Verify rules container signatures (SuperAdmin signatures)
// 3. Decode rules container (base64 -> protobuf -> model)
// 4. Verify hash coverage (metadata.hash in at least one signature.hashes)
// 5. Verify whitelist signatures (user signatures meet governance thresholds)
// 6. Parse WhitelistedAddress from verified payload
func (s *WhitelistedAddressService) GetWhitelistedAddressEnvelope(
	ctx context.Context,
	id string,
) (*model.WhitelistedAddressEnvelope, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	if s.verifier == nil {
		return nil, &model.IntegrityError{
			Message: "verification is required for GetWhitelistedAddressEnvelope but no SuperAdmin keys are configured",
		}
	}

	resp, httpResp, err := s.api.WhitelistServiceGetWhitelistedAddress(ctx, id).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("whitelisted address not found")
	}

	// Create envelope from DTO
	envelope := mapper.WhitelistedAddressEnvelopeFromDTO(resp.Result)

	// Initialize and verify the envelope
	if err := s.initializeEnvelope(envelope); err != nil {
		return nil, err
	}

	return envelope, nil
}

// initializeEnvelope performs the 6-step verification and populates the verified fields.
func (s *WhitelistedAddressService) initializeEnvelope(envelope *model.WhitelistedAddressEnvelope) error {
	if envelope == nil {
		return fmt.Errorf("envelope cannot be nil")
	}

	// Check required data for verification
	if envelope.Metadata == nil {
		return &model.IntegrityError{Message: "metadata is required for verification"}
	}
	if envelope.RulesContainer == "" {
		return &model.IntegrityError{Message: "rules container is required for verification"}
	}
	if envelope.SignedAddress == nil {
		return &model.IntegrityError{Message: "signed address is required for verification"}
	}

	// Create a temporary WhitelistedAddress for verification
	tempAddr := &model.WhitelistedAddress{
		ID:                      envelope.ID,
		Blockchain:              envelope.Blockchain,
		Network:                 envelope.Network,
		Metadata:                envelope.Metadata,
		SignedAddress:           envelope.SignedAddress,
		RulesContainer:          envelope.RulesContainer,
		RulesSignatures:         envelope.RulesSignatures,
		LinkedInternalAddresses: envelope.LinkedInternalAddresses,
		LinkedWallets:           envelope.LinkedWallets,
	}

	// Perform the 6-step verification
	result, err := s.verifier.VerifyWhitelistedAddress(
		tempAddr,
		mapper.RulesContainerFromBase64,
		mapper.UserSignaturesFromBase64,
	)
	if err != nil {
		return err
	}

	// Set verified data on envelope
	envelope.SetVerified(result.VerifiedAddress, result.RulesContainer)

	return nil
}
