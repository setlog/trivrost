package locking

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/setlog/trivrost/pkg/system"
)

func readProcessSignatureListFile(filePath string) (procSigs []system.ProcessSignature) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		panic(fmt.Sprintf("Could not read process signature file: %v", err))
	}
	return unmarshalProcessSignaturesLeniently(bytes)
}

func unmarshalProcessSignaturesLeniently(data []byte) (procSigs []system.ProcessSignature) {
	err := json.Unmarshal(data, &procSigs)
	if err != nil {
		// On Windows, sometimes files break and contain NULL bytes. In that case, ignore.
		log.Warnf("Could not unmarshal process signature json: %v", err)
		return nil
	}
	return procSigs
}

func mustWriteProcessSignatureListFile(filePath string, procSigs []system.ProcessSignature) {
	bytes, err := json.Marshal(procSigs)
	if err != nil {
		panic(fmt.Sprintf("Could not marshal process signature slice of length %d: %v", len(procSigs), err))
	}
	err = ioutil.WriteFile(filePath, bytes, 0666)
	if err != nil {
		panic(fmt.Sprintf("Could not write process signature list file \"%s\": %v", filePath, err))
	}
}
