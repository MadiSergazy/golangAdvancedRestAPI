package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

const (
	logsDir  = "logs"
	logsFile = "logs/all.log"
)

// kafka -- info debug
// file  -- error, trace
// stdout -- warning, critical
// Hooks in contexts of logrus //^ is used for attaching custom functions that will be executed when log entries of specific levels are made
// also they need for increasing perfomance, //? It specifies which writers to use for logs and which log levels to apply the hook to
type writeHook struct {
	Writer    []io.Writer    // * for being able to write to the: stdout stdr files elasticSearch kafka etc
	LogLevels []logrus.Level //for leveling logs
}

// this func will be called every time when we will write log
func (hook *writeHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}
	return nil
}

func (hook *writeHook) Levels() []logrus.Level {
	return hook.LogLevels
}

func init() { //when init writed with small letter it will called automaticly (if this package (package logging) will be used outside of the package)
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		//CallerPrettyfier is need for telling in whitch place we are logging
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}
	//* ModePerm => 0755 grants read, write, and execute permissions
	if err := os.MkdirAll(logsDir, os.ModePerm); err != nil {
		l.Info(err)
		panic(err)
	}
	//* 0644 grants read and write permissions
	//^os.O_CREATE: This flag indicates that the file should be created if it does not exist.
	//^os.O_WRONLY: This flag specifies that the file should be opened for writing.
	//^os.O_APPEND: This flag indicates that data should be appended to the end of the file rather than overwriting existing content
	allFile, err := os.OpenFile(logsFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		l.Info(err)
		panic(err)
	}

	// * sets the log output to be discarded, meaning logs won't be printed to the terminal.
	l.SetOutput(io.Discard) //need for log's did't go anywhere
	l.AddHook(&writeHook{
		Writer:    []io.Writer{allFile, os.Stdout}, //we will write to all.log file and terminal
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel) //capturing all log levels.

	e = logrus.NewEntry(l) // creates a new logrus Entry instance based on the configured logger
}

// * in some cases when you need to add new instance of logger into the application but also with additional field it might be needed
var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() Logger {
	return Logger{e}
}

func (l *Logger) GetLoggerWithField(field string, value interface{}) Logger {
	return Logger{l.WithField(field, value)}
}
