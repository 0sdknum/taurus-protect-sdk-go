// Package service provides high-level service wrappers for the Taurus-PROTECT API.
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/0sdknum/taurus-protect-sdk-go/internal/openapi"
	"github.com/0sdknum/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/0sdknum/taurus-protect-sdk-go/pkg/protect/model"
)

// WalletService provides wallet management operations.
type WalletService struct {
	api       *openapi.WalletsAPIService
	errMapper *ErrorMapper
}

// NewWalletService creates a new WalletService.
func NewWalletService(client *openapi.APIClient) *WalletService {
	return &WalletService{
		api:       client.WalletsAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetWallet retrieves a wallet by ID.
func (s *WalletService) GetWallet(ctx context.Context, walletID string) (*model.Wallet, error) {
	if walletID == "" {
		return nil, fmt.Errorf("walletID cannot be empty")
	}

	resp, httpResp, err := s.api.WalletServiceGetWalletV2(ctx, walletID).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("wallet not found")
	}

	return mapper.WalletFromDTO(resp.Result), nil
}

// ListWallets retrieves a list of wallets.
func (s *WalletService) ListWallets(ctx context.Context, opts *model.ListWalletsOptions) ([]*model.Wallet, *model.Pagination, error) {
	req := s.api.WalletServiceGetWalletsInfo(ctx)

	if opts != nil {
		if opts.Limit > 0 {
			req = req.Limit(fmt.Sprintf("%d", opts.Limit))
		}
		if opts.Offset > 0 {
			req = req.Offset(fmt.Sprintf("%d", opts.Offset))
		}
		if opts.Currency != "" {
			req = req.Currencies([]string{opts.Currency})
		}
		if opts.Query != "" {
			req = req.Query(opts.Query)
		}
		if opts.ExcludeDisabled {
			req = req.ExcludeDisabled(true)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, nil, s.errMapper.MapError(err, httpResp)
	}

	wallets := mapper.WalletsFromDTO(resp.Result)

	var pagination *model.Pagination
	if resp.TotalItems != nil || resp.Offset != nil {
		pagination = &model.Pagination{}
		if resp.TotalItems != nil {
			if total, err := strconv.ParseInt(*resp.TotalItems, 10, 64); err == nil {
				pagination.TotalItems = total
			}
		}
		if resp.Offset != nil {
			if offset, err := strconv.ParseInt(*resp.Offset, 10, 64); err == nil {
				pagination.Offset = offset
			}
		}
		if opts != nil {
			pagination.Limit = opts.Limit
		}
		pagination.HasMore = pagination.Offset+pagination.Limit < pagination.TotalItems
	}

	return wallets, pagination, nil
}

// CreateWallet creates a new wallet.
func (s *WalletService) CreateWallet(ctx context.Context, req *model.CreateWalletRequest) (*model.Wallet, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}

	createReq := openapi.TgvalidatordCreateWalletRequest{
		Name:     req.Name,
		Currency: &req.Currency,
	}

	if req.Comment != "" {
		createReq.Comment = &req.Comment
	}
	if req.CustomerID != "" {
		createReq.CustomerId = &req.CustomerID
	}
	if req.ExternalWalletID != "" {
		createReq.ExternalWalletId = &req.ExternalWalletID
	}
	if req.VisibilityGroupID != "" {
		createReq.VisibilityGroupID = &req.VisibilityGroupID
	}

	resp, httpResp, err := s.api.WalletServiceCreateWallet(ctx).
		Body(createReq).
		Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("failed to create wallet")
	}

	return mapper.WalletFromCreateDTO(resp.Result), nil
}

