package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/machinefi/w3bstream/tools/wsctl/cmd/applet"
	"github.com/machinefi/w3bstream/tools/wsctl/cmd/instance"
	"github.com/machinefi/w3bstream/tools/wsctl/cmd/project"
	"github.com/machinefi/w3bstream/tools/wsctl/cmd/publisher"
	"github.com/machinefi/w3bstream/tools/wsctl/cmd/strategy"
	"github.com/stretchr/testify/require"
)

func TestDummyApplet(t *testing.T) {
	var (
		require     = require.New(t)
		projectName = gofakeit.Noun()
		deviceName  = gofakeit.Noun()
		deviceKey   = gofakeit.Noun()
		eventType   = gofakeit.Noun()
		payload     = gofakeit.Sentence(20)
		pwd, _      = os.Getwd()
		filePath    = filepath.Join(pwd, "./testdata/log.wasm")
	)

	_, err := project.Create(_cli, projectName)
	require.NoError(err)

	appID, insID, err := applet.Create(_cli, projectName, filePath)
	require.NoError(err)

	// TODO: enable after the bug is fixed
	// err = instance.Start(_cli, insID)
	// require.NoError(err)

	pubID, pubToken, err := publisher.Create(_cli, projectName, deviceName, deviceKey)
	require.NoError(err)

	stgID, err := strategy.Create(_cli, projectName, appID, eventType, "start")
	require.NoError(err)

	err = publisher.HTTPEvent(_cli, projectName, pubToken, eventType, payload)
	require.NoError(err)

	err = strategy.Delete(_cli, stgID)
	require.NoError(err)

	err = publisher.Delete(_cli, pubID)
	require.NoError(err)

	err = instance.Stop(_cli, insID)
	require.NoError(err)

	err = instance.Delete(_cli, insID)
	require.NoError(err)

	err = applet.Delete(_cli, appID)
	require.NoError(err)

	err = project.Delete(_cli, projectName)
	require.NoError(err)
}
