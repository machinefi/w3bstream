package enums

//go:generate toolkit gen enum MonitorCmd
type MonitorCmd uint8

const (
	MONITOR_CMD_UNKNOWN MonitorCmd = iota
	MONITOR_CMD__START
	MONITOR_CMD__PAUSE
)
