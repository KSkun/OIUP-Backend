/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package util

import (
    "OIUP-Backend/config"
    "OIUP-Backend/model"
    "errors"
)

func GetUploadPath(submitID string) string {
    return config.Config.File.DirectoryUpload + "/" + submitID + "/"
}

func GetSourcePath(contestID string, problemFilename string) string {
    return config.Config.File.DirectorySource + "/" + contestID + "/" + problemFilename + "/"
}

func GetTempPath(filename string) string {
    return config.Config.File.DirectoryTemp + "/" + filename
}

func GetCodeSuffix(language int) (string, error) {
    switch language {
    case model.LanguageCPlusPlus:
        return ".cpp", nil
    case model.LanguageC:
        return ".c", nil
    case model.LanguagePascal:
        return ".pas", nil
    }
    return "", errors.New("不支持该语言类型！")
}
