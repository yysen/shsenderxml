package main

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

var vconfig config

type config struct {
	AreaCode  string
	URL       string
	FinishOut string
	UserName  string
	Password  string
	Filedir   string
	Startmonth  string
	Stopmonth   string
	Uptime    string
}

func readConfig(fileName string) error {
	bs, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(bs, &vconfig); err != nil {
		return err
	}
	return nil
}
