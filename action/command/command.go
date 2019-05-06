package command

import (
"golang.org/x/crypto/ssh"
"azh/subnginx/config"
"log"
"azh/subnginx/model"
"fmt"
"errors"
)

const err_conn  = "build up connection error"



func SshTo(ip model.Ip, buffer string, retry int) model.Call {
	for i:=1;i<=retry;i++{
		call:=sshto(ip,buffer)
		if call.Err != nil {
			log.Println(fmt.Sprintf("重试第%d次",retry))
			log.Println(fmt.Sprintf("%v",call))
			if retry == i {
				return call
			}
			continue
		}
		return call
	}
	return model.Call{
		errors.New("is over retry"),
		"xx",
		"xx",
	}
}
func sshto(ip model.Ip,buffer string) (call model.Call) {
	conn, err := ssh.Dial("tcp", ip.Address(), &ssh.ClientConfig {
		User: ip.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(ip.Pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		call.Err=err
		call.ErrorCommand=err_conn
		return
	}
	defer conn.Close()
	//run before
	for _,value:=range config.GetConfig().BeforeCommand{
		out,err:=runCommand(value,conn)
		if err != nil {
			call.Err=err
			call.ErrorCommand=value
			return
		}
		call.OutList+=out
	}
	out,err:=runCommand(buffer,conn)
	if err != nil {
		call.Err=err
		call.ErrorCommand=buffer
		return
	}
	call.OutList+=out
	out,err=runCommand(config.GetConfig().NginxRestartCommand,conn)
	if err != nil {
		call.Err=err
		call.ErrorCommand=config.GetConfig().NginxRestartCommand
		return
	}
	call.OutList+=out
	//run after
	for _,value:=range config.GetConfig().AfterCommand{
		out,err:=runCommand(value,conn)
		if err != nil {
			call.Err=err
			call.ErrorCommand=value
			return
		}
		call.OutList+=out
	}
	return
}
//./nginx –s reload
func runCommand(cmd string, conn *ssh.Client) (string,error){
	sess, err := conn.NewSession()
	if err != nil {
		log.Println(err)
		return "",err
	}
	defer sess.Close()
	b,err:=sess.Output(cmd)
	if err != nil {
		log.Println(err)
		return "",err
	}
	return string(b),nil
}