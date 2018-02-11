package base

import (
	"errors"
	"io"
	"log"
	"os"
	"time"

	"gopkg.in/configo.v2"
)

var outputName = "system.log"
var isLog = true

func init() {
	LogInitialize()
}

func LogInitialize() {
	wd, e := os.Getwd()
	path := "." + configo.GetSystemSeparator()

	if e == nil {
		path = wd + configo.GetSystemSeparator()

	}
	err := os.Rename(path+outputName, path+time.Now().Format("200102150405_")+outputName)
	if err != nil {
		//do nothing
	}
	CreateLogFile(path + outputName)
	log.Println("log start")
}

func CreateLogFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_SYNC, os.ModePerm)
	if err != nil {
		errors.New("create log file error")
	}
	w := io.MultiWriter(os.Stdout, file)
	log.SetOutput(w)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	return nil
}

func Println(v ...interface{}) {
	if isLog {
		log.Println(v)
	}
}

func Print(v ...interface{}) {
	if isLog {
		log.Print(v)
	}
}
