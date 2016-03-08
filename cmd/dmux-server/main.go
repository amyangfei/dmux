package main

import (
	"flag"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/koding/multiconfig"
	"log"
	"github.com/amyangfei/dmux/registry/controllers"
	"github.com/amyangfei/dmux/store"
	"github.com/amyangfei/dmux/utils"
	"net/http"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
)

type (
	ServerConfig struct {
		C      string
		Base   BaseConfig
		Plugin PluginConfig
	}

	BaseConfig struct {
		LogFile      string `default:"./dmux_server.log"`
		RegistryHost string `default:"0.0.0.0"`
		RegistryPort int    `default:"6000"`
		StoreDB      string `required:"true"`
		DBPath       string `default:"."`
	}

	PluginConfig struct {
		Plugins []string
	}
)

var Config *ServerConfig

func initConfig(configFile string) error {
	m := multiconfig.NewWithPath(configFile)
	Config = &ServerConfig{}
	m.MustLoad(Config)

	if !utils.In(Config.Base.StoreDB, []string{"goleveldb", "memory"}) {
		return fmt.Errorf("db model [%s] is not support\n", Config.Base.StoreDB)
	}

	return nil
}

func initStorage() error {
	var storage store.Storage
	var err error
	if Config.Base.StoreDB == "goleveldb" {
		dbpath := path.Clean(path.Join(Config.Base.DBPath, "dmux-server.db"))
		log.Printf("dbpath: %s\n", dbpath)
		storage, err = store.NewLevelStore(dbpath)
	} else {
		storage, err = store.NewMemStore()
	}
	if err != nil {
		storage.Close()
		return fmt.Errorf("store init error: %s", err)
	}
	controllers.InitEntryManager(storage)
	return nil
}

func registryService() {
	errorHandlers := map[string]func(w http.ResponseWriter, r *http.Request){
		"400": controllers.Error400,
		"404": controllers.Error404,
		"500": controllers.Error500,
	}
	beego.Router("/entry/list", &controllers.RegistryController{}, "get:ListEntries")
	beego.Router("/entry/new", &controllers.RegistryController{}, "post:NewEntry")
	beego.Router("/entry/query", &controllers.RegistryController{}, "get:GetEntry")
	beego.Router("/entry/delete", &controllers.RegistryController{}, "post:DelEntry")
	for code, handler := range errorHandlers {
		beego.ErrorHandler(code, handler)
	}
	beego.Run()
}

func singalHandle(s os.Signal) {
	switch s {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
		log.Printf("recevie signal %v, exit.", s)
		os.Exit(0)
	case syscall.SIGHUP:
		// TODO: reload
		log.Printf("server reload...")
	}
}

func main() {
	var configFile string
	var printVersion bool

	flag.BoolVar(&printVersion, "version", false, "print version")
	flag.StringVar(&configFile, "c", "config.toml", "path to config file")
	flag.Parse()

	if printVersion {
		utils.PrintVersion()
		os.Exit(0)
	}

	if err := initConfig(configFile); err != nil {
		panic(err)
	}

	if err := initStorage(); err != nil {
		panic(err)
	}

	go registryService()

	stop := make(chan os.Signal)
	signal.Notify(
		stop, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	var wg sync.WaitGroup

	select {
	case s := <-stop:
		singalHandle(s)
	}
	wg.Wait()
}
