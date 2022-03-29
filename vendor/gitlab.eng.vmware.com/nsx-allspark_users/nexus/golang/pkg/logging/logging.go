package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// HTTPLogPort is the TCP port for the HTTP log server
const HTTPLogPort = 3000

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

type logging struct {
	logger      *zap.Logger
	atom        zap.AtomicLevel
	globalLevel int
	fileLog     map[string]int // [fileName]level
}

var logInstance *logging

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

/*
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
			enc.Encode(payload{Level: logLevelToString(logInstance.globalLevel)})
		} else {
			// GET http://localhost:HTTPLogPort/log?file=<filename>
			if logInstance.fileLog.fileName != file {
				http.Error(w, "no log set for filename", http.StatusBadRequest)
				return
			}

			enc.Encode(payload{Level: logLevelToString(logInstance.fileLog.level)})
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
		enc.Encode(payload{Level: level})
	default:
		http.Error(w, "Only GET and POST are supported.", http.StatusBadRequest)
		return
	}
	w.Write([]byte(m))
}
*/
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
	if err != nil {
		panic(err)
	}

	atom := zap.NewAtomicLevel()
	// atom.SetLevel(zap.InfoLevel) // .DebugLevel)
	atom.SetLevel(zap.DebugLevel)
	logger = logger.WithOptions(zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewConsoleEncoder(
				zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(os.Stdout),
			atom,
		)
	}))
	logger.Sync()
	if err != nil {
		panic(err)
	}

	logInstance = &logging{
		logger:      logger,
		atom:        atom,
		globalLevel: Info, // Debug, // Info, // Debug,
		fileLog:     make(map[string]int),
	}
	goDebug := os.Getenv("GODEBUG")
	goDebugSplit := strings.Split(goDebug, ",")
	for _, itm := range goDebugSplit {
		debugSettingItem := strings.TrimSpace(itm)
		debugSettings := strings.Split(debugSettingItem, ":")
		if len(debugSettings) == 2 {
			SetLogLevel(debugSettings[1], debugSettings[0])
		} else {
			SetLogLevel(debugSettingItem, "")
		}
	}
	// SetLogLevel("debug", "dataModelCache")

	/*
			http.HandleFunc("/log", ServeHTTP)
			httpListenPort := fmt.Sprintf(":%d", HTTPLogPort)
		    go http.ListenAndServe(httpListenPort, nil)
	*/
}

// SetLogLevel sets the logging level to be used.
func SetLogLevel(level string, file string) bool {
	level = strings.ToLower(level)
	lvl := Debug
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
		logInstance.fileLog[file] = lvl
	} else {
		logInstance.globalLevel = lvl
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
	if fileLevel, ok := logInstance.fileLog[file]; ok {
		if fileLevel <= level {
			return true, true
		}
		return true, false
	}
	return false, false
}

func checkGlobalLevelMet(level int) bool {
	lvl := logInstance.globalLevel
	if lvl > level {
		return false
	}
	return true
}

func logFormat(msg string, args ...interface{}) string {
	s := fmt.Sprintf(msg, args...)
	m := fmt.Sprintf("%s  %s  %s", getPackageName(), getFunctionName(), s)
	return m
}

func checkLogLevel(file string, level int) bool {
	// fmt.Printf("checklog [%s]\n", file)
	// if file == "dmBus" {
	// 	return true
	// }
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
