package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	ws "shortsrv/wshort"

	"github.com/gorilla/mux"
	yaml "gopkg.in/yaml.v3"
)

var (
	URLInBase64Param = "urlinbase64"
	ShortIDParam     = "surlid"
)

type WsConfig struct {
	Addr          string   `yaml:"Address"`
	EtcdEndpoints []string `yaml:"EtcdEndpoints"`
	EtcdPrefix    string   `yaml:"EtcdPrefix"`
}

func main() {
	var conf WsConfig
	confBytes, err := os.ReadFile("/etc/wshort.conf")
	if err != nil {
		fmt.Printf("Error in reading config file: %s\n", err.Error())
		os.Exit(1)
	}
	err = yaml.Unmarshal(confBytes, &conf)
	if err != nil {
		fmt.Printf("Error in unmarshal config: %s\n", err.Error())
		os.Exit(1)
	}
	if len(conf.EtcdEndpoints) == 0 {
		fmt.Println("Error in config: EtcdEndpoints is required")
		os.Exit(1)
	}
	if conf.EtcdPrefix == "" {
		conf.EtcdPrefix = "wshort"
	}
	if conf.Addr == "" {
		fmt.Println("Error in config: Address is required")
		os.Exit(1)
	}
	if err := ws.InitStore(conf.EtcdEndpoints, conf.EtcdPrefix); err != nil {
		fmt.Printf("Error in initializing store: %s\n", err.Error())
		os.Exit(1)
	}
	defer ws.CloseStore()
	ws.Logger.Printf("starting server on %s with etcd endpoints %s and prefix %s", conf.Addr, strings.Join(conf.EtcdEndpoints, ","), conf.EtcdPrefix)

	router := mux.NewRouter()
	router.HandleFunc("/s", ShortUrlPutHandler).Methods(http.MethodPut)
	router.HandleFunc("/s/{"+ShortIDParam+"}", ShortUrlGetHandler).Methods(http.MethodGet)
	router.HandleFunc("/debug", DebugHandler).Methods(http.MethodGet)

	s := &http.Server{
		Addr:    conf.Addr,
		Handler: router,
	}
	err = s.ListenAndServe()
	if err != nil {
		fmt.Printf("server started error: %s", err.Error())
		os.Exit(1)
	}
}
