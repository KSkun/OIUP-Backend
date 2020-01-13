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
    _, err = addOutputSubmitQuery.Exec(submitID, user, md5JSON, time.Unix(), problemID, SubmitUnconfirmed)
    return err
}

func GetSubmit(submitID string) (SubmitInfo, bool, error) {
    getSubmitQuery, _ := db.Prepare("SELECT * FROM " + config.Config.DB.TableSubmit +
        " WHERE id = ?")
    rows, err := getSubmitQuery.Query(submitID)
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
    rows, err := getLatestSubmitQuery.Query(user, problemID)
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

    confirmSubmitQuery, _ := db.Prepare("UPDATE " + config.Config.DB.TableSubmit +
        " SET confirm = ? WHERE id = ?")
    _, err = confirmSubmitQuery.Exec(SubmitConfirmed, submitID)
    if err != nil {
        return err
    }

    deleteLatestSubmitQuery, _ := db.Prepare("DELETE FROM " + config.Config.DB.TableLatestSubmit +
        " WHERE user = ? AND problem_id = ?")
    _, err = deleteLatestSubmitQuery.Exec(submit.User, submit.ProblemID)
    if err != nil {
        return err
    }

    addLatestSubmitQuery, _ := db.Prepare("INSERT INTO " + config.Config.DB.TableLatestSubmit +
        " VALUES (?, ?, ?)")
    _, err = addLatestSubmitQuery.Exec(submit.User, submitID, submit.ProblemID)
    return err
}
