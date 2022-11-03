package postgres_test

import (
	"testing"
	// . "github.com/onsi/gomega"
)

func DISABLED_TestPostgres(t *testing.T) {
	/*
		var (
			pg        = postgrestestutil.Endpoint
			masterURL = []byte(postgrestestutil.Endpoint.Master.String())
			slaveURL  = []byte(postgrestestutil.Endpoint.Slave.String())
		)

		NewWithT(t).Expect(pg.Master.UnmarshalText(masterURL)).To(BeNil())
		NewWithT(t).Expect(pg.Slave.UnmarshalText(slaveURL)).To(BeNil())

		pg.SetDefault()
		pg.Init()

		{
			row, err := pg.QueryContext(context.Background(), "SELECT 1")
			NewWithT(t).Expect(err).To(BeNil())
			_ = row.Close()
		}

		{
			row, err := postgres.SwitchSlave(pg).QueryContext(context.Background(), "SELECT 1")
			NewWithT(t).Expect(err).To(BeNil())
			_ = row.Close()
		}

		NewWithT(t).Expect(pg.UseSlave()).NotTo(Equal(pg.DB))
		NewWithT(t).Expect(pg.LivenessCheck()).To(
			Equal(map[string]string{
				pg.Master.Host(): "ok",
				pg.Slave.Host():  "ok",
			}),
		)
	*/
}
