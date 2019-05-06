package ipsmanager

import (
	"azh/subnginx/model"
	"azh/subnginx/db"
	"azh/subnginx/action/template"
	"log"
	"fmt"
	"time"
	"sync"
	"azh/subnginx/action/command"
	"azh/subnginx/config"
)

type IpCenter struct {
	Hosts model.Ips
	TemplateLine model.TemplateLine
	timegap time.Duration
	round time.Timer
	//也可以在一轮结束时才将新的同步进来作为下一个目标版本（这样好像是最好的），也就是要主动去更新TemplateLine ,在轮次间隙做这个事情好像还可以
	//ip是可以随便变动的，但是在当前轮次，包括回滚操作都是一个ip版本，下一次才会生效
	model.Context
	sync.Mutex
}
//ip add update
func (ic *IpCenter)PutsIp(ip model.Ip,index int)  error{
	return ic.Hosts.Puts(ip,index)
}

//ip delete
func (ic *IpCenter)DeletesIp(ip model.Ip){
	ic.Hosts.Delete(ip)
}
//template update
func (ic *IpCenter)UpdateTemplate(){
	appconf:=db.SelectAll()
	nginxConf:=template.Template(appconf)
	log.Println(fmt.Sprintf("update template %s",nginxConf))
	ic.TemplateLine.Update(model.NewTemplate(nginxConf))
}
func (ic *IpCenter)Run(){
	ic.Lock()
	defer ic.Unlock()
}
func (ic *IpCenter)Wait(){
	<-ic.Done
}
func (ic *IpCenter)rollback(hosts []model.Ip)  {
	recent:=ic.TemplateLine.GetRecent(1)
	for index,ip:=range hosts{
		callback:=command.SshTo(ip,recent.TempText,config.GetConfig().RetryTimesWhenFail)
		if callback.Err!=nil {
			log.Println(callback.Err)
			ip.SyncStatus="OutOfSync"
		}else{
			ip.TempHashVersion=recent.HashCode
			ip.SyncStatus="Synchronized"
		}
		ip.RecentCallReport=callback
		err:=ic.PutsIp(ip,index)
		if err!=nil{
			log.Println("what happand")
		}
	}
}
func (ic *IpCenter)upgrade(){
	righthost:=ic.Hosts.GetAvailablePort()
	recent:=ic.TemplateLine.GetRecent( 0)
	for index,ip:=range righthost{
		callback:=command.SshTo(ip,recent.TempText,config.GetConfig().RetryTimesWhenFail)
		if callback.Err!=nil {
			log.Println(callback.Err)
			ip.SyncStatus="OutOfSync"
		}else{
			ip.TempHashVersion=recent.HashCode
			ip.SyncStatus="Synchronized"
		}
		ip.RecentCallReport=callback
		err:=ic.PutsIp(ip,index)
		if err!=nil{
			log.Println("what happand")
		}

		//TODO webhook
		if callback.Err!=nil {
			log.Println(callback.Err)
			//调用web hook推送消息
			switch config.GetConfig().FailAction {//goon rollback stop
			case "goon":
				log.Println("go on")
				continue
			case "rollback":
				log.Println("rollback")
				ic.rollback(righthost[:index])
				return
			case "stop":
				log.Println("stop")
				break
			default:
				continue
			}
		}

	}
}
func (ic *IpCenter)StopIpCent() {
	close(ic.Stop)
}
func (ic *IpCenter)run()  {
	go func() {
		for {
			select {
			case <-ic.round.C:
				log.Println("batch start ")
				ic.UpdateTemplate()
				ic.upgrade()
			case <-ic.Stop:
				log.Println("batch stop ")
				break
			}
			time.Sleep(ic.timegap)
		}
		close(ic.Done)
	}()
}
//在达到最大重试次数后可选择继续向后运行，停止运行，和回滚上个版本
//对版本的管理 //只增 在修改（包括删除）时触发 在达到重试最大次数的时候回滚需要替换最新的version为上一个版本，将根据最新的version作为强制同步版本进行同步

