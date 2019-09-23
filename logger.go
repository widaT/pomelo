package pomelo

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Logger interface {
	Log(msg string, err ...interface{})
}

type ErrLog struct {
	f Logger
}

func (e *ErrLog) Log(msg string, err ...interface{}) {
	msg = fmt.Sprintf(msg, err...)
	t := time.Now().Format("2006-01-02 15:04:05")
	msg = t + " [Err] " + msg + "\n"
	if e.f == nil {
		os.Stderr.WriteString(msg)
		return
	}
	e.f.Log(msg)
}

func NewErrLogger(config *Config) Logger {
	if config.ErrLog == "" {
		return &ErrLog{}
	}
	return &ErrLog{
		f: NewFileLog(config.ErrLog, config),
	}
}

type FileLog struct {
	f        *os.File
	maxSize  int64
	curSize  int64
	curNum   int
	FileName string
	Perm     os.FileMode
	MaxFiles int
	Lock     sync.Mutex
}

func NewFileLog(fileName string, config *Config) Logger {
	path, _ := filepath.Abs(fileName)
	dir := filepath.Dir(path)
	if exists, _ := PathExists(dir); !exists {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			log.Fatal("create err log path err", err)
		}
	}
	fd, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("create err log file err", err)
	}

	finfo, _ := fd.Stat()
	os.Chmod(path, 0666)
	f := &FileLog{
		f:        fd,
		curSize:  finfo.Size(),
		Perm:     0666,
		FileName: path,
		MaxFiles: config.LogMaxFiles,
		maxSize:  config.LogMaxSize,
	}
	return f
}

func (f *FileLog) doRotate() {
	f.curNum++
	if f.curNum > f.MaxFiles {
		f.curNum = 1
	}
	newFileName := fmt.Sprintf("%s.%d", f.FileName, f.curNum)
	f.f.Close()
	os.Rename(f.FileName, newFileName)
	fd, err := os.OpenFile(f.FileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("create err log file err", err)
	}
	os.Chmod(f.FileName, 0666)
	f.f = fd
}

func (f *FileLog) Log(msg string, err ...interface{}) {
	m := Str2byte(fmt.Sprintf(msg, err...))
	size := len(m)
	f.Lock.Lock()
	f.curSize += int64(size)
	if f.curSize > f.maxSize {
		f.doRotate()
		f.curSize = int64(size)
	}
	f.f.Write(m)
	f.Lock.Unlock()
}
