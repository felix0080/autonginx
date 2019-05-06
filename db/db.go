package db

import (
	"sync"
	"azh/subnginx/model"
	"github.com/kataras/iris/core/errors"
	"azh/subnginx/util"
	"log"
	"azh/subnginx/config"
)
var m *Cmap
func init() {
	var cmp Cmap
	cmp.m=make(map[interface{}]interface{},100)
	m=&cmp
}
type Cmap struct {
	m map[interface{}]interface{}
	lock sync.RWMutex
}

func (c *Cmap)Delete(key interface{})  {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.m,key)
	//持久化
	if err:=util.Persistence(config.GetConfig().DataSavePath,m);err != nil {
		log.Println(err)
	}
}
func (c *Cmap) Get(key interface{}) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.m[key]
}

func (c *Cmap) Store(key , val interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.m[key] = val
	//持久化
	if err:=util.Persistence(config.GetConfig().DataSavePath,m);err != nil {
		log.Println(err)
	}
}
func Update(conf model.AppConf)  {
	m.Store(conf.AppName,conf)
}
func Delete(conf model.AppConf){
	m.Delete(conf.AppName)
}

func SelectWithUser(user string)  []model.AppConf{
	var ret []model.AppConf
	m.Range(func(key, value interface{}) bool {
		c,ok:=value.(model.AppConf)
		if ok && c.User == user {
			ret=append(ret,c)
		}
		return true
	})
	return ret
}
func SelectWithAppName(appname string) (model.AppConf,error){
	appconf:=m.Get(appname)
	c,ok:=appconf.(model.AppConf)
	if ok {
		return c,nil
	}
	return c,errors.New("not fount")
}
//已经删除的应用应当删除此数据库中数据
func SelectWithProjectName(projectName string) []model.AppConf {
	var ret []model.AppConf
	m.Range(func(key, value interface{}) bool {
		c,ok:=value.(model.AppConf)
		if ok && (c.ProjectName==projectName){
			ret=append(ret,c)
		}
		return true
	})
	return ret
}
func SelectAll() []model.AppConf {
	var ret []model.AppConf
	m.Range(func(key, value interface{}) bool {
		c,ok:=value.(model.AppConf)
		if ok{
			ret=append(ret,c)
		}
		return true
	})
	return ret
}
func (c *Cmap)Range(funcs func(key, value interface{})bool){
	c.lock.RLock()
	defer c.lock.RUnlock()
	for key,value:=range c.m {
		if !funcs(key, value) {
			break
		}
	}
}