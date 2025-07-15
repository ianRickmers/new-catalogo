package utils

import (
	"log"
	"os"

	"github.com/natefinch/lumberjack"
)

func LoadLogFile(filepath string, filename string, maxSize int, maxBackups, maxAge int) {
	// MaxSize is the maximum size in megabytes of the log file before it gets
	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.
	// MaxBackups is the maximum number of old log files to retain.
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	log.SetOutput(&lumberjack.Logger{
		Filename:   filepath + filename + ".log",
		MaxSize:    maxSize,    // megabytes after which new file is created
		MaxBackups: maxBackups, // number of backups
		MaxAge:     maxAge,     //days
	})
}

func Debug(args ...any) {
	if os.Getenv("GIN_MODE") == "debug" {
		log.Print(args...)
	}
}
