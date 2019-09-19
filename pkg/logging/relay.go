package logging

import (
	"fmt"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"path"
	"path/filepath"
	"strings"
)

type logRelay struct {
}

func (lr *logRelay) Write(p []byte) (n int, err error) {
	fullMessage := strings.TrimRight(string(p), "\n")
	colonBeforeLineNumberIndex := strings.Index(fullMessage, ":")
	if colonBeforeLineNumberIndex != -1 {
		colonBeforeMessageIndex := strings.Index(fullMessage[colonBeforeLineNumberIndex+1:], ":") + colonBeforeLineNumberIndex + 1
		if colonBeforeMessageIndex != -1 {
			lineNumberString := fullMessage[colonBeforeLineNumberIndex+1 : colonBeforeMessageIndex]
			logMessage := strings.TrimLeft(fullMessage[colonBeforeMessageIndex+1:], " ")
			filePath := filepath.ToSlash(fullMessage[:colonBeforeLineNumberIndex])
			interestingPath := path.Join(path.Base(path.Dir(path.Dir(filePath))), path.Base(path.Dir(filePath)), path.Base(filePath))
			codeLocationHint := colorize(fmt.Sprintf("[%s:%s]", interestingPath, lineNumberString), color.FgHiBlack, true)
			fullMessage = logMessage + " " + codeLocationHint
		}
	}
	log.Warn(fullMessage)
	return len(p), nil
}
