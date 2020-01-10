/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package config

import (
	"encoding/json"
	"io/ioutil"
)

type HTTPConfig struct {
	Port      		   int  	  `json:"port"`
}

type JWTConfig struct {
	JWTSigningMethod   string	  `json:"signing_method"`
	JWTUserSecret      string 	  `json:"secret_user"`
	JWTBackstageSecret string 	  `json:"secret_backstage"`
	JWTTokenLife       int        `json:"token_life"`       // Unit: minute
}

type DBConfig struct {
	DBFile			   string     `json:"db_file"`
	TableUser		   string	  `json:"table_user"`
	TableProblem	   string	  `json:"table_problem"`
	TableSubmit 	   string	  `json:"table_submit"`
}

type FileConfig struct {
	DirectoryUpload	   string     `json:"directory_upload"`
	DirectorySource	   string     `json:"directory_source"`
}

type ConfigObject struct {
	HTTP               HTTPConfig `json:"http"`
	JWT 			   JWTConfig  `json:"jwt"`
	DB				   DBConfig   `json:"db"`
	File			   FileConfig `json:"file"`
}

var Config ConfigObject

func init() {
	configFile, _ := ioutil.ReadFile("config.json")
	_ = json.Unmarshal(configFile, &Config)
}
