package logging

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/setlog/trivrost/pkg/misc"

	log "github.com/sirupsen/logrus"
)

const codeLocationHintColor = color.FgHiBlack

type LogFormatter struct {
}

func (lf *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	messageTypeShort := getMessageTypeShort(entry.Level)
	time := colorize(getTimeString(), color.FgHiBlack, true)
	logMessage := misc.StripTrailingLineBreak(entry.Message)
	codeLocationHint := ""
	if !hasCodeLocationHint(logMessage) { // logRelay will have written code location hint for standard log package prints already.
		codeLocationHint = getCodeLocationHint(entry)
	}
	fields := getSortedFields(entry)
	printMessage := misc.JoinNonEmpty([]string{messageTypeShort, time, logMessage, codeLocationHint, fields}, " ")
	stackTrace := getStackTrace(entry)
	return []byte(printMessage + "\n" + stackTrace), nil
}

func getTimeString() string {
	now := time.Now()
	timePassed := now.Sub(initTime)
	floatSecondsPassed := timePassed.Seconds()
	if floatSecondsPassed >= 999.5 {
		return fmt.Sprintf("%.1f", floatSecondsPassed) // Crazy long execution. Don't bother formatting this.
	}
	return fmt.Sprintf("%5.3f", floatSecondsPassed)[:5]
}

func hasCodeLocationHint(logMessage string) bool {
	openingMarkerText, closingMarkerText := formatStyle(codeLocationHintColor)+"[", "]"+clearStyle()
	openingMarkerIndex := strings.LastIndex(logMessage, openingMarkerText)
	closingMarkerIndex := strings.LastIndex(logMessage, closingMarkerText)
	if openingMarkerIndex != -1 && closingMarkerIndex != -1 && closingMarkerIndex > openingMarkerIndex {
		return closingMarkerIndex == len(logMessage)-len(closingMarkerText)
	}
	return false
}

func getCodeLocationHint(entry *log.Entry) string {
	if entry.HasCaller() {
		caller := entry.Caller
		// caller.Function is of the form "gitlab.example.com/project/package/subpackage/subpackage2.CalledFunction".
		// The last part can be "subpackage2.CalledFunction.func1.1" or such for anonymous functions. The path always uses forward slashes.
		// caller.File is of the form "/home/user/[...]/project/package/subpackage/subpackage2/file_containing_called_function.go",
		// and is OS-specific, i.e. it will contains backslashes if the binary was built under Windows.

		packageFuncSegments := strings.Split(filepath.Base(caller.Function), ".")
		fileName := path.Base(filepath.ToSlash(caller.File))
		if len(packageFuncSegments) > 1 {
			shortPackagePath := limitedPackagePath(path.Join(path.Dir(caller.Function), packageFuncSegments[0])) + "/"
			if shortPackagePath == "main/" {
				shortPackagePath = ""
			}
			functionName := packageFuncSegments[1]

			return colorize(fmt.Sprintf("[%s%s:%s():%d]", shortPackagePath, fileName, functionName, caller.Line), codeLocationHintColor, true)
		}
		return colorize(fmt.Sprintf("[%s:%d]", fileName, caller.Line), codeLocationHintColor, true)
	}
	return ""
}

func getSortedFields(entry *log.Entry) string {
	fields := ""
	sortedFieldNames := getSortedFieldNames(entry.Data)
	lastIndex, i := len(entry.Data)-1, 0
	for _, fieldName := range sortedFieldNames {
		fieldValue := entry.Data[fieldName]
		q := fmt.Sprintf("%s%s%v", colorize(fieldName, color.FgCyan, false), colorize("=", color.FgHiBlack, true), fieldValue)
		fields += q
		if i != lastIndex {
			fields += colorize(", ", color.FgHiBlack, true)
		}
		i++
	}
	return fields
}

func getSortedFieldNames(data log.Fields) []string {
	names, i := make([]string, len(data)), 0
	for k := range data {
		names[i] = k
		i++
	}
	sort.Strings(names)
	return names
}

func getMessageTypeShort(level log.Level) string {
	switch level {
	case log.TraceLevel:
		return colorize("T", color.FgHiBlack, true)
	case log.DebugLevel:
		return colorize("D", color.FgHiBlack, true)
	case log.InfoLevel:
		return "I"
	case log.WarnLevel:
		return colorize("W", color.FgYellow, true)
	case log.ErrorLevel:
		return colorize("E", color.FgRed, true)
	case log.FatalLevel:
		return colorize("F", color.FgRed, true)
	case log.PanicLevel:
		return colorize("P", color.FgMagenta, true)
	default:
		return colorize("?", color.FgRed, true)
	}
}

func getStackTrace(entry *log.Entry) string {
	stackTrace := ""
	if entry.Level == log.TraceLevel {
		stackTrace = string(debug.Stack())
	}
	return stackTrace
}

func limitedPackagePath(packagePath string) string {
	segments := strings.Split(packagePath, "/")
	segmentCount := len(segments)

	const maxSegments int = 2
	if segmentCount > maxSegments {
		return path.Join(segments[segmentCount-maxSegments:]...)
	}
	return path.Join(segments...)
}
