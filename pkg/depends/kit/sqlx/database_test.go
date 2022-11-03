package sqlx_test

import (
	"context"
	"database/sql/driver"
	"testing"
	"time"

	// . "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var connectors map[string]driver.Connector

func init() {
	connectors = make(map[string]driver.Connector)
	// TODO add other database connector for testing
	// mysqlConnector = &mysql.MysqlConnector{
	// 	Host:  "root@tcp(0.0.0.0:3306)",
	// 	Extra: "charset=utf8mb4&parseTime=true&interpolateParams=true&autocommit=true&loc=Local",
	// }

	// postgres
	/*
		{
			ep := postgrestestutil.Endpoint
			connectors["postgres"] = &postgres.Connector{
				Extra:      "sslmode=disable",
				Extensions: []string{"postgis"},
				Host: fmt.Sprintf(
					"postgresql://%s:%s@127.0.0.1:5432",
					ep.Master.Username, ep.Master.Password,
				),
			}
		}
	*/
}

func Background() context.Context {
	return log.WithLogger(context.Background(), log.Std())
}

type OperateTime struct {
	CreatedAt types.Datetime `db:"f_created_at,default=CURRENT_TIMESTAMP,onupdate=CURRENT_TIMESTAMP"`
	UpdatedAt int64          `db:"f_updated_at,default='0'"`
}

type Gender int

const (
	GenderMale Gender = iota + 1
	GenderFemale
)

func (Gender) EnumType() string {
	return "Gender"
}

func (Gender) Enums() map[int][]string {
	return map[int][]string{
		int(GenderMale):   {"male", "男"},
		int(GenderFemale): {"female", "女"},
	}
}

func (g Gender) String() string {
	switch g {
	case GenderMale:
		return "male"
	case GenderFemale:
		return "female"
	}
	return ""
}

type User struct {
	ID       uint64 `db:"f_id,autoincrement"`
	Name     string `db:"f_name,size=255,default=''"`
	Nickname string `db:"f_nickname,size=255,default=''"`
	Username string `db:"f_username,default=''"`
	Gender   Gender `db:"f_gender,default='0'"`

	OperateTime
}

func (user *User) Comments() map[string]string {
	return map[string]string{
		"Name": "姓名",
	}
}

func (user *User) TableName() string {
	return "t_user"
}

func (user *User) PrimaryKey() []string {
	return []string{"ID"}
}

func (user *User) Indexes() builder.Indexes {
	return builder.Indexes{
		"i_nickname": {"Nickname"},
	}
}

func (user *User) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"i_name": {"Name"},
	}
}

type User2 struct {
	ID       uint64 `db:"f_id,autoincrement"`
	Nickname string `db:"f_nickname,size=255,default=''"`
	Gender   Gender `db:"f_gender,default='0'"`
	Name     string `db:"f_name,deprecated=f_real_name"`
	RealName string `db:"f_real_name,size=255,default=''"`
	Age      int32  `db:"f_age,default='0'"`
	Username string `db:"f_username,deprecated"`
}

func (user *User2) TableName() string {
	return "t_user"
}

func (user *User2) PrimaryKey() []string {
	return []string{"ID"}
}

func (user *User2) Indexes() builder.Indexes {
	return builder.Indexes{
		"i_nickname": {"Nickname"},
	}
}

func (user *User2) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"i_name": {"RealName"},
	}
}

