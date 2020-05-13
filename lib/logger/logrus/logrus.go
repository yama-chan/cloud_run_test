package logrus

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/taisukeyamashita/test/lib/logger"
)

type logrusLogger struct {
	logger *logrus.Logger
}

var _ logger.Logger = logrusLogger{}

// NewLogger 新規Loggerを生成
func NewLogger() logger.Logger {

	/*
		PanicLevel Level = iota
			// PanicLevel level, highest level of severity. Logs and then calls panic with the
			// message passed to Debug, Info, ...

		FatalLevel
			// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
			// logging level is set to Panic.

		ErrorLevel
			// ErrorLevel level. Logs. Used for errors that should definitely be noted.
			// Commonly used for hooks to send errors to an error tracking service.

		WarnLevel
			// WarnLevel level. Non-critical entries that deserve eyes.

		InfoLevel
			// InfoLevel level. General operational entries about what's going on inside the
			// application.

		DebugLevel
			// DebugLevel level. Usually only enabled when debugging. Very verbose logging.

		TraceLevel
			// TraceLevel level. Designates finer-grained informational events than the Debug.
	*/
	logrusLogger := logrusLogger{
		//ログの記載を'Lock'するためにsync.Mutexをフィールドに持っているためポインタ(&)で渡すこと
		logger: &logrus.Logger{
			Out:          os.Stdout,
			Formatter:    new(logrus.TextFormatter),
			Hooks:        make(logrus.LevelHooks),
			Level:        logrus.DebugLevel,
			ExitFunc:     os.Exit,
			ReportCaller: true, //Caller情報をログに記録するかどうかのフラグ（デフォルトではオフ）
		},
	}
	return logrusLogger
}

// Debug Debugレベルのログを出力する関数
func (l logrusLogger) Debug(message string) {
	l.logger.Debugf("Caller: [%s]\nMessage: %s", getCaller(), message)
}

// Info Infoレベルのログを出力する関数
func (l logrusLogger) Info(message string) {
	l.logger.Infof("Caller: [%s]\nMessage: %s", getCaller(), message)
}

// Warn Warnレベルのログを出力する関数
func (l logrusLogger) Warn(message string) {
	l.logger.Warnf("Caller: [%s]\nMessage: %s", getCaller(), message)
}

// Error Errorレベルのログを出力する関数
func (l logrusLogger) Error(message string) {
	l.logger.Errorf("Caller: [%s]\nMessage: %s", getCaller(), message)
}

// Fatal Fatalレベルのログを出力する関数
func (l logrusLogger) Fatal(message string) {
	l.logger.Fatalf("Caller: [%s]\nMessage: %s", getCaller(), message)
}

// Panic Panicレベルのログを出力する関数
func (l logrusLogger) Panic(message string) {
	l.logger.Panicf("Caller: [%s]\nMessage: %s", getCaller(), message)
}

func getCaller() string {
	pc, file, line, _ := runtime.Caller(2)

	fmt.Printf("pc : %v", pc)
	fmt.Printf("file : %v", file)
	fmt.Printf("line : %v", line)

	// var (
	// 	f        = runtime.FuncForPC(pc)
	// 	parts    = strings.Split(f.Name(), "/")
	// 	funcName = parts[len(parts)-1]
	// )

	f := runtime.FuncForPC(pc)
	fnName := f.Name()
	splitedFnName := strings.Split(fnName, ".")
	packageName := splitedFnName[0]
	callerFuncName := splitedFnName[1]

	log.Printf("f: %v\n", f)
	log.Printf("packageName: %s\n", packageName)
	log.Printf("functionName: %s\n", callerFuncName)

	return fmt.Sprintf("file=%s、pkg=%s、func=%s、line=%v", file, packageName, callerFuncName, line)
	// return fmt.Sprintf("%s\n#%s:L%v", file, funcName, line)
}
