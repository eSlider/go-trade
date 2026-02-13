// Package currency provides embedded fiat and cryptocurrency reference data
// with 170+ fiat currencies and 60+ crypto tokens, loaded from a bundled YAML file.
package currency

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

// Kind classifies a currency as fiat or crypto.
type Kind string

const (
	Fiat   Kind = "fiat"
	Crypto Kind = "crypto"
)

// Currency represents a single fiat or cryptocurrency.
type Currency struct {
	ID          string `yaml:"id" json:"id"`
	Code        string `yaml:"code" json:"code"`
	Description string `yaml:"description" json:"description"`
	Kind        Kind   `yaml:"kind" json:"kind"`
	Country     string `yaml:"country" json:"country,omitempty"`
}

// IsFiat returns true if this is a fiat currency.
func (c Currency) IsFiat() bool { return c.Kind == Fiat }

// IsCrypto returns true if this is a cryptocurrency.
func (c Currency) IsCrypto() bool { return c.Kind == Crypto }

// Unit represents a measurement unit (mass, volume, energy, etc.).
type Unit struct {
	ID          string `yaml:"id" json:"id"`
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	Type        string `yaml:"type" json:"type"`
}

// Units is a list of measurement units.
type Units []Unit

// Currencies is a list of currencies with filter and lookup methods.
type Currencies []Currency

// Cryptos returns all cryptocurrency entries.
func (c Currencies) Cryptos() Currencies {
	var list Currencies
	for _, v := range c {
		if v.IsCrypto() {
			list = append(list, v)
		}
	}
	return list
}

// Fiats returns all fiat currency entries.
func (c Currencies) Fiats() Currencies {
	var list Currencies
	for _, v := range c {
		if v.IsFiat() {
			list = append(list, v)
		}
	}
	return list
}

// Get returns a currency by its code (e.g. "BTC", "USD"), or nil if not found.
func (c Currencies) Get(code string) *Currency {
	for _, v := range c {
		if v.Code == code {
			return &v
		}
	}
	return nil
}

// Provider holds the complete currency and unit reference data.
type Provider struct {
	Currencies Currencies `yaml:"currencies"`
	Units      Units      `yaml:"units"`
}

//go:embed currencies.yml
var fileString []byte

// New creates a new Provider with all embedded currency and unit data loaded.
func New() (*Provider, error) {
	provider := &Provider{}
	err := yaml.Unmarshal(fileString, provider)
	return provider, err
}
