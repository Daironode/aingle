 package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func logPrint() {
	Debug("debug")
	Info("info")
	Warn("warn")
	Error("error")
	Fatal("fatal")
	Trace("trace")

	testValue := 1
	Debugf("debug %v", testValue)
	Infof("info %v", testValue)
	Warnf("warn %v", testValue)
	Errorf("error %v", testValue)
	Fatalf("fatal %v", testValue)
	Tracef("trace %v", testValue)
}

func TestLog(t *testing.T) {
	defer func() {
		os.RemoveAll("Log/")
	}()

	InitLog(InfoLog, PATH, Stdout)
	Log().SetDebugLevel(DebugLog)
	logPrint()

	Log().SetDebugLevel(WarnLog)

	logPrint()

	err := ClosePrintLog()
	assert.Nil(t, err)
}

func TestNewLogFile(t *testing.T) {
	defer func() {
		os.RemoveAll("Log/")
	}()
	InitLog(InfoLog, PATH, Stdout)
	logfileNum1, err1 := ioutil.ReadDir("Log/")
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	logPrint()
	ClosePrintLog()
	time.Sleep(time.Second * 2)
	InitLog(InfoLog, PATH, Stdout)
	logfileNum2, err2 := ioutil.ReadDir("Log/")
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	assert.Equal(t, len(logfileNum1), len(logfileNum2)-1)
}

func TestDataRace(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			InitLog(DebugLog)
		}
	}()

	for i := 0; i < 10; i++ {
		Debug("aaaa")
	}

	wg.Wait()
}
