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
	"time"
)

type HTTPConfig struct {
	Port      		   int  	     `json:"port"`
	BackstageKey       string        `json:"backstage_key"`
}

type JWTConfig struct {
	JWTSigningMethod   string	     `json:"signing_method"`
	JWTSecret      string 	         `json:"secret"`
	JWTTokenLife       int           `json:"token_life"`       // Unit: Minute
}

type DBConfig struct {
	DBFile			   string        `json:"db_file"`
	TableUser		   string
	TableSubmit 	   string
	TableLatestSubmit  string
	ChannelBuffer      int           `json:"channel_buffer"`
	RecordsPerPage     int           `json:"records_per_page"`
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

func ApplyConfig(configObj ConfigObject) error {
	configObj.DB.TableUser = "user"
	configObj.DB.TableSubmit = "submit"
	configObj.DB.TableLatestSubmit = "latest_submit"

	configObj.Contest.Problems = make(ProblemMap, 0)
	for _, problem := range configObj.Contest.ProblemSet {
		configObj.Contest.Problems[problem.ID] = problem
	}

	startTime, err := time.ParseInLocation("2006-01-02 15:04", configObj.Contest.StartTimeStr, time.Local)
	if err != nil {
		return errors.New("config error: invalid start_time, " + err.Error())
	}
	configObj.Contest.StartTime = startTime

	Config = configObj
	return nil
}

func LoadConfig() error {
	var configObj ConfigObject

	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(configFile, &configObj)
	if err != nil {
		return err
	}

	return ApplyConfig(configObj)
}

func SaveConfig() error {
	configJSON, err := json.MarshalIndent(Config, "", "  ")
	if err != nil {
		return err
	}
	err = os.Rename("config.json", "config-" + time.Now().Format("200601021504") + ".json.bak")
	if err != nil {
		return err
	}
	return ioutil.WriteFile("config.json", configJSON, os.ModePerm)
}

func GetProblemConfig(problemID int) (ProblemInfo, bool) {
	problem, found := Config.Contest.Problems[problemID]
	if !found {
		return ProblemInfo{}, false
	}
	return problem, true
}

func init() {
	err := LoadConfig()
	if err != nil {
		panic(err)
	}
}
