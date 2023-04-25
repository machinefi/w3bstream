package enums

//go:generate toolkit gen enum DeployCmd
type DeployCmd uint8

const (
	DEPLOY_CMD_UNKNOWN DeployCmd = iota
	_
	DEPLOY_CMD__START
	DEPLOY_CMD__STOP
	DEPLOY_CMD__REMOVE
	_
)
