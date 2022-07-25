package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func init() {
	// start at the bottom of the stack before the package-name cache is primed
	minimumCallerDepth = 1

	logrus.SetFormatter(&MyFormatter{})
}

var (
	// qualified package name, cached at first use
	packageName string

	// Positions in the call stack when tracing to report the calling method
	minimumCallerDepth int

	// Used for caller information initialisation
	callerInitOnce sync.Once
)

const (
	maximumCallerDepth int = 25
	knownLogrusFrames  int = 4
)

type MyFormatter struct {
}

type myFields logrus.Fields

func (m myFields) String() string {
	if len(m) == 0 {
		return ""
	}
	b, _ := json.Marshal(m)
	return string(b)
}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog string

	caller := getCaller()

	var (
		funcName string
		file     string
		line     int
	)
	if caller != nil {
		s := strings.Split(filepath.Base(caller.Function), ".")
		if len(s) != 0 {
			funcName = s[len(s)-1]
		}
		file = filepath.Base(caller.File)
		line = caller.Line
	} else {
		funcName = "?"
		file = "?"
		line = 0
	}

	newLog = fmt.Sprintf("[%s][%s][%s->%s:%d]%s %s\n",
		timestamp, entry.Level, file, funcName, line, myFields(entry.Data), entry.Message)

	b.WriteString(newLog)
	return b.Bytes(), nil
}

// getCaller retrieves the name of the first log calling function
//copied from logrus
func getCaller() *runtime.Frame {
	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(0, pcs)

		// dynamic get the package name and the minimum caller depth
		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "getCaller") {
				packageName = getPackageName(funcName)
				break
			}
		}

		minimumCallerDepth = knownLogrusFrames
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != packageName {
			//return &f //nolint:scopelint
			fmt.Println(pkg, f.Function, f.Line)
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

// getPackageName reduces a fully qualified function name to the package name
// There really ought to be a better way...
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return logrus.WithField(key, value)
}
