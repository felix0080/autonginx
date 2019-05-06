package model

import (
	"fmt"
	"errors"
	"azh/subnginx/util"
	"time"
	"sync"
)

type Ip struct {
	Addr string 	`192.168.0.1`
	Port int
	User string 	`root`
	Pass string		`Cmb@2`
	TempHashVersion string  // asdasojcxich123jsajd
	SyncStatus string //	Synchronized/OutOfSync
	Alive bool 	//the hosts is alive?
	ScanPort int	`80`
	RecentCallReport Call
}


func (ip Ip)Address() string {
	return fmt.Sprintf("%s:%d",ip.Addr,ip.Port)
}
type Ips struct {
	sync.Mutex
	ips []Ip
	len int
	index int
}
func (ip Ip) Scan () bool {
	return util.IsOpen(ip.Address())
}
// replace if the ip addr repeat
//	update
// just for speed ， suppose the index is not change
func (p *Ips) Puts(ip Ip,index int)(error){
	p.Lock()
	defer p.Unlock()
	isExist:=-1
	if index != -1 {
		if ip.Addr == p.ips[index].Addr {
			isExist = index
		}
	}
	if isExist == -1{
		for i:=0;i<=p.index;i++{
			if  p.ips[i].Addr == ip.Addr {
				isExist=i
				break
			}
		}
	}
	if p.index == p.len-1 && isExist != -1 {
		return errors.New("index out of range ")
	}
	if isExist!=-1 {
		p.ips[isExist]=ip
	}else{
		p.index++
		p.ips[p.index] = ip
	}
	return nil
}
func (p *Ips) Delete(ip Ip){
	p.Lock()
	defer p.Unlock()
	isExist:=-1
	for i:=0;i<=p.index;i++{
		if  p.ips[i].Addr == ip.Addr {
			isExist=i
			break
		}
	}
	if isExist!=-1 {
		p.ips[isExist]=ip
		p.ips = append(p.ips[:isExist],p.ips[isExist+1:]...)
		p.index--
	}
}
// 返回 ok,err
func (p *Ips)ScanDog(stop chan struct{}){
	timer:=time.NewTimer(time.Minute)
	for {
		select {
		case <-stop:
			return
		case <-timer.C:
			p.Lock()
			for index,value:=range p.ips{
				hostAlived:=value.Scan()
				value.Alive=hostAlived
				p.ips[index]=value
			}
			p.Unlock()
		}
	}
}
// 返回 ok,err
func (p Ips) GetAvailablePort() []Ip {
	var ret1 []Ip
	for _,value:=range p.ips{
		if value.Alive {
			ret1=append(ret1,value)
		}
	}
	return ret1
}
type Call struct {
	Err error
	ErrorCommand string
	OutList string
}
