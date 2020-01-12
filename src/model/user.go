/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package model

import (
    "OIUP-Backend/config"
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
    addUserQuery, _ := db.Prepare("INSERT INTO " + config.Config.DB.TableUser +
        " VALUES (?, ?, ?, ?, ?)")
    _, err := addUserQuery.Exec(user.Name, user.School, user.ContestID, user.PersonID, user.Language)
    return err
}

func DeleteUser(contestID string) error {
    deleteUserQuery, _ := db.Prepare("DELETE FROM " + config.Config.DB.TableUser +
        " WHERE contest_id = ?")
    _, err := deleteUserQuery.Exec(contestID)
    return err
}

func GetUser(contestID string) (UserInfo, bool, error) {
    getUserQuery, _ := db.Prepare("SELECT * FROM " + config.Config.DB.TableUser +
        " WHERE contest_id = ?")
    var user UserInfo
    rows, err := getUserQuery.Query(contestID)
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
    updateUserQuery, _ := db.Prepare("UPDATE "  + config.Config.DB.TableUser +
        " SET name = ?, school = ?, person_id = ?, language = ? WHERE contest_id = ?")
    _, err := updateUserQuery.Exec(user.Name, user.School, user.PersonID, user.Language, user.ContestID)
    return err
}
