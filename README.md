# Taurus-PROTECT Go SDK Fork

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

An independent Go-only fork of the official [Taurus-PROTECT SDK](https://github.com/taurushq-io/taurus-protect-sdk). It provides a typed client for the Taurus-PROTECT API, TPV1 authentication, integrity verification, custody operations, transaction management, and Taurus Network services.

> [!IMPORTANT]
> This repository is not an official Taurus SA distribution and is not endorsed or supported by Taurus SA. For the official multi-language SDK, use the [upstream repository](https://github.com/taurushq-io/taurus-protect-sdk).

## Fork Status

- Extracted from the upstream `taurus-protect-sdk-go` module without the Java, Python, and TypeScript SDKs.
- Distributed as the standalone module `github.com/0sdknum/taurus-protect-sdk-go`.
- Preserves the upstream public API while adding public high-level wrappers for Taurus-PROTECT operations that were previously available only through the generated internal client.
- The `v0.2.0` release adds outgoing contract calls and deployments, whitelisted-address creation, and richer API error details.
- Generated OpenAPI and protobuf code remains under `internal/` and must not be edited manually.

Fork-specific changes should remain generic and suitable for contribution back to upstream. See the repository releases for the exact upstream base and changes included in each published version.

## Documentation

| Document                                                                     | Description                                      |
| ---------------------------------------------------------------------------- | ------------------------------------------------ |
| [Key Concepts](docs/CONCEPTS.md)                                             | Go model types, exceptions, and domain concepts  |
| [SDK Overview](docs/SDK_OVERVIEW.md)                                         | Architecture, packages, and design patterns      |
| [Authentication](docs/AUTHENTICATION.md)                                     | TPV1 authentication and cryptographic operations |
| [Services Reference](docs/SERVICES.md)                                       | Complete API documentation for all 43 services   |
| [Usage Examples](docs/USAGE_EXAMPLES.md)                                     | Go code examples and common patterns             |
| [Whitelisted Address Verification](docs/WHITELISTED_ADDRESS_VERIFICATION.md) | 6-step verification flow                         |

## Quick Start

### Prerequisites

- Go 1.26 or higher
- Go modules enabled

### Installation

Install the latest tagged version using Go modules:

```bash
go get github.com/0sdknum/taurus-protect-sdk-go@latest
```

Applications import the fork directly. `go get` records the resolved version in `go.mod`; a `replace` directive is not required. Transitive dependencies are declared by this module and must not be installed individually.

### Dependencies

The SDK requires the following packages (installed automatically):

| Package                                      | Purpose                         |
| -------------------------------------------- | ------------------------------- |
| `github.com/google/uuid`                     | TPV1 request nonces             |
| `github.com/grpc-ecosystem/grpc-gateway/v2` | Generated protobuf JSON support |
| `golang.org/x/net`                           | Generated HTTP helpers          |
| `google.golang.org/genproto`                 | Generated Google API types      |
| `google.golang.org/protobuf`                 | Protocol buffer support         |
| Standard library                             | HTTP client, crypto, encoding   |

### Client Initialization

Credentials, SuperAdmin public keys, and a positive minimum signature count are required when creating a client:

```go
package main

import (
    "context"
    "log"

    "github.com/0sdknum/taurus-protect-sdk-go/pkg/protect"
)

const (
    apiHost   = "https://your-taurus-protect-host"
    apiKey    = "your-api-key"
    apiSecret = "your-hex-encoded-api-secret"
)

var superAdminKeys = []string{
    `-----BEGIN PUBLIC KEY-----
REPLACE_WITH_BASE64_DER_PUBLIC_KEY
-----END PUBLIC KEY-----`,
}

func main() {
    client, err := protect.NewClient(
        apiHost,
        protect.WithCredentials(apiKey, apiSecret),
        protect.WithSuperAdminKeysPEM(superAdminKeys),
        protect.WithMinValidSignatures(1),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    wallets, _, err := client.Wallets().ListWallets(context.Background(), nil)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("wallets: %d", len(wallets))
}
```

Replace every placeholder locally. Never commit API credentials or real public/private key material.

See [Authentication](docs/AUTHENTICATION.md) for more initialization options.

## Services

The SDK provides 43 services organized into core services and the TaurusNetwork namespace.

### Core Services (38 services)

| Service                       | Access                          | Purpose                                      |
| ----------------------------- | ------------------------------- | -------------------------------------------- |
| `WalletService`               | `client.Wallets()`              | Wallet creation, retrieval, balance history  |
| `AddressService`              | `client.Addresses()`            | Address management, proof of reserve         |
| `RequestService`              | `client.Requests()`             | Transfers, contract calls/deployments, approvals |
| `TransactionService`          | `client.Transactions()`         | Transaction queries and export               |
| `BalanceService`              | `client.Balances()`             | Balance queries across assets                |
| `CurrencyService`             | `client.Currencies()`           | Currency and blockchain information          |
| `GovernanceRuleService`       | `client.GovernanceRules()`      | Governance rules with signature verification |
| `WhitelistedAddressService`   | `client.WhitelistedAddresses()` | Address creation, lookup, and verification   |
| `WhitelistedAssetService`     | `client.WhitelistedAssets()`    | Asset/contract whitelisting                  |
| `AuditService`                | `client.Audits()`               | Audit log queries                            |
| `ChangeService`               | `client.Changes()`              | Configuration change tracking                |
| `FeeService`                  | `client.Fees()`                 | Transaction fee information                  |
| `PriceService`                | `client.Prices()`               | Price data and history                       |
| `AirGapService`               | `client.AirGap()`               | Air-gap signing operations                   |
| `StakingService`              | `client.Staking()`              | Multi-chain staking information              |
| `WhitelistedContractService`  | `client.WhitelistedContracts()` | Smart contract whitelisting                  |
| `BusinessRuleService`         | `client.BusinessRules()`        | Business rule management                     |
| `ReservationService`          | `client.Reservations()`         | Balance reservations                         |
| `MultiFactorSignatureService` | `client.MultiFactorSignature()` | Multi-factor signature operations            |
| `UserService`                 | `client.Users()`                | User management                              |
| `GroupService`                | `client.Groups()`               | User group management                        |
| `VisibilityGroupService`      | `client.VisibilityGroups()`     | Visibility group management                  |
| `ConfigService`               | `client.Config()`               | System configuration                         |
| `WebhookService`              | `client.Webhooks()`             | Webhook management                           |
| `WebhookCallService`          | `client.WebhookCalls()`         | Webhook call history                         |
| `TagService`                  | `client.Tags()`                 | Tag management                               |
| `AssetService`                | `client.Assets()`               | Asset information                            |
| `ActionService`               | `client.Actions()`              | Action management                            |
| `BlockchainService`           | `client.Blockchains()`          | Blockchain information                       |
| `ExchangeService`             | `client.Exchanges()`            | Exchange integration                         |
| `FiatService`                 | `client.Fiat()`                 | Fiat currency operations                     |
| `FeePayerService`             | `client.FeePayers()`            | Fee payer management                         |
| `HealthService`               | `client.Health()`               | API health checks                            |
| `JobService`                  | `client.Jobs()`                 | Background job management                    |
| `ScoreService`                | `client.Scores()`               | Risk scoring operations                      |
| `StatisticsService`           | `client.Statistics()`           | Platform statistics                          |
| `TokenMetadataService`        | `client.TokenMetadata()`        | Token metadata information                   |
| `UserDeviceService`           | `client.UserDevices()`          | User device management                       |

### TaurusNetwork Services (5 services)

| Service              | Access                                  | Purpose                       |
| -------------------- | --------------------------------------- | ----------------------------- |
| `ParticipantService` | `client.TaurusNetwork().Participants()` | Participant management        |
| `PledgeService`      | `client.TaurusNetwork().Pledges()`      | Pledge lifecycle operations   |
| `LendingService`     | `client.TaurusNetwork().Lending()`      | Lending offers and agreements |
| `SettlementService`  | `client.TaurusNetwork().Settlements()`  | Settlement operations         |
| `SharingService`     | `client.TaurusNetwork().Sharing()`      | Address and asset sharing     |

See [Services Reference](docs/SERVICES.md) for complete API documentation.

## Basic Usage

### List Wallets

```go
ctx := context.Background()

// List wallets with pagination
wallets, pagination, err := client.Wallets().ListWallets(ctx, &model.ListWalletsOptions{
    Limit: 50,
})
if err != nil {
    return err
}

for _, wallet := range wallets {
    fmt.Printf("%s: %s (%s/%s)\n", wallet.Name, wallet.Currency, wallet.Blockchain, wallet.Network)
}

if pagination != nil {
    fmt.Printf("Total wallets: %d\n", pagination.TotalItems)
}
```

### Create a Wallet

```go
wallet, err := client.Wallets().CreateWallet(ctx, &model.CreateWalletRequest{
    Name:     "Trading Wallet",
    Currency: "ETH",
})
if err != nil {
    return err
}
fmt.Printf("Created wallet ID: %s\n", wallet.ID)
```

### Create a CMTA20 Contract Call

Contract calls are real state-changing Taurus-PROTECT requests. Amounts and ABI
integer arguments use decimal base-unit strings. CMTA20 mutable operations use
the legacy `Method` form; do not also set `Call` in the same request.

```go
request, err := client.Requests().CreateOutgoingCallContractRequest(ctx, &model.CreateOutgoingCallContractRequest{
    FromAddressID:          "your-source-address-id",
    ToWhitelistedAddressID: "your-contract-whitelisted-address-id",
    ContractType:           "CMTA20",
    Amount:                 "0",
    ExternalRequestID:      "your-idempotency-key",
    Method: model.ContractCall{
        FunctionSignature: "mint(address,uint256)",
        Arguments: []model.ContractArgument{
            {
                Name:  "to",
                Type:  "address",
                Value: &model.ContractArgumentValue{Primitive: "your-recipient-address"},
            },
            {
                Name:  "amount",
                Type:  "uint256",
                Value: &model.ContractArgumentValue{Primitive: "1"},
            },
        },
    },
})
if err != nil {
    return err
}
fmt.Printf("Created request ID: %s\n", request.ID)
```

### Create a Whitelisted Address

```go
id, err := client.WhitelistedAddresses().CreateWhitelistedAddress(ctx, &model.CreateWhitelistedAddressRequest{
    Address:    "your-external-address",
    Blockchain: "ETH",
    Network:    "your-network",
    Label:      "your-label",
})
if err != nil {
    return err
}
fmt.Printf("Created whitelisted address ID: %s\n", id)
```

### Approve Transaction Requests

```go
import (
    "crypto/ecdsa"
    "github.com/0sdknum/taurus-protect-sdk-go/pkg/protect/crypto"
)

// Load your private key (takes a PEM string, not bytes)
privateKey, err := crypto.DecodePrivateKeyPEM(pemString)
if err != nil {
    return err
}

// Get requests pending approval
requests, _, err := client.Requests().ListRequestsForApproval(ctx, &model.ListRequestsOptions{
    Limit: 10,
})
if err != nil {
    return err
}

if len(requests) > 0 {
    // Approve with ECDSA signature
    signedCount, err := client.Requests().ApproveRequests(ctx, requests, privateKey)
    if err != nil {
        return err
    }
    fmt.Printf("Approved %d request(s)\n", signedCount)
}
```

### TaurusNetwork Operations

```go
// Get my participant info
me, err := client.TaurusNetwork().Participants().GetMyParticipant(ctx)
if err != nil {
    return err
}
fmt.Printf("Participant: %s\n", me.Name)

// List pledges
pledges, cursor, err := client.TaurusNetwork().Pledges().ListPledges(ctx, &taurusnetwork.ListPledgesOptions{
    Limit: 10,
})
if err != nil {
    return err
}
for _, pledge := range pledges {
    fmt.Printf("Pledge %s: %s\n", pledge.ID, pledge.Status)
}

// List shared addresses
addresses, _, err := client.TaurusNetwork().Sharing().ListSharedAddresses(ctx, nil)
if err != nil {
    return err
}
```

### Error Handling

```go
import (
    "errors"
    "github.com/0sdknum/taurus-protect-sdk-go/pkg/protect"
)

wallet, err := client.Wallets().GetWallet(ctx, "999999")
if err != nil {
    if errors.Is(err, protect.ErrNotFound) {
        fmt.Println("Wallet not found")
    } else if errors.Is(err, protect.ErrRateLimit) {
        if apiErr, ok := protect.IsAPIError(err); ok && apiErr.RetryAfter > 0 {
            fmt.Printf("Rate limited, retry after %v\n", apiErr.RetryAfter)
            time.Sleep(apiErr.RetryAfter)
        }
    } else if protect.IsIntegrityError(err) {
        // Security error - DO NOT retry
        fmt.Printf("Integrity verification failed: %v\n", err)
    } else if apiErr, ok := protect.IsAPIError(err); ok {
        if apiErr.IsRetryable() {
            fmt.Printf("Retryable error: %s\n", apiErr.Message)
        } else {
            fmt.Printf("Non-retryable error: %s\n", apiErr.Message)
        }
    }
    return err
}
```

See [Usage Examples](docs/USAGE_EXAMPLES.md) for comprehensive examples.

## Build Commands

```bash
# Build and test (default) - USE THIS TO VERIFY CHANGES
./build.sh

# Build only
./build.sh build

# Run unit tests only
./build.sh unit

# Run linter (requires golangci-lint)
./build.sh lint

# Code generation requires the upstream monorepo resources directory
./build.sh generate

# Clean build artifacts
./build.sh clean
```

## Development

### Running Tests

Cross-SDK crypto vectors are stored in package-local `testdata`, so the complete unit-test gate is self-contained and works from a downloaded Go module.

```bash
# Run all unit tests (includes verbose output and coverage)
./build.sh unit

# Run a single test by name pattern
./build.sh unit-one TestMapWallet

# Run all integration tests (requires API access, see below)
./build.sh integration

# Run a single integration test by name pattern
./build.sh integration-one TestListWallets
```

### Integration Tests

Integration tests require environment configuration:

```bash
export PROTECT_INTEGRATION_TEST=true
export PROTECT_API_HOST="https://your-api-host.com"
export PROTECT_API_KEY="your-api-key"
export PROTECT_API_SECRET="your-hex-encoded-secret"

./build.sh integration
```

## License

The original code is copyright (c) 2026 Taurus SA. The original code and fork modifications are distributed under the [MIT license](./LICENSE).
