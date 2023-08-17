package postgres

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
)

func CreateUserIfNotExists(d sqlx.DBExecutor, usename, passwd string) (error, bool) {
	exists := ptrx.Ptr(false)
	err := d.QueryAndScan(
		builder.Expr(fmt.Sprintf("SELECT 1 FROM pg_catalog.pg_user AS u WHERE u.usename='%s'", usename)),
		exists,
	)
	if err != nil && !sqlx.DBErr(err).IsNotFound() {
		return err, false
	}
	if *exists {
		// DO NOT update password otherwise connection will be lost, if this is current user
		return nil, false
	}

	_, err = d.Exec(builder.Expr("CREATE USER " + usename + " WITH PASSWORD '" + passwd + "'"))
	if err != nil {
		return err, false
	}
	return nil, true
}

func DropUser(d sqlx.DBExecutor, usename string) error {
	currentUser, err := CurrentUser(d)
	if currentUser == usename {
		return ErrDropCurrentUser
	}

	_, err = d.Exec(builder.Expr("DROP USER " + usename))
	return err
}

func CurrentUser(d sqlx.DBExecutor) (string, error) {
	usename := ""
	err := d.QueryAndScan(builder.Expr("SELECT USER"), &usename)
	if err != nil {
		return "", err
	}
	return usename, nil
}

func GrantAllPrivileges(d sqlx.DBExecutor, on PrivilegeDomain, onName, usename string) error {
	_, err := d.Exec(builder.Expr(fmt.Sprintf(
		"GRANT ALL PRIVILEGES ON %s %s TO %s",
		on, onName, usename,
	)))
	return err
}

func AlterUserConnectionLimit(d sqlx.DBExecutor, usename string, lmt int) error {
	_, err := d.Exec(builder.Expr(fmt.Sprintf(
		"ALTER USER %s CONNECTION LIMIT %d", usename, lmt,
	)))
	return err
}

var (
	ErrDropCurrentUser = errors.New("cannot drop current user")
)

type PrivilegeDomain string

const (
	PrivilegeDomainDatabase PrivilegeDomain = "DATABASE"
	PrivilegeDomainTABLE    PrivilegeDomain = "TABLE"
)
