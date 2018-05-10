package util

import (
	"log"
	"log/syslog"
	"runtime"
	"fmt"
	"os"
	"time"
)

type Logger struct {
	*log.Logger
}

func DoLog(c string, cate string) {
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	c = "From:" + f.Name() + c

	raddr := GetConfig("service", "log")
	sysLog, err := syslog.Dial("tcp", raddr, syslog.LOG_INFO, cate)
	if err != nil {
		localLog(c, cate)
	} else {
		sysLog.Info(c)
	}
}

func localLog(c, cat string) {
	err := os.Mkdir("./logs", 0777)
	if err != nil && os.IsExist(err) == false {
		fmt.Print(err)
	}

	f, err := os.OpenFile("./logs/" + cat + "_" + time.Now().Format(`20060102`) + ".log", os.O_APPEND|os.O_WRONLY, 0)
	if os.IsNotExist(err) == true {
		f, err = os.Create( "./logs/" + cat + "_" + time.Now().Format(`20060102`) + ".log")
	}

	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()
	_, erre := f.Write([]byte(c + "\n"))
	if erre != nil {
		fmt.Println(erre)
	}
}

func FUNCTION() string{
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}