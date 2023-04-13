package device

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-gw-device/types"
	basetypes "github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
)

func GetDeviceByID(ctx context.Context, id basetypes.SFID) (*models.Device, error) {
	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)
	m := &models.Device{
		RelDevice: models.RelDevice{DeviceID: id},
	}

	_, l = l.Start(ctx, "GetDeviceByID")
	defer l.End()

	if err := m.FetchByDeviceID(d); err != nil {
		l.Error(err)
		return nil, status.CheckDatabaseError(err, "GetDeviceByID")
	}

	return m, nil
}
