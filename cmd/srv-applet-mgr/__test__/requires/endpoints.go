package requires

import (
	"log"
	"os"
	"time"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/__test__/clients/applet_mgr"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/global"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/client"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/types"
)

func Client(transports ...client.HttpTransport) *applet_mgr.Client {
	c := &client.Client{
		Host: "127.0.0.1",
		Port: uint16(global.ServerMgr.Port),
	}
	c.SetDefault()

	c.Transports = append(c.Transports, transports...)
	return applet_mgr.NewClient(c)
}

func AuthClient(transports ...client.HttpTransport) *applet_mgr.Client {
	return Client(NewAuthPatchRT())
}

func Serve() (stop func()) {
	go kit.Run(apis.RootMgr, global.Server())

	time.Sleep(3 * time.Second)

	return func() {
		p, _ := os.FindProcess(os.Getpid())
		_ = p.Signal(os.Interrupt)
		time.Sleep(3 * time.Second)
	}
}

func DropTempWasmDatabase(projectID *types.SFID) {
	if *projectID == 0 {
		return
	}

	d := types.MustMgrDBExecutorFromContext(global.Context)

	_, err := d.Exec(builder.Expr("DROP DATABASE w3b_" + projectID.String()))
	if err != nil {
		log.Println(err)
	}
	log.Printf("database: %v dropped", projectID)
}
