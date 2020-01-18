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

func getSQLConditionsStr(filters map[string]interface{}) string {
    str := ""
    for key := range filters {
        if key == "language" {
            str += key + " = ? AND "
            continue
        }
        str += key + " LIKE ? AND "
    }
    if len(filters) != 0 {
        str = str[0:len(str) - 5]
    }
    return str
}

func SearchUser(filters map[string]interface{}, page int) ([]UserInfo, error) {
    conditionsStr := getSQLConditionsStr(filters)

    searchUserSubQueryStr := "SELECT contest_id FROM " + config.Config.DB.TableUser
    if len(conditionsStr) > 0 {
        searchUserSubQueryStr += " WHERE " + conditionsStr
    }
    searchUserSubQueryStr += " ORDER BY contest_id LIMIT " +
        strconv.Itoa(config.Config.DB.RecordsPerPage * (page - 1))

    searchUserQueryStr := "SELECT * FROM " + config.Config.DB.TableUser +
        " WHERE contest_id NOT IN (" + searchUserSubQueryStr + ")"
    if len(conditionsStr) > 0 {
        searchUserQueryStr += " AND " + conditionsStr
    }
    searchUserQueryStr += " ORDER BY contest_id LIMIT " +
        strconv.Itoa(config.Config.DB.RecordsPerPage * page)

    searchUserQuery, err := db.Prepare(searchUserQueryStr)
    if err != nil {
        return nil, err
    }
    defer searchUserQuery.Close()
    values := make([]interface{}, 0)
    for key, value := range filters {
        if key == "language" {
            values = append(values, value)
            continue
        }
        values = append(values, value.(string) + "%")
    }
    values = append(values, values...)
    rows, err := searchUserQuery.Query(values...)
    defer rows.Close()
    if err != nil {
        return nil, err
    }

    users := make([]UserInfo, 0)
    for rows.Next() {
        var user UserInfo
        err = rows.Scan(&user.Name, &user.School, &user.ContestID, &user.PersonID, &user.Language)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
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
