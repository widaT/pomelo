package pomelo

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Elogger interface {
	Error(msg string, err ...interface{})
}

type ErrLog struct {
	f *os.File
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func NewErrLogger(filePath string) Elogger {
	if filePath == "" {
		return &ErrLog{}
	}
	path, _ := filepath.Abs(filePath)
	dir := filepath.Dir(path)
	if exists, _ := PathExists(dir); !exists {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			log.Fatal("create err log path err", err)
		}
	}
	fd, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("create err log file err", err)
	}
	os.Chmod(filePath, 0666)
	return &ErrLog{f: fd}
}

func (e *ErrLog) Error(msg string, err ...interface{}) {
	if e.f == nil {
		e.f = os.Stderr
	}
	msg = fmt.Sprintf(msg, err...)
	t := time.Now().Format("2006-01-02 15:04:05")
	e.f.WriteString(t + " [Err] " + msg + "\n")
}