func DISABLED_TestMigrate(t *testing.T) {
	/*
		os.Setenv("PROJECT_FEATURE", "test1")
		defer func() {
			os.Remove("PROJECT_FEATURE")
		}()

		dbTest := sqlx.NewDatabase("test_for_migrate")

		for name, connector := range connectors {
			t.Run(name, func(t *testing.T) {
				for _, schema := range []string{"import", "public", "backup"} {
					dbTest.Tables.Range(func(table *builder.Table, idx int) {
						db := dbTest.OpenDB(connector).WithSchema(schema)
						_, _ = db.Exec(db.Dialect().DropTable(table))
					})

					t.Run("CreateTable_"+schema, func(t *testing.T) {
						dbTest.Register(&User{})
						db := dbTest.OpenDB(connector).WithSchema(schema)

						t.Run("First", func(t *testing.T) {
							err := migration.Migrate(db, nil)
							NewWithT(t).Expect(err).To(BeNil())
						})

						t.Run("Again", func(t *testing.T) {
							_ = migration.Migrate(db, os.Stdout)
							err := migration.Migrate(db, nil)
							NewWithT(t).Expect(err).To(BeNil())
						})
					})

					t.Run("NoMigrate_"+schema, func(t *testing.T) {
						dbTest.Register(&User{})
						db := dbTest.OpenDB(connector).WithSchema(schema)
						err := migration.Migrate(db, nil)
						NewWithT(t).Expect(err).To(BeNil())

						t.Run("migrate to user2", func(t *testing.T) {
							dbTest.Register(&User2{})
							db := dbTest.OpenDB(connector).WithSchema(schema)
							err := migration.Migrate(db, nil)
							NewWithT(t).Expect(err).To(BeNil())
						})

						t.Run("migrate to user2 again", func(t *testing.T) {
							dbTest.Register(&User2{})
							db := dbTest.OpenDB(connector).WithSchema(schema)
							err := migration.Migrate(db, nil)
							NewWithT(t).Expect(err).To(BeNil())
						})
					})

					t.Run("MigrateToUser_"+schema, func(t *testing.T) {
						db := dbTest.OpenDB(connector).WithSchema(schema)
						err := migration.Migrate(db, os.Stdout)
						NewWithT(t).Expect(err).To(BeNil())
						err = migration.Migrate(db, nil)
						NewWithT(t).Expect(err).To(BeNil())
					})

					dbTest.Tables.Range(func(table *builder.Table, idx int) {
						db := dbTest.OpenDB(connector).WithSchema(schema)
						_, _ = db.Exec(db.Dialect().DropTable(table))
					})
				}
			})
		}
	*/
}

func DISABLED_TestCRUD(t *testing.T) {
	/*
		dbTest := sqlx.NewDatabase("test_crud")

		for name, connector := range connectors {
			t.Run(name, func(t *testing.T) {
				d := dbTest.OpenDB(connector)
				db := d.WithContext(metax.ContextWithMeta(d.Context(), metax.ParseMeta("_id=11111")))
				userTable := dbTest.Register(&User{})
				err := migration.Migrate(db, nil)

				NewWithT(t).Expect(err).To(BeNil())

				t.Run("InsertSingle", func(t *testing.T) {
					user := User{
						Name:   uuid.New().String(),
						Gender: GenderMale,
					}

					t.Run("Canceled", func(t *testing.T) {
						ctx, cancel := context.WithCancel(Background())
						db2 := db.WithContext(ctx)

						go func() {
							time.Sleep(5 * time.Millisecond)
							cancel()
						}()

						err := sqlx.NewTasks(db2).
							With(
								func(db sqlx.DBExecutor) error {
									_, err := db.Exec(sqlx.InsertToDB(db, &user, nil))
									return err
								},
								func(db sqlx.DBExecutor) error {
									time.Sleep(10 * time.Millisecond)
									return nil
								},
							).
							Do()

						NewWithT(t).Expect(err).NotTo(BeNil())
					})

					_, err := db.Exec(sqlx.InsertToDB(db, &user, nil))
					NewWithT(t).Expect(err).To(BeNil())

					t.Run("Update", func(t *testing.T) {
						user.Gender = GenderFemale
						_, err := db.Exec(
							builder.Update(dbTest.T(&user)).
								Set(sqlx.AsAssignments(db, &user)...).
								Where(
									userTable.ColByFieldName("Name").Eq(user.Name),
								),
						)
						NewWithT(t).Expect(err).To(BeNil())
					})
					t.Run("Select", func(t *testing.T) {
						userForSelect := User{}
						err := db.QueryAndScan(
							builder.Select(nil).From(
								userTable,
								builder.Where(userTable.ColByFieldName("Name").Eq(user.Name)),
								builder.Comment("FindUser"),
							),
							&userForSelect)

						NewWithT(t).Expect(err).To(BeNil())

						NewWithT(t).Expect(user.Name).To(Equal(userForSelect.Name))
						NewWithT(t).Expect(user.Gender).To(Equal(userForSelect.Gender))
					})
					t.Run("Conflict", func(t *testing.T) {
						_, err := db.Exec(sqlx.InsertToDB(db, &user, nil))
						NewWithT(t).Expect(sqlx.DBErr(err).IsConflict()).To(BeTrue())
					})
				})
				db.(*sqlx.DB).Tables.Range(func(table *builder.Table, idx int) {
					_, err := db.Exec(db.Dialect().DropTable(table))
					NewWithT(t).Expect(err).To(BeNil())
				})
			})
		}
	*/
}

