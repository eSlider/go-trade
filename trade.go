// Package trade provides a unified data model for aggregating market data
// across multiple exchanges and data feed providers.
//
// It defines exchange-agnostic types for candles, time-and-sale transactions,
// order books, orders, price clusters, and market instruments. These models
// serve as the common schema for normalizing data from Binance, CME, NASDAQ,
// and other exchanges into a single queryable format.
//
// The currency sub-package provides embedded fiat and cryptocurrency
// reference data with 170+ fiat currencies and 60+ crypto tokens.
//
// # Architecture
//
// Data flows from exchange-specific connectors through these shared types:
//
//	Exchange A ──┐
//	Exchange B ──┼── Connector → trade.TimeAndSale / trade.Candle → Storage
//	Exchange C ──┘
//
// This decouples storage and analysis from exchange-specific wire formats.
package trade
