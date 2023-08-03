package enums

//go:generate toolkit gen enum ChainName
type ChainName string

const (
	CHAIN_NAME_IOTEX_MAINNET    ChainName = "iotex-mainnet"
	CHAIN_NAME_IOTEX_TESTNET    ChainName = "iotex-testnet"
	CHAIN_NAME_ETHEREUM_MAINNET ChainName = "ethereum-mainnet"
	CHAIN_NAME_GOERLI           ChainName = "goerli"
	CHAIN_NAME_POLYGON_MAINNET  ChainName = "polygon-mainnet"
	CHAIN_NAME_MUMBAI           ChainName = "mumbai"
	CHAIN_NAME_SOLANA_MAINNET   ChainName = "solana-mainnet"
	CHAIN_NAME_SOLANA_TESTNET   ChainName = "solana-testnet"
	CHAIN_NAME_SOLANA_DEVNET    ChainName = "solana-devnet"
)
