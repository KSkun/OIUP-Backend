package model

import (
    "OIUP-Backend/config"
    "encoding/json"
    "errors"
    "time"
)

const (
    SubmitUnconfirmed = 0
    SubmitConfirmed   = 1
)

type SubmitInfo struct {
    ID        string
    User      string
    MD5       string
    MD5Set    []MD5Info
    Time      int
    ProblemID int
    Confirm   int
}

type MD5Info struct {
    MD5    string `json:"md5"`
    TestID int    `json:"test_id"`
}

func AddCodeSubmit(submitID string, user string, md5 string, time time.Time, problemID int) error {
    addCodeSubmitQuery, _ := db.Prepare("INSERT INTO " + config.Config.DB.TableSubmit +
        " VALUES (?, ?, ?, ?, ?, ?)")
    defer addCodeSubmitQuery.Close()
    _, err := addCodeSubmitQuery.Exec(submitID, user, md5, time.Unix(), problemID, SubmitUnconfirmed)
    return err
}

func AddOutputSubmit(submitID string, user string, md5Set []MD5Info, time time.Time, problemID int) error {
    md5JSON, err := json.Marshal(md5Set)
    if err != nil {
        return err
    }

    addOutputSubmitQuery, _ := db.Prepare("INSERT INTO " + config.Config.DB.TableSubmit +
        " VALUES (?, ?, ?, ?, ?, ?)")
    defer addOutputSubmitQuery.Close()
    _, err = addOutputSubmitQuery.Exec(submitID, user, md5JSON, time.Unix(), problemID, SubmitUnconfirmed)
    return err
}

func GetSubmit(submitID string) (SubmitInfo, bool, error) {
    getSubmitQuery, _ := db.Prepare("SELECT * FROM " + config.Config.DB.TableSubmit +
        " WHERE id = ?")
    defer getSubmitQuery.Close()
    rows, err := getSubmitQuery.Query(submitID)
    defer rows.Close()
    if err != nil {
        return SubmitInfo{}, false, err
    }

    var submit SubmitInfo
    if !rows.Next() {
        return submit, false, nil
    }
    err = rows.Scan(&submit.ID, &submit.User, &submit.MD5, &submit.Time, &submit.ProblemID, &submit.Confirm)

    // MD5Info array unmarshal
    problemConfig, _ := config.GetProblemConfig(submit.ProblemID)
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
    getLatestSubmitQuery, _ := db.Prepare("SELECT * FROM " + config.Config.DB.TableLatestSubmit +
        " WHERE user = ? AND problem_id = ?")
    defer getLatestSubmitQuery.Close()
    rows, err := getLatestSubmitQuery.Query(user, problemID)
    defer rows.Close()
    if err != nil {
        return SubmitInfo{}, false, err
    }

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
    confirmSubmitQuery, _ := db.Prepare("UPDATE " + config.Config.DB.TableSubmit +
        " SET confirm = ? WHERE id = ?")
    defer confirmSubmitQuery.Close()
    _, err = confirmSubmitQuery.Exec(SubmitConfirmed, submitID)
    if err != nil {
        return err
    }

    deleteLatestSubmitQuery, _ := db.Prepare("DELETE FROM " + config.Config.DB.TableLatestSubmit +
        " WHERE user = ? AND problem_id = ?")
    defer deleteLatestSubmitQuery.Close()
    _, err = deleteLatestSubmitQuery.Exec(submit.User, submit.ProblemID)
    if err != nil {
        return err
    }

    addLatestSubmitQuery, _ := db.Prepare("INSERT INTO " + config.Config.DB.TableLatestSubmit +
        " VALUES (?, ?, ?)")
    defer addLatestSubmitQuery.Close()
    _, err = addLatestSubmitQuery.Exec(submit.User, submitID, submit.ProblemID)
    return err
}
