package postgres_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	confpostgres "github.com/machinefi/w3bstream/pkg/depends/conf/postgres"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/driver/postgres"
)

var d *confpostgres.Endpoint

func init() {
	d = &confpostgres.Endpoint{
		Master: types.Endpoint{
			Hostname: "localhost",
			Port:     15432,
			Base:     "postgres",
			Username: "root",
			Password: "test_passwd",
			Param:    map[string][]string{"sslmode": {"disable"}},
		},
		Database: sqlx.NewDatabase("postgres"),
	}
	d.SetDefault()

	if err := d.Init(); err != nil {
		panic(err)
	}
}

func TestCurrentUser(t *testing.T) {
	usename, err := postgres.CurrentUser(d)
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(usename).To(Equal("root"))
}

func TestCreateUserIfNotExists(t *testing.T) {
	var (
		usename = "notexist"
		passwd  = "nopasswd"
	)

	err, isCreated := postgres.CreateUserIfNotExists(d, usename, passwd)
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(isCreated).To(BeTrue())

	err, isCreated = postgres.CreateUserIfNotExists(d, usename, passwd)
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(isCreated).To(BeFalse())

	err = postgres.DropUser(d, usename)
	NewWithT(t).Expect(err).To(BeNil())
}

func TestDropUser(t *testing.T) {
	err := postgres.DropUser(d, "root")
	NewWithT(t).Expect(err).To(Equal(postgres.ErrDropCurrentUser))

	var (
		usename = "notexist"
		passwd  = "nopasswd"
	)
	err, isCreated := postgres.CreateUserIfNotExists(d, usename, passwd)
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(isCreated).To(BeTrue())

	err = postgres.DropUser(d, usename)
	NewWithT(t).Expect(err).To(BeNil())
}

func TestGrantAllPrivileges(t *testing.T) {
	err := postgres.GrantAllPrivileges(d, postgres.PrivilegeDomainDatabase, "postgres", "root")
	NewWithT(t).Expect(err).To(BeNil())
}

func TestAlterUserConnectionLimit(t *testing.T) {
	err := postgres.AlterUserConnectionLimit(d, "root", 100)
	NewWithT(t).Expect(err).To(BeNil())
}
