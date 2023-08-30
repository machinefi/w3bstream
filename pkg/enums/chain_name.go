package enums

//go:generate toolkit gen enum ChainName
type ChainName string

const (
	IOTEX_MAINNET       ChainName = "iotex-mainnet"
	IOTEX_TESTNET       ChainName = "iotex-testnet"
	ETHEREUM_MAINNET    ChainName = "ethereum-mainnet"
	GOERLI              ChainName = "goerli"
	POLYGON_MAINNET     ChainName = "polygon-mainnet"
	MUMBAI              ChainName = "mumbai"
	SOLANA_MAINNET_BETA ChainName = "solana-mainnet-beta"
	SOLANA_TESTNET      ChainName = "solana-testnet"
	SOLANA_DEVNET       ChainName = "solana-devnet"
	ARBITRUM_ONE        ChainName = "arbitrum-one"
	ARBITRUM_GOERLI     ChainName = "arbitrum-goerli"
)
