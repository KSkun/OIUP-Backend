/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package model

import (
    "OIUP-Backend/config"
    "encoding/json"
    "errors"
    "strconv"
    "time"
)

const (
    SubmitUnconfirmed = 0
    SubmitConfirmed   = 1
)

type SubmitInfo struct {
    ID        string    `json:"id"`
    User      string    `json:"user"`
    MD5       string    `json:"md5"`
    MD5Set    []MD5Info `json:"-"`
    Time      int       `json:"time"`
    ProblemID int       `json:"problem_id"`
    Confirm   int       `json:"confirm"`
}

type MD5Info struct {
    MD5    string `json:"md5"`
    TestID int    `json:"test_id"`
}

func AddCodeSubmit(submitID string, user string, md5 string, time time.Time, problemID int) error {
    addCodeSubmitQuery, err := db.Prepare("INSERT INTO " + config.Config.DB.TableSubmit +
        " VALUES (?, ?, ?, ?, ?, ?)")
    if err != nil {
        return err
    }
    queryCh <- DBWriteQuery{
        Stmt:       addCodeSubmitQuery,
        Parameters: []interface{}{submitID, user, md5, time.Unix(), problemID, SubmitUnconfirmed},
    }
    err = <-errCh
    return err
}

func AddOutputSubmit(submitID string, user string, md5Set []MD5Info, time time.Time, problemID int) error {
    md5JSON, err := json.Marshal(md5Set)
    if err != nil {
        return err
    }

    addOutputSubmitQuery, err := db.Prepare("INSERT INTO " + config.Config.DB.TableSubmit +
        " VALUES (?, ?, ?, ?, ?, ?)")
    if err != nil {
        return err
    }
    queryCh <- DBWriteQuery{
        Stmt:       addOutputSubmitQuery,
        Parameters: []interface{}{submitID, user, md5JSON, time.Unix(), problemID, SubmitUnconfirmed},
    }
    err = <-errCh
    return err
}

func GetSubmit(submitID string) (SubmitInfo, bool, error) {
    getSubmitQuery, err := db.Prepare("SELECT * FROM " + config.Config.DB.TableSubmit +
        " WHERE id = ?")
    if err != nil {
        return SubmitInfo{}, false, err
    }
    defer getSubmitQuery.Close()
    rows, err := getSubmitQuery.Query(submitID)
    if err != nil {
        return SubmitInfo{}, false, err
    }
    defer rows.Close()

    var submit SubmitInfo
    if !rows.Next() {
        return submit, false, nil
    }
    err = rows.Scan(&submit.ID, &submit.User, &submit.MD5, &submit.Time, &submit.ProblemID, &submit.Confirm)

    // MD5Info array unmarshal
    problemConfig, found := config.GetProblemConfig(submit.ProblemID)
    if !found {
        return SubmitInfo{}, false, errors.New("problem with id " + strconv.Itoa(submit.ProblemID) + " not found")
    }
    if problemConfig.Type == config.ProblemAnswer {
        err = json.Unmarshal([]byte(submit.MD5), &submit.MD5Set)
        if err != nil {
            return submit, false, err
        }
    }

    return submit, true, nil
}

type LatestSubmitInfo struct {
    User      string
    SubmitID  string
    ProblemID int
}

func GetLatestSubmit(user string, problemID int) (SubmitInfo, bool, error) {
    getLatestSubmitQuery, err := db.Prepare("SELECT * FROM " + config.Config.DB.TableLatestSubmit +
        " WHERE user = ? AND problem_id = ?")
    if err != nil {
        return SubmitInfo{}, false, err
    }
    defer getLatestSubmitQuery.Close()
    rows, err := getLatestSubmitQuery.Query(user, problemID)
    if err != nil {
        return SubmitInfo{}, false, err
    }
    defer rows.Close()

    var latestSubmit LatestSubmitInfo
    if !rows.Next() {
        return SubmitInfo{}, false, nil
    }
    err = rows.Scan(&latestSubmit.User, &latestSubmit.SubmitID, &latestSubmit.ProblemID)

    submit, found, err := GetSubmit(latestSubmit.SubmitID)
    if err != nil {
        return submit, false, err
    }
    if !found {
        return submit, false, errors.New("wrong latest_submit record")
    }
    return submit, true, nil
}

