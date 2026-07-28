package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/0sdknum/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewTransactionService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestTransactionService_ListTransactions_ForwardsSynchronizationFilters(t *testing.T) {
	from := time.Date(2026, time.July, 1, 10, 11, 12, 0, time.UTC)
	to := from.Add(2 * time.Hour)
	server := newTransactionListServer(t, map[string]string{
		"limit":           "25",
		"offset":          "50",
		"currency":        "PPEUR",
		"direction":       "incoming",
		"blockchain":      "ETH",
		"query":           "request-1",
		"from":            from.Format(time.RFC3339),
		"to":              to.Format(time.RFC3339),
		"fromBlockNumber": "100",
		"toBlockNumber":   "200",
	})

	transactionService := NewTransactionService(newServiceTestAPIClient(server))
	_, _, err := transactionService.ListTransactions(context.Background(), &model.ListTransactionsOptions{
		Limit:           25,
		Offset:          50,
		Currency:        "PPEUR",
		Direction:       "incoming",
		Blockchain:      "ETH",
		Query:           "request-1",
		From:            &from,
		To:              &to,
		FromBlockNumber: "100",
		ToBlockNumber:   "200",
	})
	if err != nil {
		t.Fatalf("ListTransactions() error = %v", err)
	}
}

func TestTransactionService_ListTransactionsByAddress_ForwardsSynchronizationFilters(t *testing.T) {
	from := time.Date(2026, time.July, 1, 10, 11, 12, 0, time.UTC)
	to := from.Add(2 * time.Hour)
	server := newTransactionListServer(t, map[string]string{
		"address":         "0xabc",
		"limit":           "25",
		"offset":          "50",
		"currency":        "PPEUR",
		"direction":       "outgoing",
		"blockchain":      "ETH",
		"from":            from.Format(time.RFC3339),
		"to":              to.Format(time.RFC3339),
		"fromBlockNumber": "100",
		"toBlockNumber":   "200",
	})

	transactionService := NewTransactionService(newServiceTestAPIClient(server))
	_, _, err := transactionService.ListTransactionsByAddress(
		context.Background(),
		"0xabc",
		&model.ListTransactionsByAddressOptions{
			Limit:           25,
			Offset:          50,
			Currency:        "PPEUR",
			Direction:       "outgoing",
			Blockchain:      "ETH",
			From:            &from,
			To:              &to,
			FromBlockNumber: "100",
			ToBlockNumber:   "200",
		},
	)
	if err != nil {
		t.Fatalf("ListTransactionsByAddress() error = %v", err)
	}
}

func newTransactionListServer(t *testing.T, expectedQuery map[string]string) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet || request.URL.Path != "/api/rest/v1/transactions" {
			t.Errorf("request = %s %s", request.Method, request.URL.Path)
		}
		for key, expected := range expectedQuery {
			if actual := request.URL.Query().Get(key); actual != expected {
				t.Errorf("query %s = %q, want %q", key, actual, expected)
			}
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
	return server
}
