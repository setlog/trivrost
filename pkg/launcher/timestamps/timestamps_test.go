package timestamps_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/setlog/trivrost/pkg/launcher/timestamps"
)

const timestampsContent string = `{
		"DeploymentConfig": "2012-01-05 17:22:49",
		"Bundles": {
		"testapp": "2017-04-03 09:07:59",
		"runtime": "2017-09-22 14:23:06"
	}
}`

const testAppBundle = "testapp"

func TestReadTimestampsFromReader(t *testing.T) {
	reader := strings.NewReader(timestampsContent)
	timestamps := timestamps.ReadTimestampsFromReader(reader)
	tests := []struct {
		value, valueName, expected string
	}{
		{timestamps.DeploymentConfig, "DeploymentConfig", "2012-01-05 17:22:49"},
		{timestamps.Bundles["testapp"], "Bundles[\"testapp\"]", "2017-04-03 09:07:59"},
		{timestamps.Bundles["runtime"], "Bundles[\"runtime\"]", "2017-09-22 14:23:06"},
	}
	for _, test := range tests {
		if test.value != test.expected {
			t.Errorf("TestReadTimestampsFromReader test failed for %s. Got: %s; Expected: %s.", test.valueName, test.value, test.expected)
		}
	}
}

func TestWriteTimestampsToWriter(t *testing.T) {
	timestamps := timestamps.Timestamps{DeploymentConfig: "2012-01-05 17:22:49", Bundles: map[string]string{"testapp": "2017-04-03 09:07:59", "runtime": "2017-09-22 14:23:06"}}
	writer := new(bytes.Buffer)
	timestamps.WriteToWriter(writer)
	result := writer.String()
	if !areJsonEqual(timestampsContent, result, t) {
		t.Errorf("The generated json is not correct.\nGot:\n%s\nExpected:\n%s", result, timestampsContent)
	}
}

func areJsonEqual(str1, str2 string, t *testing.T) bool {
	return reflect.DeepEqual(parseAbstractJson(str1, t), parseAbstractJson(str2, t))
}

func parseAbstractJson(str string, t *testing.T) interface{} {
	var parsed interface{}
	err := json.Unmarshal([]byte(str), &parsed)
	if err != nil {
		t.Fatalf("The following string doesn't seem to be a valid json: %s", str)
	}
	return parsed
}

func TestCheckAndSetDeploymentConfigTimestampWithNewVersion(t *testing.T) {
	reader := strings.NewReader(timestampsContent)
	timestamps := timestamps.ReadTimestampsFromReader(reader)
	newTimestamp := "2012-01-05 17:22:50"
	timestamps.CheckAndSetDeploymentConfigTimestamp(newTimestamp)
	if timestamps.DeploymentConfig != newTimestamp {
		t.Errorf("The timestamp of the deployment-config was not set correctly. Got: %s; Expected: %s.", timestamps.DeploymentConfig, newTimestamp)
	}
}

func TestCheckAndSetDeploymentConfigTimestampWithSameVersion(t *testing.T) {
	reader := strings.NewReader(timestampsContent)
	timestamps := timestamps.ReadTimestampsFromReader(reader)
	newTimestamp := "2012-01-05 17:22:49"
	timestamps.CheckAndSetDeploymentConfigTimestamp(newTimestamp)
	if timestamps.DeploymentConfig != newTimestamp {
		t.Errorf("The timestamp of the deployment-config was not set correctly. Got: %s; Expected: %s.", timestamps.DeploymentConfig, newTimestamp)
	}
}

func TestCheckAndSetDeploymentConfigTimestampWithOldVersion(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The new timestamp of the deployment-config was before the old timestamp, but the code did not panic as expected.")
		}
	}()

	reader := strings.NewReader(timestampsContent)
	timestamps := timestamps.ReadTimestampsFromReader(reader)
	newTimestamp := "2012-01-05 17:22:48"
	timestamps.CheckAndSetDeploymentConfigTimestamp(newTimestamp)
}

func TestCheckAndSetDeploymentConfigTimestampOnFirstStartup(t *testing.T) {
	timestamps := timestamps.Timestamps{DeploymentConfig: "", Bundles: nil}
	newTimestamp := "2012-01-05 17:22:49"
	timestamps.CheckAndSetDeploymentConfigTimestamp(newTimestamp)
	if timestamps.DeploymentConfig != newTimestamp {
		t.Errorf("The timestamp of the deployment-config was not set correctly. Got: %s; Expected: %s.", timestamps.DeploymentConfig, newTimestamp)
	}
}

func TestCheckAndSetBundleInfoTimestampWithNewVersion(t *testing.T) {
	reader := strings.NewReader(timestampsContent)
	timestamps := timestamps.ReadTimestampsFromReader(reader)
	newTimestamp := "2017-04-03 09:08:00"
	timestamps.CheckAndSetBundleInfoTimestamp(testAppBundle, newTimestamp)
	if timestamps.Bundles[testAppBundle] != newTimestamp {
		t.Errorf("The timestamp of the bundle was not set correctly. Got: %s; Expected: %s.", timestamps.DeploymentConfig, newTimestamp)
	}
}

func TestCheckAndSetBundleInfoTimestampWithSameVersion(t *testing.T) {
	reader := strings.NewReader(timestampsContent)
	timestamps := timestamps.ReadTimestampsFromReader(reader)
	newTimestamp := "2017-04-03 09:07:59"
	timestamps.CheckAndSetBundleInfoTimestamp(testAppBundle, newTimestamp)
	if timestamps.Bundles[testAppBundle] != newTimestamp {
		t.Errorf("The timestamp of the bundle was not set correctly. Got: %s; Expected: %s.", timestamps.DeploymentConfig, newTimestamp)
	}
}

func TestCheckAndSetBundleInfoTimestampWithOldVersion(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The new timestamp of the bundle was before the old timestamp, but the code did not panic as expected.")
		}
	}()

	reader := strings.NewReader(timestampsContent)
	timestamps := timestamps.ReadTimestampsFromReader(reader)
	newTimestamp := "2012-01-05 17:22:47"
	timestamps.CheckAndSetBundleInfoTimestamp(testAppBundle, newTimestamp)
}

func TestCheckAndSetBundleInfoTimestampWithNewBundle(t *testing.T) {
	reader := strings.NewReader(timestampsContent)
	timestamps := timestamps.ReadTimestampsFromReader(reader)
	bundle := "new_bundle"
	newTimestamp := "2019-02-03 16:34:16"
	timestamps.CheckAndSetBundleInfoTimestamp(bundle, newTimestamp)
	if timestamps.Bundles[bundle] != newTimestamp {
		t.Errorf("The timestamp of the bundle was not set correctly. Got: %s; Expected: %s.", timestamps.DeploymentConfig, newTimestamp)
	}
}
