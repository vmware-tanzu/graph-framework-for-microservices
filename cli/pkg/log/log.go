package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger = logrus.New()

	Panicf   = log.Panicf
	Fatalf   = log.Fatalf
	Errorf   = log.Errorf
	Warnf    = log.Warnf
	Infof    = log.Infof
	Debugf   = log.Debugf
	Tracef   = log.Tracef
	SetLevel = log.SetLevel

	PanicLevel = logrus.PanicLevel
	FatalLevel = logrus.FatalLevel
	ErrorLevel = logrus.ErrorLevel
	WarnLevel  = logrus.WarnLevel
	InfoLevel  = logrus.InfoLevel
	DebugLevel = logrus.DebugLevel
	TraceLevel = logrus.TraceLevel
)

func init() {
	log.Level = logrus.ErrorLevel
	log.Formatter.(*logrus.TextFormatter).DisableColors = true
	log.Out = os.Stdout
}
