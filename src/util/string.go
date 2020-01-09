/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package util

import (
    "regexp"
)

var contestIDRegex = regexp.MustCompile("HB-\\d+")

func CheckContestID(contestID string) bool {
    if len(contestID) != 8 {
        return false // check length
    }
    return contestIDRegex.MatchString(contestID)
}

var stringRegex = regexp.MustCompile("[0-9A-Za-z\\-]+")

func CheckString(str string) bool {
    return stringRegex.MatchString(str)
}
