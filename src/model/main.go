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

type DBWriteQuery struct {
    Stmt       *sql.Stmt
    Parameters []interface{}
}

var queryCh chan DBWriteQuery
var errCh chan error

func doSQLWrite(query DBWriteQuery) error {
    _, err := query.Stmt.Exec(query.Parameters...)
    if err != nil {
        return err
    }
    return query.Stmt.Close()
}

func asyncSQLWrite() {
    for {
        query := <-queryCh
        errCh <-doSQLWrite(query)
    }
}

func init() {
    _db, err := sql.Open("sqlite3", config.Config.DB.DBFile)
    if err != nil {
        panic(err)
    }
    db = _db

    // Init database
    _, err = db.Exec("CREATE TABLE IF NOT EXISTS " + config.Config.DB.TableUser +
        "(name TEXT, school TEXT, contest_id TEXT, person_id TEXT, language INTEGER)")
    if err != nil {
        panic(err)
    }
    _, err = db.Exec("CREATE TABLE IF NOT EXISTS " + config.Config.DB.TableSubmit +
        "(id TEXT, user TEXT, md5 TEXT, time INTEGER, problem_id INTEGER, confirm INTEGER)")
    if err != nil {
        panic(err)
    }
    _, err = db.Exec("CREATE TABLE IF NOT EXISTS " + config.Config.DB.TableLatestSubmit +
        "(user TEXT, submit_id TEXT, problem_id INTEGER)")
    if err != nil {
        panic(err)
    }

    queryCh = make(chan DBWriteQuery, config.Config.DB.ChannelBuffer)
    errCh = make(chan error, config.Config.DB.ChannelBuffer)
    go asyncSQLWrite()
}
