package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
)

func New() *logrus.Logger {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return "", fmt.Sprintf(" %s:%d", filename, f.Line)
		},
	}
	return l
}
