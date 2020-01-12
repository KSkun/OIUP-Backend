package model

import (
    "OIUP-Backend/config"
    "encoding/json"
    "errors"
    "github.com/satori/go.uuid"
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

func AddCodeSubmit(user string, md5 string, time time.Time, problemID int) (string, error) {
    addCodeSubmitQuery, _ := db.Prepare("INSERT INTO " + config.Config.DB.TableSubmit +
        " VALUES (?, ?, ?, ?, ?, ?)")
    submitID := uuid.NewV4()
    _, err := addCodeSubmitQuery.Exec(submitID.String(), user, md5, time.Unix(), problemID, 0)
    return submitID.String(), err
}

func AddOutputSubmit(user string, md5Set []MD5Info, time time.Time, problemID int) (string, error) {
    md5JSON, err := json.Marshal(md5Set)
    if err != nil {
        return "", err
    }

    addOutputSubmitQuery, _ := db.Prepare("INSERT INTO " + config.Config.DB.TableSubmit +
        " VALUES (?, ?, ?, ?, ?, ?)")
    submitID := uuid.NewV4()
    _, err = addOutputSubmitQuery.Exec(submitID.String(), user, md5JSON, time.Unix(), problemID, 0)
    return submitID.String(), err
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
    problemConfig, err := config.GetProblemConfig(submit.ProblemID)
    if err != nil {
        return submit, false, err
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
    getLatestSubmitQuery, _ := db.Prepare("SELECT * FROM " + config.Config.DB.TableLatestSubmit +
        " WHERE user = ?, problem_id = ?")
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
    confirmSubmitQuery, _ := db.Prepare("UPDATE " + config.Config.DB.TableSubmit +
        " SET confirm = ? WHERE id = ?")
    _, err := confirmSubmitQuery.Exec(1, submitID)
    return err
}
