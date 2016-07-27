package services

import (
	"io/ioutil"
	"os"
	"encoding/json"
	"os/signal"
	"syscall"
	"splash/logger"
	"runtime"
	"time"
)

type Locator struct {
}

func NewLocator() *Locator {
	return &Locator{}
}

func (self *Locator) LoadConfig(configPath *string) (map[string]interface{}, error) {

	file, err := ioutil.ReadFile(*configPath)

	if nil != err {

		return nil, err
	}

	config := map[string]interface{}{}

	err = json.Unmarshal(file, &config)

	if (nil != err) {
		return nil, err
	}

	return config, nil

}

func (*Locator) Logger() *logger.Logger {
	return logger.NewLogger()
}

func (*Locator) BlockIndefinitely() {

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	println("Blocking indefinitely...")
	<-sigc
	println("Bye Bye!")
}

func (self *Locator) GetAsTimestamp(nanoseconds int64) time.Time {
	return time.Unix(0, nanoseconds)
}

func (self *Locator) Stats() {

	logger := self.Logger()

	for {

		logger.Info("Number of Go routines:", runtime.NumGoroutine())
		time.Sleep(3 * time.Second)
	}
}
