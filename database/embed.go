package database

import _ "embed"

//go:embed init.sql
var InitSQL []byte
