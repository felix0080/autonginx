package server

import (
	"os"
	"os/signal"
	"net/http"
	"time"
	"log"
	"io/ioutil"
	"encoding/json"
	"azh/subnginx/model"
	"azh/subnginx/db"
)

func HttpServer() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	mux := http.NewServeMux()
	mux.HandleFunc("/apply",apply)
	server := &http.Server{
		Addr:         ":80",
		WriteTimeout: time.Second * 5,
		Handler:      mux,
	}
	go func() {
		<-quit
		if err := server.Close(); err != nil {
			log.Fatal("server has quit:", err)
		}
	}()
	log.Println("server started")
	//go action.DeelWith()
	err := server.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			log.Fatal("Server closed under request")
		} else {
			log.Fatal("Server closed unexpected", err)
		}
	}
	log.Fatal("Server exited")
}
func apply(w http.ResponseWriter, r *http.Request) {
	bs,err:=ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	var app model.AppConf
	err=json.Unmarshal(bs,&app)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	db.Update(app)
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}
func selectWithUser(w http.ResponseWriter, r *http.Request) {
	bs,err:=ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	var app model.AppConf
	err=json.Unmarshal(bs,&app)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	db.Update(app)
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}
func SelectWithAppName(w http.ResponseWriter, r *http.Request) {
	bs,err:=ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	var app model.AppConf
	err=json.Unmarshal(bs,&app)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	db.Update(app)
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}