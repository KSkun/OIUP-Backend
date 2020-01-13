/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type HTTPConfig struct {
	Port      		   int  	     `json:"port"`
}

type JWTConfig struct {
	JWTSigningMethod   string	     `json:"signing_method"`
	JWTUserSecret      string 	     `json:"secret_user"`
	JWTBackstageSecret string 	     `json:"secret_backstage"`
	JWTTokenLife       int           `json:"token_life"`       // Unit: Minute
}

type DBConfig struct {
	DBFile			   string        `json:"db_file"`
	TableUser		   string	     `json:"table_user"`
	TableSubmit 	   string	     `json:"table_submit"`
	TableLatestSubmit  string        `json:"table_latest_submit"`
}

type FileConfig struct {
	DirectoryUpload	   string        `json:"directory_upload"`
	DirectorySource	   string        `json:"directory_source"`
}

const (
	ProblemClassic  = 1
	ProblemAnswer   = 2
	ProblemInteract = 3
)

type ProblemInfo struct {
	ID 		           int           `json:"id"`
	Name 	           string        `json:"name"`
	Filename           string        `json:"filename"`
	TimeLimit          string        `json:"time_limit"`
	SpaceLimit         string        `json:"space_limit"`
	Type               int           `json:"type"`
}

const (
	ContestStatusDefault int = 1
	ContestStatusError   int = -1
)

type ProblemMap map[int]ProblemInfo

type ContestConfig struct {
	Name			   string 	     `json:"name"`
	Status             int           `json:"status"`
	Message            string        `json:"message"`
	StartTimeStr	   string        `json:"start_time"`
	StartTime          time.Time
	Duration           float32       `json:"duration"`      // Unit: Hour
	Download           string        `json:"download"`
	UnzipToken         string        `json:"unzip_token"`
	UnzipShift         int           `json:"unzip_shift"`   // Unit: Minute
	ProblemSet		   []ProblemInfo `json:"problems"`
	Problems		   ProblemMap
}

type ConfigObject struct {
	HTTP               HTTPConfig    `json:"http"`
	JWT 			   JWTConfig     `json:"jwt"`
	DB				   DBConfig      `json:"db"`
	File			   FileConfig    `json:"file"`
	Contest            ContestConfig `json:"contest"`
}

var Config ConfigObject

func init() {
	configFile, _ := ioutil.ReadFile("config.json")
	err := json.Unmarshal(configFile, &Config)
	if err != nil {
		panic(err)
	}

	Config.Contest.Problems = make(ProblemMap, 0)
	for _, problem := range Config.Contest.ProblemSet {
		Config.Contest.Problems[problem.ID] = problem
	}

	startTime, err := time.ParseInLocation("2006-01-02 15:04", Config.Contest.StartTimeStr, time.Local)
	if err != nil {
		panic("config error: invalid start_time, " + err.Error())
	}
	Config.Contest.StartTime = startTime
}

func SaveConfig() {
	configJSON, _ := json.Marshal(Config)
	_ = os.Rename("config.json", "config-" + time.Now().String() + ".json.bak")
	_ = ioutil.WriteFile("config.json", configJSON, os.ModePerm)
}

func GetProblemConfig(problemID int) (ProblemInfo, bool) {
	problem, found := Config.Contest.Problems[problemID]
	if !found {
		return ProblemInfo{}, false
	}
	return problem, true
}
