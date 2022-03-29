package logging

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// HTTPLogPort is the TCP port for the HTTP log server
	HTTPLogPort = 3000

	// HTTPLogPath is the HTTP path for the log server
	HTTPLogPath = "/log"
)

const (
	// Debug defines the debug logging level.
	Debug = iota
	// Info defines the info logging level.
	Info
	// Warn defines the warn logging level.
	Warn
	// Error defines the error logging level.
	Error
	// Fatal defines the Fatal logging level.
	Fatal
)

type payload struct {
	Level string `json:"level"`
}

type fileLogLevel struct {
	fileName string
	level    int
}

type logging struct {
	logger      *zap.Logger
	atom        zap.AtomicLevel
	globalLevel int
	fileLog     fileLogLevel
}

var logInstance *logging
var lock sync.RWMutex

// getFunctionName
func getFunctionName() string {
	pc, _, _, _ := runtime.Caller(3)
	function := runtime.FuncForPC(pc)
	FuncSplit := strings.Split(function.Name(), ".")
	return FuncSplit[len(FuncSplit)-1]
}

// getPackageName
func getPackageName() string {
	pc, _, _, _ := runtime.Caller(3)
	functionFullName := runtime.FuncForPC(pc)
	function := filepath.Base(functionFullName.Name())
	FuncSplit := strings.Split(function, ".")
	return FuncSplit[0]
}

// getCaller
func getFileName() string {
	_, absFilename, _, _ := runtime.Caller(2)
	fileFullName := filepath.Base(absFilename)
	fileName := strings.TrimSuffix(fileFullName, filepath.Ext(fileFullName))
	return fileName
}

func logLevelToString(level int) string {
	switch level {
	case Debug:
		return "debug"
	case Info:
		return "info"
	case Warn:
		return "warn"
	case Error:
		return "error"
	case Fatal:
		return "fatal"
	}
	return "unknown"
}

// ServeHTTP handles the following HTTP requests:
// GET http://localhost:HTTPLogPort/log
// GET http://localhost:HTTPLogPort/log?file=<filename>
// PUT http://localhost:HTTPLogPort/log?file=<filename>&level=<level>
// PUT http://localhost:HTTPLogPort/log?level=<level>
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var m string
	var file, level string
	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, "bad syntax", http.StatusBadRequest)
		return
	}
	query := u.Query()
	file = query.Get("file")
	level = query.Get("level")
	enc := json.NewEncoder(w)
	switch r.Method {
	case http.MethodGet:
		if file == "" {
			// GET http://localhost:HTTPLogPort/log
			// GET http://localhost:HTTPLogPort/log?file=
			_ = enc.Encode(payload{Level: logLevelToString(logInstance.globalLevel)})
		} else {
			// GET http://localhost:HTTPLogPort/log?file=<filename>
			if logInstance.fileLog.fileName != file {
				http.Error(w, "no log set for filename", http.StatusBadRequest)
				return
			}

			_ = enc.Encode(payload{Level: logLevelToString(logInstance.fileLog.level)})
		}
	case http.MethodPost:
		if level == "" {
			http.Error(w, "bad syntax", http.StatusBadRequest)
			return
		}
		// PUT http://localhost:HTTPLogPort/log?file=<filename>&level=<level>
		// PUT http://localhost:HTTPLogPort/log?level=<level>
		levelOk := SetLogLevel(level, file)
		if !levelOk {
			http.Error(w, "Invalid log level", http.StatusBadRequest)
			return
		}
		_ = enc.Encode(payload{Level: level})
	default:
		http.Error(w, "Only GET and POST are supported.", http.StatusBadRequest)
		return
	}
	_, _ = w.Write([]byte(m))
}

