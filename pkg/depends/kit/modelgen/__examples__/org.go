package example

// @def primary ID

// Org describe organization information
type Org struct {
	ID   uint64 `db:"f_id,autoincrement"`
	Name string `db:"f_name,default=''"`
	// @rel User.ID
	// User relation...
	UserID string `db:"f_user_id"`
}