type UserSet map[string]*User

func (UserSet) New() interface{} {
	return &User{}
}

func (u UserSet) Next(v interface{}) error {
	user := v.(*User)
	u[user.Name] = user
	time.Sleep(500 * time.Microsecond)
	return nil
}

func DISABLED_TestSelect(t *testing.T) {
	/*
		dbTest := sqlx.NewDatabase("test_for_select")

		for name, connector := range connectors {
			t.Run(name, func(t *testing.T) {
				db := dbTest.OpenDB(connector)
				table := dbTest.Register(&User{})

				db.Tables.Range(func(t *builder.Table, idx int) {
					_, _ = db.Exec(db.Dialect().DropTable(t))
				})

				err := migration.Migrate(db, nil)
				NewWithT(t).Expect(err).To(BeNil())

				{
					columns := table.MustColsByFieldNames("Name", "Gender")
					values := make([]interface{}, 0)

					for i := 0; i < 1000; i++ {
						values = append(values, uuid.New().String(), GenderMale)
					}

					_, err := db.Exec(
						builder.Insert().Into(table).Values(columns, values...),
					)
					NewWithT(t).Expect(err).To(BeNil())
				}

				t.Run("SelectToSlice", func(t *testing.T) {
					users := make([]User, 0)
					err := db.QueryAndScan(
						builder.Select(nil).From(
							table,
							builder.Where(
								table.ColByFieldName("Gender").Eq(GenderMale),
							),
						),
						&users,
					)
					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(users).To(HaveLen(1000))
				})

				t.Run("SelectToIter", func(t *testing.T) {
					userSet := UserSet{}
					err := db.QueryAndScan(
						builder.Select(nil).From(
							table,
							builder.Where(
								table.ColByFieldName("Gender").Eq(GenderMale),
							),
						),
						userSet,
					)
					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(userSet).To(HaveLen(1000))
				})

				t.Run("NotFound", func(t *testing.T) {
					user := User{}
					err := db.QueryAndScan(
						builder.Select(nil).From(
							table,
							builder.Where(
								table.ColByFieldName("ID").Eq(1001),
							),
						),
						&user,
					)
					NewWithT(t).Expect(sqlx.DBErr(err).IsNotFound()).To(BeTrue())
				})

				t.Run("SelectCount", func(t *testing.T) {
					count := 0
					err := db.QueryAndScan(
						builder.Select(builder.Count()).From(table),
						&count,
					)
					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(count).To(Equal(1000))
				})

				t.Run("Canceled", func(t *testing.T) {
					ctx, cancel := context.WithCancel(Background())
					db2 := db.WithContext(ctx)

					go func() {
						time.Sleep(3 * time.Millisecond)
						cancel()
					}()

					userSet := UserSet{}
					err := db2.QueryAndScan(
						builder.Select(nil).From(
							table,
							builder.Where(
								table.ColByFieldName("Gender").Eq(GenderMale),
							),
						),
						userSet,
					)
					NewWithT(t).Expect(err).NotTo(BeNil())
				})

				db.Tables.Range(func(tab *builder.Table, idx int) {
					_, _ = db.Exec(db.Dialect().DropTable(tab))
				})
			})
		}
	*/
}
