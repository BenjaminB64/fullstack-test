package database

import _ "embed"

//go:embed init_db.sql
var InitDBSQL string
