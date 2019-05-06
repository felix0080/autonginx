package util

import (
	"log"
	"os"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"net"
	"time"
)
//read list
func Getlist(path string)([]byte,bool,error) {
	_,err := os.Stat(path)
	if err == nil {
		bs,err:=ioutil.ReadFile(path)
		return bs,true,err
	}
	if os.IsNotExist(err) {
		return []byte{},false,err
	}
	return []byte{},true,err
}
//i 传递指针
func Persistence(path string ,i interface{})error{
	b,err:=json.Marshal(i)
	if err!=nil {
		return err
	}
	err=savelist(path,b)
	if err !=nil {
		return err
	}
	return nil
}
//save in list
func savelist(path string,b []byte)error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Println(err)
		return err
	}
	defer f.Close()
	_,err=f.Write(b)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func IsOpen(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr,2*time.Second)
	if err != nil {
		fmt.Printf("Fail to connect, %s\n", err)
		return false
	}
	conn.Close()
	return true
}
