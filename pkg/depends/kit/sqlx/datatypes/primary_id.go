package datatypes

type PrimaryID struct {
	ID uint64 `db:"f_id,autoincrement" json:"-"`
}
