package test

//go:generate mockgen -source=../depends/kit/sqlx/database.go -destination=./mock_depends_kit_sqlx/database.go -package=mock_sqlx
//go:generate mockgen -source=../depends/conf/storage/storage.go -destination=./mock_depends_conf_storage/storage.go -package=mock_conf_storage
//go:generate mockgen -source=./mock_base_types/interfaces.go -destination=./mock_base_types/mock_base_types.go -package=mock_base_types

// TODO gen modules monkey patches for each function
// TODO gen models monkey patches for each method
