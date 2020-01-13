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

func GetCodeSuffix(language int) (string, error) {
    switch language {
    case model.LanguageCPlusPlus:
        return ".cpp", nil
    case model.LanguageC:
        return ".c", nil
    case model.LanguagePascal:
        return ".pas", nil
    }
    return "", errors.New("invalid language type")
}
