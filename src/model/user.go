/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package model

import (
    "OIUP-Backend/config"
    _ "github.com/mattn/go-sqlite3"
)

type LanguageType int32

const (
    LanguageCPlusPlus LanguageType = 1
    LanguageC         LanguageType = 2
    LanguagePascal    LanguageType = 3
)

type UserInfo struct {
    Name      string
    School    string
    ContestID string
    PersonID  string
    Language  LanguageType
}

var addUserQuery, _ = db.Prepare("INSERT INTO " + config.Config.DB.TableUser + " VALUES (?, ?, ?, ?, ?)")

func AddUser(user UserInfo) error {
    _, err := addUserQuery.Exec(user.Name, user.School, user.ContestID, user.PersonID, user.Language)
    return err
}

var deleteUserQuery, _ = db.Prepare("DELETE FROM " + config.Config.DB.TableUser + " WHERE contest_id = ?")

func DeleteUser(contestID string) error {
    _, err := deleteUserQuery.Exec(contestID)
    return err
}

var getUserQuery, _ = db.Prepare("SELECT * FROM " + config.Config.DB.TableUser + " WHERE contest_id = ?")

func GetUser(contestID string) (UserInfo, error) {
    var user UserInfo
    rows, err := getUserQuery.Query(contestID)
    if err != nil {
        return user, err
    }
    err = rows.Scan(&user.Name, &user.School, &user.ContestID, &user.PersonID, &user.Language)
    return user, err
}

var uploadUserQuery, _ = db.Prepare("UPDATE "  + config.Config.DB.TableUser + " SET name = ?, school = ?, person_id = ?, language = ? WHERE contest_id = ?")

func UpdateUser(user UserInfo) error {
    _, err := uploadUserQuery.Exec(user.Name, user.School, user.PersonID, user.Language, user.ContestID)
    return err
}
