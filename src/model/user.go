/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package model

import (
    "OIUP-Backend/config"
    "errors"
    "strconv"
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

    addUserQuery, err := db.Prepare("INSERT INTO " + config.Config.DB.TableUser +
        " VALUES (?, ?, ?, ?, ?)")
    if err != nil {
        return err
    }
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

    deleteUserQuery, err := db.Prepare("DELETE FROM " + config.Config.DB.TableUser +
        " WHERE contest_id = ?")
    if err != nil {
        return err
    }
    queryCh <- DBWriteQuery{
        Stmt:       deleteUserQuery,
        Parameters: []interface{}{contestID},
    }
    err = <-errCh
    return err
}

func GetUser(contestID string) (UserInfo, bool, error) {
    getUserQuery, err := db.Prepare("SELECT * FROM " + config.Config.DB.TableUser +
        " WHERE contest_id = ?")
    if err != nil {
        return UserInfo{}, false, err
    }
    defer getUserQuery.Close()
    var user UserInfo
    rows, err := getUserQuery.Query(contestID)
    if err != nil {
        return user, false, err
    }
    defer rows.Close()

    if !rows.Next() {
        return user, false, nil
    }
    err = rows.Scan(&user.Name, &user.School, &user.ContestID, &user.PersonID, &user.Language)
    return user, true, nil
}

func SearchUser(filters map[string]interface{}, page int) ([]UserInfo, int, error) {
    conditionsStr := getSQLConditionsStr(filters)
    values := make([]interface{}, 0)
    for key, value := range filters {
        if key == "language" {
            values = append(values, value)
            continue
        }
        values = append(values, value.(string) + "%")
    }

    var count int
    countUserQueryStr := "SELECT COUNT(*) FROM " + config.Config.DB.TableUser
    if len(conditionsStr) > 0 {
        countUserQueryStr += " WHERE " + conditionsStr
    }
    countUserQuery, err := db.Prepare(countUserQueryStr)
    if err != nil {
        return nil, 0, err
    }
    defer countUserQuery.Close()
    countRows, err := countUserQuery.Query(values...)
    if err != nil {
        return nil, 0, err
    }
    defer countRows.Close()
    countRows.Next()
    err = countRows.Scan(&count)
    if err != nil {
        return nil, 0, err
    }

    searchUserQueryStr := "SELECT * FROM " + config.Config.DB.TableUser
    if len(conditionsStr) > 0 {
        searchUserQueryStr += " WHERE " + conditionsStr
    }
    searchUserQueryStr += " ORDER BY contest_id LIMIT " +
        strconv.Itoa(config.Config.DB.RecordsPerPage * page) + " OFFSET " +
        strconv.Itoa(config.Config.DB.RecordsPerPage * (page - 1))
    searchUserQuery, err := db.Prepare(searchUserQueryStr)
    if err != nil {
        return nil, 0, err
    }
    defer searchUserQuery.Close()
    rows, err := searchUserQuery.Query(values...)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    users := make([]UserInfo, 0)
    for rows.Next() {
        var user UserInfo
        err = rows.Scan(&user.Name, &user.School, &user.ContestID, &user.PersonID, &user.Language)
        if err != nil {
            return nil, 0, err
        }
        users = append(users, user)
    }
    return users, count, nil
}

func UpdateUser(user UserInfo) error {
    _, found, err := GetUser(user.ContestID)
    if err != nil {
        return err
    }
    if !found {
        return errors.New("user with contest_id " + user.ContestID + " not found")
    }

    updateUserQuery, err := db.Prepare("UPDATE "  + config.Config.DB.TableUser +
        " SET name = ?, school = ?, person_id = ?, language = ? WHERE contest_id = ?")
    if err != nil {
        return err
    }
    queryCh <- DBWriteQuery{
        Stmt:       updateUserQuery,
        Parameters: []interface{}{user.Name, user.School, user.PersonID, user.Language, user.ContestID},
    }
    err = <-errCh
    return err
}
