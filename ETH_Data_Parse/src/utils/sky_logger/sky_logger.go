package sky_logger

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

const (
	DebugLevel int = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

type LogFile struct {
	level     int
	log_time  int64
	log_path  string
	file_name string
	
	file_fd *os.File
}

var sky_log LogFile

func InitLogger(log_path string, level int) error {
	sky_log.log_path = log_path
	sky_log.file_name = fmt.Sprintf("%s/skydp_log_%d", log_path, os.Getpid())
	sky_log.level = level
	sky_log.file_fd = nil

	fmt.Printf("Init logger module, log_path is %s\n", sky_log.log_path)
	err := os.MkdirAll(sky_log.log_path, os.ModePerm)
	if err != nil {
		fmt.Printf("Fail to mkdir for log_path(%s), err_desc=%s\n",
			sky_log.log_path, err.Error())
		return err
	}

	fmt.Sprintf("Success to mkdir for log_path(%s)\n", sky_log.log_path)

	sky_log.create_log_file()

	log.SetOutput(&sky_log)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	return nil
}

func SetLevel(level int) {
	sky_log.level = level
}

func GetLevel() int {
	return sky_log.level
}

func Debugf(format string, args ...interface{}) {
	if sky_log.level <= DebugLevel {
		log.SetPrefix("[Debug] ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Infof(format string, args ...interface{}) {
	if sky_log.level <= InfoLevel {
		log.SetPrefix("[Info ] ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Warnf(format string, args ...interface{}) {
	if sky_log.level <= WarnLevel {
		log.SetPrefix("[Warn ] ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Errorf(format string, args ...interface{}) {
	if sky_log.level <= ErrorLevel {
		log.SetPrefix("[Error] ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Fatalf(format string, args ...interface{}) {
	if sky_log.level <= FatalLevel {
		log.SetPrefix("[Fatal] ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func (me *LogFile) Write(buf []byte) (n int, err error) {
	if me.file_name == "" {
		return len(buf), nil
	}

	if me.log_time+3600 < time.Now().Unix() || me.file_fd == nil {
		me.create_log_file()
		me.log_time = time.Now().Unix()
	}
	if me.file_fd == nil {
		return len(buf), nil
	}

	return me.file_fd.Write(buf)
}

func (sky_log *LogFile) create_log_file() {

	var err error

	_, err = os.Stat(sky_log.file_name)
	if err == nil {
		now := time.Now()
		filename := fmt.Sprintf("%s_%04d%02d%02d_%02d%02d", sky_log.file_name, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
		err = os.Rename(sky_log.file_name, filename)

		if err == nil {
			go func() {
				tarCmd := exec.Command("tar", "-zcf", filename+".tar.gz", filename, "--remove-files")
				tarCmd.Run()

				rmCmd := exec.Command("/bin/sh", "-c", "find "+sky_log.log_path+` -type f -mtime +6 -exec rm {} \;`)
				rmCmd.Run()
			}()
		} else {
			Errorf("Fail to rename (%s) to (%s)", sky_log.file_name, filename)
		}
	}

	for index := 0; index < 10; index++ {
		fd, err := os.OpenFile(sky_log.file_name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
		if nil == err {
			sky_log.file_fd.Sync()
			sky_log.file_fd.Close()
			sky_log.file_fd = fd
			break
		}

		time.Sleep(10 * time.Millisecond)

		sky_log.file_fd = nil
	}
}