func ConfirmSubmit(submitID string) error {
    submit, found, err := GetSubmit(submitID)
    if err != nil {
        return err
    }
    if !found {
        return errors.New("submit with id " + submitID + " not found")
    }
    /*
       Note of SQL Error: database is locked
       Causes:
         1) 'rows' in other queries remain open after operation
         2) DataGrip and DB Browser connections keep alive
       The 2 causes locks the db and UPDATE event meets a timed-out error.
       Links:
         - https://www.jianshu.com/p/54a76cb84bf5
         - https://blog.csdn.net/LOVETEDA/article/details/82690498
     */
    confirmSubmitQuery, err := db.Prepare("UPDATE " + config.Config.DB.TableSubmit +
        " SET confirm = ? WHERE id = ?")
    if err != nil {
        return err
    }
    queryCh <- DBWriteQuery{
        Stmt:       confirmSubmitQuery,
        Parameters: []interface{}{SubmitConfirmed, submitID},
    }
    err = <-errCh
    if err != nil {
        return err
    }

    deleteLatestSubmitQuery, err := db.Prepare("DELETE FROM " + config.Config.DB.TableLatestSubmit +
        " WHERE user = ? AND problem_id = ?")
    if err != nil {
        return err
    }
    queryCh <- DBWriteQuery{
        Stmt:       deleteLatestSubmitQuery,
        Parameters: []interface{}{submit.User, submit.ProblemID},
    }
    err = <-errCh
    if err != nil {
        return err
    }

    addLatestSubmitQuery, err := db.Prepare("INSERT INTO " + config.Config.DB.TableLatestSubmit +
        " VALUES (?, ?, ?)")
    if err != nil {
        return err
    }
    queryCh <- DBWriteQuery{
        Stmt:       addLatestSubmitQuery,
        Parameters: []interface{}{submit.User, submitID, submit.ProblemID},
    }
    err = <-errCh
    return err
}

func SearchSubmit(filters map[string]interface{}, page int) ([]SubmitInfo, int, error) {
    conditionsStr := getSQLConditionsStr(filters)
    values := make([]interface{}, 0)
    for key, value := range filters {
        if key == "problem_id" {
            values = append(values, value)
            continue
        }
        values = append(values, value.(string) + "%")
    }

    var count int
    countSubmitQueryStr := "SELECT COUNT(*) FROM " + config.Config.DB.TableSubmit
    if len(conditionsStr) > 0 {
        countSubmitQueryStr += " WHERE " + conditionsStr
    }
    countSubmitQuery, err := db.Prepare(countSubmitQueryStr)
    if err != nil {
        return nil, 0, err
    }
    defer countSubmitQuery.Close()
    countRows, err := countSubmitQuery.Query(values...)
    if err != nil {
        return nil, 0, err
    }
    defer countRows.Close()
    countRows.Next()
    err = countRows.Scan(&count)
    if err != nil {
        return nil, 0, err
    }

    searchSubmitQueryStr := "SELECT * FROM " + config.Config.DB.TableSubmit
    if len(conditionsStr) > 0 {
        searchSubmitQueryStr += " WHERE " + conditionsStr
    }
    searchSubmitQueryStr += " ORDER BY user LIMIT " +
        strconv.Itoa(config.Config.DB.RecordsPerPage * page) + " OFFSET " +
        strconv.Itoa(config.Config.DB.RecordsPerPage * (page - 1))
    searchSubmitQuery, err := db.Prepare(searchSubmitQueryStr)
    if err != nil {
        return nil, 0, err
    }
    defer searchSubmitQuery.Close()
    rows, err := searchSubmitQuery.Query(values...)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    submits := make([]SubmitInfo, 0)
    for rows.Next() {
        var submit SubmitInfo
        err = rows.Scan(&submit.ID, &submit.User, &submit.MD5, &submit.Time, &submit.ProblemID, &submit.Confirm)
        if err != nil {
            return nil, 0, err
        }
        submits = append(submits, submit)
    }
    return submits, count, nil
}
