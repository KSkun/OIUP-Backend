/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package model

import (
    "OIUP-Backend/config"
    "errors"
)

const (
    LanguageCPlusPlus = 1
    LanguageC         = 2
    LanguagePascal    = 3
)

type UserInfo struct {
    Name      string       `json:"name"`
    School    string       `json:"school"`
    ContestID string       `json:"contest_id"`
    PersonID  string       `json:"person_id"`
    Language  int          `json:"language"`
}

func AddUser(user UserInfo) error {
    _, found, err := GetUser(user.ContestID)
    if err != nil {
        return err
    }
    if found {
        return errors.New("user with contest_id " + user.ContestID + " already exists")
    }

    addUserQuery, _ := db.Prepare("INSERT INTO " + config.Config.DB.TableUser +
        " VALUES (?, ?, ?, ?, ?)")
    queryCh <- DBWriteQuery{
        Stmt:       addUserQuery,
        Parameters: []interface{}{user.Name, user.School, user.ContestID, user.PersonID, user.Language},
    }
    err = <-errCh
    return err
}

func DeleteUser(contestID string) error {
    _, found, err := GetUser(contestID)
    if err != nil {
        return err
    }
    if !found {
        return errors.New("user with contest_id " + contestID + " not found")
    }

    deleteUserQuery, _ := db.Prepare("DELETE FROM " + config.Config.DB.TableUser +
        " WHERE contest_id = ?")
    queryCh <- DBWriteQuery{
        Stmt:       deleteUserQuery,
        Parameters: []interface{}{contestID},
    }
    err = <-errCh
    return err
}

func GetUser(contestID string) (UserInfo, bool, error) {
    getUserQuery, _ := db.Prepare("SELECT * FROM " + config.Config.DB.TableUser +
        " WHERE contest_id = ?")
    defer getUserQuery.Close()
    var user UserInfo
    rows, err := getUserQuery.Query(contestID)
    defer rows.Close()
    if err != nil {
        return user, false, err
    }

    if !rows.Next() {
        return user, false, nil
    }
    err = rows.Scan(&user.Name, &user.School, &user.ContestID, &user.PersonID, &user.Language)
    return user, true, nil
}

func UpdateUser(user UserInfo) error {
    _, found, err := GetUser(user.ContestID)
    if err != nil {
        return err
    }
    if found {
        return errors.New("user with contest_id " + user.ContestID + " already exists")
    }

    updateUserQuery, _ := db.Prepare("UPDATE "  + config.Config.DB.TableUser +
        " SET name = ?, school = ?, person_id = ?, language = ? WHERE contest_id = ?")
    queryCh <- DBWriteQuery{
        Stmt:       updateUserQuery,
        Parameters: []interface{}{user.Name, user.School, user.PersonID, user.Language, user.ContestID},
    }
    err = <-errCh
    return err
}
