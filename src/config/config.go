/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
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

type ProblemType int

const (
	ProblemClassic  ProblemType = 1
	ProblemAnswer   ProblemType = 2
	ProblemInteract ProblemType = 3
)

type ProblemInfo struct {
	ID 		           int           `json:"id"`
	Name 	           string        `json:"name"`
	Filename           string        `json:"filename"`
	TimeLimit          string        `json:"time_limit"`
	SpaceLimit         string        `json:"space_limit"`
	Type               ProblemType   `json:"type"`
}

type ContestStatus int

const (
	ContestStatusDefault ContestStatus = 1
	ContestStatusError   ContestStatus = -1
)

type ProblemMap map[int]ProblemInfo

type ContestConfig struct {
	Name			   string 	     `json:"name"`
	Status             ContestStatus `json:"status"`
	Message            string        `json:"message"`
	StartTime		   string        `json:"start_time"`
	Duration           float32       `json:"duration"`      // Unit: Hour
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
	_ = json.Unmarshal(configFile, &Config)

	Config.Contest.Problems = make(ProblemMap, 0)
	for _, problem := range Config.Contest.Problems {
		Config.Contest.Problems[problem.ID] = problem
	}
}

func SaveConfig() {
	configJSON, _ := json.Marshal(Config)
	_ = os.Rename("config.json", "config-" + time.Now().String() + ".json.bak")
	_ = ioutil.WriteFile("config.json", configJSON, os.ModePerm)
}

func GetProblemConfig(problemID int) (ProblemInfo, error) {
	problem, found := Config.Contest.Problems[problemID]
	if !found {
		return problem, errors.New("problem with id " + strconv.Itoa(problemID) + " not found")
	}
	return problem, nil
}
