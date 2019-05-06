package config

import (
	"azh/subnginx/util"
	"encoding/json"
	"azh/subnginx/model"
)

var config model.Config
//get config from file
func init() {
	bs,ise,err:=util.Getlist("/etc/subnginx/config")
	if err != nil || !ise{
		panic(err)
	}
	err=json.Unmarshal(bs,&config)
	if err != nil {
		panic(err)
	}
}
func GetConfig() model.Config {
	return config
}