// CreateWalletAttribute creates a custom attribute on a wallet.
func (s *WalletService) CreateWalletAttribute(ctx context.Context, walletID, key, value string) error {
	if walletID == "" {
		return fmt.Errorf("walletID cannot be empty")
	}
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	attrReq := openapi.TgvalidatordCreateWalletAttributeRequest{}
	attrReq.SetKey(key)
	attrReq.SetValue(value)

	body := openapi.WalletServiceCreateWalletAttributesBody{
		Attributes: []openapi.TgvalidatordCreateWalletAttributeRequest{attrReq},
	}

	_, httpResp, err := s.api.WalletServiceCreateWalletAttributes(ctx, walletID).
		Body(body).
		Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// DeleteWalletAttribute deletes a custom attribute from a wallet.
func (s *WalletService) DeleteWalletAttribute(ctx context.Context, walletID, attributeID string) error {
	if walletID == "" {
		return fmt.Errorf("walletID cannot be empty")
	}
	if attributeID == "" {
		return fmt.Errorf("attributeID cannot be empty")
	}

	_, httpResp, err := s.api.WalletServiceDeleteWalletAttribute(ctx, walletID, attributeID).
		Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// GetWalletBalanceHistory retrieves balance history for a wallet.
// intervalHours specifies the interval between balance snapshots in hours.
func (s *WalletService) GetWalletBalanceHistory(ctx context.Context, walletID string, intervalHours int) ([]*model.BalanceHistoryPoint, error) {
	if walletID == "" {
		return nil, fmt.Errorf("walletID cannot be empty")
	}

	req := s.api.WalletServiceGetWalletBalanceHistory(ctx, walletID)
	if intervalHours > 0 {
		req = req.IntervalHours(fmt.Sprintf("%d", intervalHours))
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.BalanceHistoryPointsFromDTO(resp.Result), nil
}

// GetWalletTokens retrieves token balances for a wallet.
// limit specifies the maximum number of tokens to return.
func (s *WalletService) GetWalletTokens(ctx context.Context, walletID string, limit int) ([]*model.AssetBalance, error) {
	if walletID == "" {
		return nil, fmt.Errorf("walletID cannot be empty")
	}

	req := s.api.WalletServiceGetWalletTokens(ctx, walletID)
	if limit > 0 {
		req = req.Limit(fmt.Sprintf("%d", limit))
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.AssetBalancesFromDTO(resp.Balances), nil
}

// ErrorMapper maps OpenAPI errors to domain errors.
type ErrorMapper struct{}

// NewErrorMapper creates a new ErrorMapper.
func NewErrorMapper() *ErrorMapper {
	return &ErrorMapper{}
}

// MapError converts an OpenAPI error to a domain error.
func (m *ErrorMapper) MapError(err error, resp *http.Response) error {
	if err == nil {
		return nil
	}

	// Try to extract OpenAPI error details
	if openAPIErr, ok := err.(*openapi.GenericOpenAPIError); ok {
		return mapOpenAPIError(openAPIErr, resp)
	}

	return err
}

func mapOpenAPIError(err *openapi.GenericOpenAPIError, resp *http.Response) error {
	if resp == nil {
		return err
	}

	code := resp.StatusCode
	responseBody := strings.TrimSpace(string(err.Body()))
	message, errorCode := parseOpenAPIErrorBody(err.Body())
	if message == "" {
		message = err.Error()
	}

	apiError := &APIError{
		Code:         code,
		StatusCode:   code,
		Message:      message,
		ErrorCode:    errorCode,
		ResponseBody: responseBody,
		Err:          err,
	}

	switch {
	case code == 400:
		apiError.Description = "Bad Request"
	case code == 401:
		apiError.Description = "Unauthorized"
	case code == 403:
		apiError.Description = "Forbidden"
	case code == 404:
		apiError.Description = "Not Found"
	case code == 429:
		apiError.Description = "Rate Limited"
		apiError.RetryAfter = parseRetryAfter(resp.Header.Get("Retry-After"), time.Now())
	case code >= 500:
		apiError.Description = "Server Error"
	}

	return apiError
}

func parseOpenAPIErrorBody(body []byte) (string, string) {
	var payload map[string]any
	if len(body) == 0 || json.Unmarshal(body, &payload) != nil {
		return "", ""
	}

	message := firstString(payload, "message", "detail", "description", "error_description")
	errorCode := firstScalar(payload, "errorCode", "error_code", "code")

	if nested, ok := payload["error"].(map[string]any); ok {
		if message == "" {
			message = firstString(nested, "message", "detail", "description")
		}
		if errorCode == "" {
			errorCode = firstScalar(nested, "errorCode", "error_code", "code")
		}
	} else if message == "" {
		message, _ = payload["error"].(string)
	}

	return message, errorCode
}

func firstString(payload map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := payload[key].(string); ok && value != "" {
			return value
		}
	}
	return ""
}

func firstScalar(payload map[string]any, keys ...string) string {
	for _, key := range keys {
		switch value := payload[key].(type) {
		case string:
			if value != "" {
				return value
			}
		case float64:
			return strconv.FormatFloat(value, 'f', -1, 64)
		}
	}
	return ""
}

func parseRetryAfter(value string, now time.Time) time.Duration {
	seconds, err := strconv.ParseInt(value, 10, 64)
	if err == nil && seconds > 0 {
		return time.Duration(seconds) * time.Second
	}

	retryAt, err := http.ParseTime(value)
	if err != nil || !retryAt.After(now) {
		return 0
	}
	return retryAt.Sub(now)
}

// APIError is the API error type returned by high-level services.
type APIError = model.APIError
