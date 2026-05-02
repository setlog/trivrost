package logging

import (
	"fmt"
	"maps"
	"path"
	"path/filepath"
	"runtime/debug"
	"slices"
	"strings"
	"time"

	"github.com/setlog/trivrost/pkg/misc"

	log "github.com/sirupsen/logrus"
)

const codeLocationHintColor = ansiHiBlack

type LogFormatter struct {
}

func (lf *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	if (entry.Message == "" || entry.Message == "\n") && isLevelUncritical(entry.Level) && len(entry.Data) == 0 {
		return []byte("\n"), nil
	}
	messageTypeShort := getMessageTypeShort(entry.Level)
	time := colorize(getTimeString(), ansiHiBlack, true)
	logMessage := misc.StripTrailingLineBreak(entry.Message)
	codeLocationHint := ""
	if !hasCodeLocationHint(logMessage) { // logRelay will have written code location hint for standard log package prints already.
		codeLocationHint = getCodeLocationHint(entry)
	}
	fields := getSortedFields(entry)
	logMessage, extraBreaks := misc.SplitTrailing(logMessage, "\n")
	printMessage := misc.JoinNonEmpty([]string{messageTypeShort, time, logMessage, codeLocationHint, fields}, " ")
	stackTrace := getStackTrace(entry)
	return []byte(printMessage + "\n" + extraBreaks + stackTrace), nil
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
		q := fmt.Sprintf("%s%s%v", colorize(fieldName, ansiCyan, false), colorize("=", ansiHiBlack, true), fieldValue)
		fields += q
		if i != lastIndex {
			fields += colorize(", ", ansiHiBlack, true)
		}
		i++
	}
	return fields
}

func getSortedFieldNames(data log.Fields) []string {
	return slices.Sorted(maps.Keys(data))
}

func getMessageTypeShort(level log.Level) string {
	switch level {
	case log.TraceLevel:
		return colorize("T", ansiHiBlack, true)
	case log.DebugLevel:
		return colorize("D", ansiHiBlack, true)
	case log.InfoLevel:
		return "I"
	case log.WarnLevel:
		return colorize("W", ansiYellow, true)
	case log.ErrorLevel:
		return colorize("E", ansiRed, true)
	case log.FatalLevel:
		return colorize("F", ansiRed, true)
	case log.PanicLevel:
		return colorize("P", ansiMagenta, true)
	default:
		return colorize("?", ansiRed, true)
	}
}

func isLevelUncritical(level log.Level) bool {
	return level == log.TraceLevel || level == log.DebugLevel || level == log.InfoLevel
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