// init implements zap log initialization. Enables Debug logs based on:
// - global log level
// - file specific log level: FILE_NAME_DEBUG, where 'FILE_NAME' is the name of
//   the file the init is called from.
//   Example: logs from file 'graph.go' require an env variable: GRAPH_DEBUG
func init() {
	if logInstance != nil {
		return
	}

	logger, err := zap.NewProduction(zap.AddCallerSkip(1), zap.AddStacktrace(zap.FatalLevel))

	atom := zap.NewAtomicLevel()
	atom.SetLevel(zap.DebugLevel)
	logger = logger.WithOptions(zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewConsoleEncoder(
				zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(os.Stdout),
			atom,
		)
	}))
	_ = logger.Sync()

	if err != nil {
		panic(err)
	}

	logInstance = &logging{
		logger:      logger,
		atom:        atom,
		globalLevel: Debug,
	}

	if os.Getenv("ENABLE_LOGGING_ENDPOINT") == "true" {
		StartDefaultEndpoint()
	}
}

// StartDefaultEndpoint enables the logging controller endpoint on
// on the default port on localhost, at the default path ("/log").
func StartDefaultEndpoint() *http.Server {
	return StartEndpoint(fmt.Sprintf(":%d", HTTPLogPort), "/log")
}

// StartEndpoint enables the logging controller endpoint on HTTPLogPort.
func StartEndpoint(addr, path string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(path, ServeHTTP)

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(http.ErrServerClosed, err) {
			Errorf("logging server failed: %v", err)
		}
	}()

	return srv
}

// SetLogLevel sets the logging level to be used.
func SetLogLevel(level string, file string) bool {
	level = strings.ToLower(level)
	var lvl int
	switch level {
	case "debug":
		lvl = Debug
	case "info":
		lvl = Info
	case "warn":
		lvl = Warn
	case "error":
		lvl = Error
	case "fatal":
		lvl = Fatal
	default:
		return false
	}
	if file != "" {
		lock.Lock()
		logInstance.fileLog.fileName = file
		logInstance.fileLog.level = lvl
		lock.Unlock()
	} else {
		lock.Lock()
		logInstance.globalLevel = lvl
		lock.Unlock()
	}
	return true
}

// func checkfileLogLevelMet check if file exists and level is met.
// Input:
// - file name
// - log level to check against
// Returns:
// - log level is set for input file
// - log level is met for input file
func checkfileLogLevelMet(file string, level int) (bool, bool) {
	lock.RLock()
	defer lock.RUnlock()
	if logInstance.fileLog.fileName == file {
		if logInstance.fileLog.level <= level {
			return true, true
		}
		return true, false
	}
	return false, false
}

func checkGlobalLevelMet(level int) bool {
	lock.RLock()
	defer lock.RUnlock()
	return logInstance.globalLevel <= level
}

func logFormat(msg string, args ...interface{}) string {
	s := fmt.Sprintf(msg, args...)
	m := fmt.Sprintf("%s  %s  %s", getPackageName(), getFunctionName(), s)
	return m
}

func checkLogLevel(file string, level int) bool {
	globalLvlMet := false
	fileFound, fileLvlMet := checkfileLogLevelMet(file, level)
	if !fileFound {
		globalLvlMet = checkGlobalLevelMet(level)
	}
	if fileLvlMet || globalLvlMet {
		return true
	}
	return false
}

// Debugf implements function to log at Debug level
func Debugf(msg string, args ...interface{}) {
	if checkLogLevel(getFileName(), Debug) {
		s := logFormat(msg, args...)
		logInstance.logger.Debug(s)
	}
}

// Infof implements function to log at Info level
func Infof(msg string, args ...interface{}) {
	if checkLogLevel(getFileName(), Info) {
		s := logFormat(msg, args...)
		logInstance.logger.Info(s)
	}
}

// Warnf implements function to log at Warn level
func Warnf(msg string, args ...interface{}) {
	if checkLogLevel(getFileName(), Warn) {
		s := logFormat(msg, args...)
		logInstance.logger.Warn(s)
	}
}

// Errorf implements function to log at Error level
func Errorf(msg string, args ...interface{}) {
	if checkLogLevel(getFileName(), Error) {
		s := logFormat(msg, args...)
		logInstance.logger.Error(s)
	}
}

// Fatalf implements function to log at fatal level
func Fatalf(msg string, args ...interface{}) {
	s := logFormat(msg, args...)
	logInstance.logger.Fatal(s)
}
