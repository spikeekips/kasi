package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/op/go-logging"
	"github.com/spikeekips/kasi"
	"github.com/spikeekips/kasi/conf"
	"github.com/spikeekips/kasi/util"

	"gopkg.in/alecthomas/kingpin.v1"
)

var log *logging.Logger
var (
	app      = kingpin.New("kasi", "API gateway.")
	debug    = app.Flag("debug", "Enable debug mode.").Bool()
	logLevel = app.Flag(
		"loglevel",
		fmt.Sprintf("log level, {%s}", strings.Join(kasi_util.LogLevelNames, ", ")),
	).String()

	start         = app.Command("start", "Start new server")
	startConf     = start.Flag("conf", "Configuration file").Default("kasi.yml").String()
	stop          = app.Command("stop", "Stop server")
	stopConf      = stop.Flag("conf", "Configuration file").Default("kasi.yml").String()
	stopSock      = stop.Flag("sock", "socket file").Default("kasi.sock").String()
	reload        = app.Command("reload", "Reload server")
	reloadNewConf = reload.Flag("new-conf", "New configuration file").String()
	reloadSock    = reload.Flag("sock", "socket file").Default("kasi.sock").String()
)

func checkConf(conf string) (*kasi_conf.CoreSetting, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, errors.New("failed to get current directory.")
	}
	confTemp := kasi_util.JoinPath(wd, conf)

	state, err := os.Stat(confTemp)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New(fmt.Sprintf("conf file, `%s` no such file.", confTemp))
		}
		return nil, err
	}

	if state.Mode().IsDir() {
		return nil, errors.New(fmt.Sprintf("conf file, `%s` must be file.", confTemp))
	}

	// parse
	yml, err := ioutil.ReadFile(conf)
	if err != nil {
		return nil, err
	}
	setting, err := kasi.ParseConfig(string(yml))
	if err != nil {
		return nil, err
	}

	return setting, nil
}

func main() {
	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	// set loglevel
	if *debug {
		*logLevel = "DEBUG"
	}
	logLevelFound, logLevelNotFound := kasi_util.GetLogLevel(*logLevel)
	if logLevelNotFound != nil {
		logLevelFound = kasi.DefaultLogLevel
	}
	log = kasi.SetLogging(logLevelFound)

	var setting *kasi_conf.CoreSetting
	var err error

	// load config
	switch command {
	case start.FullCommand():
		// check conf
		setting, err = checkConf(*startConf)
		if len(*logLevel) < 1 && setting.Env.LogLevel >= 0 {
			log = kasi.SetLogging(setting.Env.LogLevel)
		} else {
			setting.Env.LogLevel = logLevelFound
		}
		if err != nil {
			kingpin.UsageErrorf(err.Error())
		}

	case stop.FullCommand():
		log.Debug("stop")

		// send stop signal thru socket

	case reload.FullCommand():
		// socket must be given

		log.Debug("reload")

		// send reload signal thru socket
	}

	switch command {
	case start.FullCommand():
		log.Debug("hello! this is `kasi`.")
		log.Debug("start")
		log.Debug("setting: %v", kasi_util.ToJson(setting))

		kasi.Start(setting)
	}

	os.Exit(0)
}
