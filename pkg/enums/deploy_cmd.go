package enums

//go:generate toolkit gen enum DeployCmd
type DeployCmd uint8

const (
	DEPLOY_CMD_UNKNOWN DeployCmd = iota
	_
	DEPLOY_CMD__START  // start wasm vm
	DEPLOY_CMD__HUNGUP // stop wasm vm
	_
	_
)
