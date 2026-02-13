# go-trade

[![CI](https://github.com/eSlider/go-trade/actions/workflows/ci.yml/badge.svg)](https://github.com/eSlider/go-trade/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/eslider/go-trade.svg)](https://pkg.go.dev/github.com/eslider/go-trade)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8.svg)](https://go.dev)
[![Latest Release](https://img.shields.io/github/v/tag/eSlider/go-trade?sort=semver&label=release)](https://github.com/eSlider/go-trade/releases)
[![GitHub Stars](https://img.shields.io/github/stars/eSlider/go-trade?style=social)](https://github.com/eSlider/go-trade/stargazers)

Go library providing a **unified data model for aggregating market data across multiple exchanges**. Defines exchange-agnostic types for candles, time-and-sale transactions, order books, orders, instruments, symbols, and currencies — serving as the common schema for normalizing data from Binance, CME, NASDAQ, and other trading venues.

## Architecture

```mermaid
graph TB
    subgraph "Exchange Connectors"
        BIN["Binance"]
        CME["CME Group"]
        NAS["NASDAQ"]
        FX["Forex Brokers"]
        CEX["Other CEX/DEX"]
    end

    subgraph "go-trade — Unified Data Model"
        TAS["TimeAndSale<br/>Atomic trade events"]
        CND["Candle<br/>OHLC + microstructure"]
        OB["OrderBook<br/>Bid/Ask snapshots"]
        ORD["Order<br/>Trade orders"]
        INS["Instrument<br/>Tradeable assets"]
        MKT["Market<br/>Trading pairs"]
        SYM["Symbol<br/>Hierarchical asset tree"]
        CUR["currency.Provider<br/>170+ fiat, 60+ crypto"]
    end

    subgraph "Storage & Analysis"
        DB["TimescaleDB<br/>PostgreSQL"]
        PQ["Parquet Files"]
        DASH["Analytics<br/>Dashboard"]
    end

    BIN -->|"normalize"| TAS
    CME -->|"normalize"| TAS
    NAS -->|"normalize"| TAS
    FX -->|"normalize"| TAS
    CEX -->|"normalize"| TAS

    TAS --> CND
    TAS --> OB

    TAS --> DB
    CND --> DB
    OB --> PQ
    ORD --> DB

    DB --> DASH
    PQ --> DASH

    INS -.->|"classifies"| TAS
    MKT -.->|"identifies"| INS
    SYM -.->|"defines"| MKT
    CUR -.->|"resolves"| SYM
```

## Data Flow

```mermaid
sequenceDiagram
    participant Exchange
    participant Connector
    participant Model as go-trade Model
    participant Store

    Exchange->>Connector: Raw trade event (JSON/WebSocket)
    Connector->>Connector: Parse exchange-specific format
    Connector->>Model: trade.TimeAndSale{Ticker, Price, Volume, Side}
    Model->>Store: Normalized row

    Note over Model: Same struct regardless<br/>of source exchange

    Exchange->>Connector: Raw candle data
    Connector->>Model: trade.Candle{OHLC, Clusters, DeltaLevels}
    Model->>Store: Normalized candle

    Store->>Store: Cross-exchange aggregation
```

## Cross-Exchange Aggregation

```mermaid
graph LR
    subgraph "BTC/USDT Trades"
        B1["Binance<br/>$97,420 × 0.5 BTC"]
        B2["Coinbase<br/>$97,418 × 1.2 BTC"]
        B3["Kraken<br/>$97,422 × 0.8 BTC"]
    end

    subgraph "Unified TimeAndSale"
        T["trade.TimeAndSale<br/>• Ticker: BTCUSDT<br/>• ExchangeID: varies<br/>• Price / Volume / Side"]
    end

    subgraph "Aggregated Candle"
        C["trade.Candle<br/>• OHLC across all venues<br/>• PriceClusters by level<br/>• DeltaLevels (buy vs sell)<br/>• VolumeLevels"]
    end

    B1 --> T
    B2 --> T
    B3 --> T
    T --> C
```

## Symbol Hierarchy

Symbols form a tree that links base assets to their derivatives and stablecoins, enabling cross-exchange normalization:

```mermaid
graph TD
    USD["USD (fiat)"]
    BTC["BTC (crypto)"]
    ETH["ETH (crypto)"]

    USD --> USDT["USDT — stablecoin"]
    USD --> USDC["USDC — stablecoin"]
    USD --> EUR["EUR (fiat)"]
    USD --> GBP["GBP (fiat)"]

    BTC --> BTCFT["BTCFT — futures token"]
    BTC --> BTCM24["BTCM24 — Jun 2024 future"]
    BTC --> WBTC["WBTC — wrapped"]

    ETH --> STETH["stETH — liquid staking"]
    ETH --> ETHM24["ETHM24 — Jun 2024 future"]

    style USD fill:#4CAF50,color:#fff
    style BTC fill:#F7931A,color:#fff
    style ETH fill:#627EEA,color:#fff
```

## Installation

```bash
go get github.com/eslider/go-trade
```

## Features

- **Candle model** — OHLC + ask/bid, price clusters, delta/volume levels, price snake
- **TimeAndSale** — Atomic trade events with exchange ID, data feed provider, aggressor side
- **OrderBook** — Bid/ask snapshots with timestamps
- **Order** — Custom JSON unmarshaling for string-encoded fields from exchange APIs
- **Symbol** — Hierarchical asset tree (fiat/crypto) with parent-child relationships for derivatives
- **Instrument / Market** — Asset classification (spot, future, option, FX) and trading pairs
- **Currency provider** — Embedded 170+ fiat currencies and 60+ crypto tokens with lookup and filtering
- **DateTime / UUID** — JSON-aware wrappers for exchange date formats and UUIDs

## Quick Start

### Candle Analysis

```go
candle := trade.Candle{
    Open: 100.0, High: 110.0, Low: 95.0, Close: 108.0,
    Ask: 108.5, Bid: 107.5,
    TimeOpen:  time.Now().Add(-5 * time.Minute),
    TimeClose: time.Now(),
}

fmt.Println("Bullish:", candle.IsBullish())   // true
fmt.Println("Range:", candle.Range())          // 15.0
fmt.Println("Spread:", candle.Spread())        // 1.0
fmt.Println("Duration:", candle.Duration())    // 5m0s
```

### Time & Sale Recording

```go
event := trade.TimeAndSale{
    Ticker:     "BTCUSDT",
    ExchangeID: 1,  // Binance
    Time:       time.Now(),
    Sale: trade.Sale{
        Price:         97420.50,
        AggressorSide: trade.AggressorBuy,
        Volume:        5,
    },
}
fmt.Printf("%s: %s @ $%.2f (%s)\n",
    event.Ticker, event.Sale.AggressorSide,
    event.Sale.Price, event.Time.Format(time.RFC3339))
```

### Currency Lookup

```go
import "github.com/eslider/go-trade/currency"

provider, _ := currency.New()

btc := provider.Currencies.Get("BTC")
fmt.Println(btc.Description) // "Bitcoin"
fmt.Println(btc.IsCrypto())  // true

fiats := provider.Currencies.Fiats()
fmt.Printf("%d fiat currencies loaded\n", len(fiats))

cryptos := provider.Currencies.Cryptos()
fmt.Printf("%d cryptocurrencies loaded\n", len(cryptos))
```

### Symbol Tree

```go
symbols := trade.Symbols{
    {ID: 1, Type: trade.SymbolFiat, Code: "USD", Name: "US Dollar"},
    {ID: 2, Type: trade.SymbolCrypto, ParentID: 1, Code: "USDT", Name: "Tether"},
    {ID: 3, Type: trade.SymbolCrypto, ParentID: 1, Code: "USDC", Name: "USD Coin"},
    {ID: 14, Type: trade.SymbolCrypto, Code: "BTC", Name: "Bitcoin"},
    {ID: 22, Type: trade.SymbolCrypto, ParentID: 14, Code: "BTCFT", Name: "BTC Futures Token"},
}

// Tree navigation
roots := symbols.Roots()               // USD, BTC
btcDerivs := symbols.Children(14)      // BTCFT
usdStables := symbols.Children(1)      // USDT, USDC

// Lookup
btc := symbols.GetByCode("BTC")
fmt.Println(btc.IsCrypto(), btc.IsRoot()) // true true

// Filter
fiats := symbols.Fiats()               // USD
cryptos := symbols.Cryptos()            // USDT, USDC, BTC, BTCFT
```

### Order Deserialization

```go
// Exchange APIs often return numeric fields as strings
raw := []byte(`{
    "order_id": "b68d69564a79dea4776afa33d1d2fcab",
    "customer_id": "41",
    "order_status": "shipped",
    "order_approved_at": "2018-02-28 10:40:35"
}`)

var order trade.Order
json.Unmarshal(raw, &order)
fmt.Println(order.OrderID)    // b68d6956-4a79-dea4-776a-fa33d1d2fcab
fmt.Println(order.CustomerID) // 41 (int64, not string)
```

## Package Structure

```
go-trade/
├── trade.go               # Package doc
├── candle.go              # Candle (OHLC + microstructure)
├── time_and_sale.go       # Atomic trade events
├── order.go               # Trading orders
├── instrument.go          # Instruments and markets
├── symbol.go              # Hierarchical asset symbols
├── aggressor_side.go      # Buy/sell side enum
├── datetime.go            # Exchange datetime parser
├── uuid.go                # UUID JSON wrapper
├── price_clusters.go      # Volume at price level
├── candle_delta_levels.go # Delta/volume levels
├── trade_test.go          # Unit tests
├── symbol_test.go         # Symbol tests
└── currency/
    ├── currency.go        # Fiat + crypto provider
    ├── currency_test.go   # Currency tests
    └── currencies.yml     # 170+ fiat, 60+ crypto definitions
```

## API Reference

### Core Types

| Type | Description |
|---|---|
| `Candle` | OHLC candlestick with price clusters, delta levels, and volume profile |
| `TimeAndSale` | Atomic trade event: price, volume, side, exchange, timestamp |
| `Order` | Trading order with auto-deserialization from string-encoded JSON |
| `Symbol` | Trading symbol with type (fiat/crypto) and parent-child hierarchy |
| `Symbols` | Collection with `GetByCode`, `GetByID`, `Children`, `Roots`, `Fiats`, `Cryptos` |
| `SymbolType` | Enum: `SymbolFiat`, `SymbolCrypto` |
| `Instrument` | Tradeable asset (spot, future, option, FX) |
| `Market` | Trading pair (FROM/TO symbols) |
| `OrderBook` / `OrderBookEntry` | Bid/ask snapshot at a point in time |
| `Sale` | Core trade data (price, aggressor side, volume) |
| `PriceClusters` | Volume and time distribution at a price level |
| `CandleDeltaLevels` | Min/max delta values within a candle |
| `DateTime` | JSON-compatible datetime for `"2006-01-02 15:04:05"` format |
| `UUID` | JSON-compatible UUID wrapper |

### Candle Methods

| Method | Description |
|---|---|
| `IsBullish()` | Close > Open |
| `IsBearish()` | Close < Open |
| `IsEmpty()` | No data |
| `Spread()` | Ask − Bid |
| `Range()` | High − Low |
| `Duration()` | TimeClose − TimeOpen |

### Symbol Methods

| Method | Description |
|---|---|
| `IsRoot()` | ParentID == 0 (top-level asset) |
| `IsFiat()` | Type == SymbolFiat |
| `IsCrypto()` | Type == SymbolCrypto |

### Symbols Collection

| Method | Description |
|---|---|
| `GetByCode(code)` | Find symbol by ticker code |
| `GetByID(id)` | Find symbol by numeric ID |
| `Children(parentID)` | All symbols with given parent |
| `Roots()` | All top-level symbols (parent == 0) |
| `Fiats()` | Filter fiat symbols |
| `Cryptos()` | Filter crypto symbols |

### Currency Package

| Function | Description |
|---|---|
| `currency.New()` | Load all embedded currency/unit data |
| `Currencies.Get(code)` | Lookup by code ("BTC", "USD") |
| `Currencies.Cryptos()` | Filter crypto only |
| `Currencies.Fiats()` | Filter fiat only |

## Environment

- **Go 1.22+**
- **CI**: GitHub Actions (Go 1.22, 1.23, 1.24) with race detection and golangci-lint

## Related

- [go-ollama](https://github.com/eSlider/go-ollama) — Ollama/Open WebUI API client
- [go-matrix-bot](https://github.com/eSlider/go-matrix-bot) — Matrix bot with AI integration
- [go-onlyoffice](https://github.com/eSlider/go-onlyoffice) — OnlyOffice Project Management API
- [go-gitea-helpers](https://github.com/eSlider/go-gitea-helpers) — Gitea pagination helpers

## License

[MIT](LICENSE)
