/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package model

import (
    "OIUP-Backend/config"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
    db, _ = sql.Open("sqlite3", config.Config.DB.DBFile)
}
