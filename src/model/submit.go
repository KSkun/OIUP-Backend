/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package model

import (
    "OIUP-Backend/config"
    "encoding/json"
    "errors"
    "fmt"
    "strconv"
    "time"
)

const (
    SubmitUnconfirmed = 0
    SubmitConfirmed   = 1
)

type SubmitInfo struct {
    ID        string     `json:"id"`
    User      string     `json:"user"`
    Hash      string     `json:"hash"`
    HashSet   []HashInfo `json:"-"`
    Time      int        `json:"time"`
    ProblemID int        `json:"problem_id"`
    Confirm   int        `json:"confirm"`
}

type HashInfo struct {
    Hash   string `json:"hash"`
    TestID int    `json:"test_id"`
}

func AddCodeSubmit(submitID string, user string, hash string, time time.Time, problemID int) error {
    addCodeSubmitQuery, err := db.Prepare("INSERT INTO " + config.Config.DB.TableSubmit +
        " VALUES (?, ?, ?, ?, ?, ?)")
    if err != nil {
        return err
    }
    queryCh <- DBWriteQuery{
        Stmt:       addCodeSubmitQuery,
        Parameters: []interface{}{submitID, user, hash, time.Unix(), problemID, SubmitUnconfirmed},
    }
    err = <-errCh
    return err
}

func AddOutputSubmit(submitID string, user string, hashSet []HashInfo, time time.Time, problemID int) error {
    hashJSON, err := json.Marshal(hashSet)
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
        Parameters: []interface{}{submitID, user, hashJSON, time.Unix(), problemID, SubmitUnconfirmed},
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
    err = rows.Scan(&submit.ID, &submit.User, &submit.Hash, &submit.Time, &submit.ProblemID, &submit.Confirm)

    // HashInfo array unmarshal
    problemConfig, found := config.GetProblemConfig(submit.ProblemID)
    if !found {
        return SubmitInfo{}, false, errors.New("找不到编号为 " + strconv.Itoa(submit.ProblemID) + " 的题目！")
    }
    if problemConfig.Type == config.ProblemAnswer {
        err = json.Unmarshal([]byte(submit.Hash), &submit.HashSet)
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
        return submit, false, errors.New("latest_submit 记录错误！")
    }
    return submit, true, nil
}

func ConfirmSubmit(submitID string) error {
    submit, found, err := GetSubmit(submitID)
    if err != nil {
        return err
    }
    if !found {
        return errors.New("找不到编号为 " + submitID + " 的提交！")
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
    if err != nil {
        return err
    }

    fmt.Println(submit.User + ": " + time.Unix(int64(submit.Time), 0).Format("2006-01-02 15:04") +
        " 题目 " + strconv.Itoa(submit.ProblemID) + " 提交 " + submit.ID)
    return nil
}

func SearchSubmit(filters map[string]interface{}) ([]SubmitInfo, error) {
    conditionsStr := getSQLConditionsStr(filters)
    values := make([]interface{}, 0)
    for key, value := range filters {
        if key == "problem_id" {
            values = append(values, value)
            continue
        }
        values = append(values, value.(string) + "%")
    }

    searchSubmitQueryStr := "SELECT * FROM " + config.Config.DB.TableSubmit
    if len(conditionsStr) > 0 {
        searchSubmitQueryStr += " WHERE " + conditionsStr
    }
    searchSubmitQuery, err := db.Prepare(searchSubmitQueryStr)
    if err != nil {
        return nil, err
    }
    defer searchSubmitQuery.Close()
    rows, err := searchSubmitQuery.Query(values...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    submits := make([]SubmitInfo, 0)
    for rows.Next() {
        var submit SubmitInfo
        err = rows.Scan(&submit.ID, &submit.User, &submit.Hash, &submit.Time, &submit.ProblemID, &submit.Confirm)
        if err != nil {
            return nil, err
        }
        submits = append(submits, submit)
    }
    return submits, nil
}
