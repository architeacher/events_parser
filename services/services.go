package services

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"splash/logger"
	"syscall"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
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

	if nil != err {
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

var src = rand.NewSource(time.Now().UnixNano())

func (self *Locator) RandString(prefix string, n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return prefix + string(b)
}

func (self *Locator) Stats() {

	logger := self.Logger()

	for {

		logger.Info("Number of Go routines:", runtime.NumGoroutine())
		time.Sleep(3 * time.Second)
	}
}
