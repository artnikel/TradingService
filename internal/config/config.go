// Package config with environment variables
package config

// Variables is a struct with environment variables
type Variables struct {
	TokenSignature      string `env:"TOKEN_SIGNATURE"`
	PostgresConnTrading string `env:"POSTGRES_CONN_TRADING"`
	CompanyShares       string `env:"COMPANY_SHARES"`
}
