/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package model

import (
    "OIUP-Backend/config"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "time"
)

var db *sql.DB

type DBWriteQuery struct {
    Stmt       *sql.Stmt
    Parameters []interface{}
}

var writeCh chan DBWriteQuery
var errCh chan error

func doSQLWriteQuery() {
    for {
        select {
        case query := <-writeCh:
            _, err := query.Stmt.Exec(query.Parameters...)
            errCh <-err
            query.Stmt.Close()
        default:
            time.Sleep(time.Duration(config.Config.DB.WriteCheckGap) * time.Millisecond)
        }
    }
}

func init() {
    db, _ = sql.Open("sqlite3", config.Config.DB.DBFile)

    writeCh = make(chan DBWriteQuery)
    errCh = make(chan error)
    go doSQLWriteQuery()
}
