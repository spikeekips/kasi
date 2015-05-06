package kasi

import (
	"net"
	"os"

	"github.com/spikeekips/kasi/conf"
)

var stopChannel chan int
var setting *conf.CoreSetting

func Start(settingInput *conf.CoreSetting) {
	log = SetLogging(settingInput.Env.LogLevel)

	setting = settingInput

	HTTPServe(setting)

	go func() {
		log.Debug("start management unixsocket.")
		os.Remove("./kasi.sock")
		l, err := net.Listen("unix", "./kasi.sock")
		if err != nil {
			log.Fatal("listen error:", err)
		}
		for {
			_, err := l.Accept()
			if err != nil {
				log.Fatal("accept error:", err)
			}
		}
	}()

	stopChannel = make(chan int)
	switch <-stopChannel {
	case 1:
		println("stop")
	}
}

func Stop(sock string) {
	stopChannel <- 1
}

// reload configuration
func Reload(conf string, sock string) {
	settingInput, err := ParseConfig(conf)
	if err != nil {
		return
	}

	log = SetLogging(settingInput.Env.LogLevel)
	setting = settingInput
}
