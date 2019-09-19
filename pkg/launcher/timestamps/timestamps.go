package timestamps

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

const timeFormat = "2006-01-02 15:04:05" // can be generated using the following unix command: date +"%Y-%m-%d %H:%M:%S"

type Timestamps struct {
	DeploymentConfig string            `json:"DeploymentConfig"`
	Bundles          map[string]string `json:"Bundles"`
}

func VerifyDeploymentConfigTimestamp(newTimestamp, filePath string) {
	timestamps := readTimestamps(filePath)
	timestamps.CheckAndSetDeploymentConfigTimestamp(newTimestamp)
	timestamps.write(filePath)
}

func VerifyBundleInfoTimestamp(uniqueBundleName, newTimestamp, filePath string) {
	timestamps := readTimestamps(filePath)
	timestamps.CheckAndSetBundleInfoTimestamp(uniqueBundleName, newTimestamp)
	timestamps.write(filePath)
}

func readTimestamps(filePath string) *Timestamps {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.WithFields(log.Fields{"err": err, "filePath": filePath}).Infof("Could not read timestamps " +
				"because the file does not yet exist: looks like this is the first run of the launcher.")
			return &Timestamps{DeploymentConfig: "", Bundles: map[string]string{}}
		}
		panic(err)
	}
	defer file.Close()

	return ReadTimestampsFromReader(file)
}

func (timestamps *Timestamps) write(filePath string) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(fmt.Sprintf("Could not open timestamps file \"%s\" for writing: %v", filePath, err))
	}
	defer file.Close()

	timestamps.WriteToWriter(file)
}

func ReadTimestampsFromReader(reader io.Reader) *Timestamps {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(fmt.Sprintf("Could not read from reader: %v", err))
	}

	var timestamps Timestamps
	err = json.Unmarshal(bytes, &timestamps)
	if err != nil {
		panic(fmt.Sprintf("Could not unmarshal timestamps `%s`: %v", string(bytes), err))
	}
	return &timestamps
}

func (timestamps *Timestamps) WriteToWriter(writer io.Writer) {
	json, err := json.Marshal(timestamps)
	if err != nil {
		panic(fmt.Sprintf("Could not marshal timestamps %+v: %v", timestamps, err))
	}
	_, err = writer.Write(json)
	if err != nil {
		panic(fmt.Sprintf("Could not write timestamps json: %v", err))
	}
}

func (timestamps *Timestamps) CheckAndSetDeploymentConfigTimestamp(newTimestampAsString string) {
	if timestamps.DeploymentConfig != "" {
		checkTimestamp(timestamps.DeploymentConfig, newTimestampAsString)
	} else {
		log.Info("No old timestamp found for deployment-config, seems that the launcher is started for the first time.")
	}

	timestamps.DeploymentConfig = newTimestampAsString
}

func (timestamps *Timestamps) CheckAndSetBundleInfoTimestamp(uniqueBundleName, newTimestampAsString string) {
	oldTimestampAsString, foundBundle := timestamps.Bundles[uniqueBundleName]
	if foundBundle {
		checkTimestamp(oldTimestampAsString, newTimestampAsString)
	} else {
		log.WithFields(log.Fields{"bundle": uniqueBundleName}).Info("No old timestamp found for bundle.")
	}

	timestamps.Bundles[uniqueBundleName] = newTimestampAsString
}

func checkTimestamp(oldTimestampAsString, newTimestampAsString string) {
	if oldTimestampAsString == "" {
		log.Info("No old timestamp found.")
	} else {
		oldTimestamp, err := time.Parse(timeFormat, oldTimestampAsString)
		if err != nil {
			panic(fmt.Sprintf("Could not parse old timestamp \"%s\": %v", oldTimestampAsString, err))
		}

		newTimestamp, err := time.Parse(timeFormat, newTimestampAsString)
		if err != nil {
			panic(fmt.Sprintf("Could not parse new timestamp \"%s\": %v", newTimestampAsString, err))
		}

		if newTimestamp.Before(oldTimestamp) {
			panic(fmt.Sprintf("New timestamp \"%s\" is older than old timestamp \"%s\". This may indicate an attack or a misconfiguration.", newTimestamp, oldTimestamp))
		}
	}
}
