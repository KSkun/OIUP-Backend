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
    _db, err := sql.Open("sqlite3", config.Config.DB.DBFile)
    if err != nil {
        panic(err)
    }
    db = _db

    // Init database
    _, err = db.Exec("CREATE TABLE IF NOT EXISTS user(name TEXT, school TEXT, contest_id TEXT, person_id TEXT, language INTEGER)")
    if err != nil {
        panic(err)
    }
    _, err = db.Exec("CREATE TABLE IF NOT EXISTS submit(id TEXT, user TEXT, md5 TEXT, time INTEGER, problem_id INTEGER, confirm INTEGER)")
    if err != nil {
        panic(err)
    }
    _, err = db.Exec("CREATE TABLE IF NOT EXISTS latest_submit(user TEXT, submit_id TEXT, problem_id INTEGER)")
    if err != nil {
        panic(err)
    }

    writeCh = make(chan DBWriteQuery)
    errCh = make(chan error)
    go doSQLWriteQuery()
}
