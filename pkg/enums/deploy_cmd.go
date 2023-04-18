package enums

//go:generate toolkit gen enum DeployCmd
type DeployCmd uint8

const (
	DEPLOY_CMD_UNKNOWN DeployCmd = iota

	_ // Deprecated DEPLOY_CMD__CREATE
	DEPLOY_CMD__START
	DEPLOY_CMD__STOP
	_ // Deprecated DEPLOY_CMD__REMOVE
	DEPLOY_CMD__RESTART
)
