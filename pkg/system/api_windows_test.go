package system

import (
	"reflect"
	"testing"
)

func TestRemoveEnv(t *testing.T) {
	envs := []string{"A=Foo", "B=Bar", "C=Delta"}
	envs = removeEnv(envs, "b")
	if !reflect.DeepEqual(envs, []string{"A=Foo", "C=Delta"}) {
		t.Errorf("Expected %v, but got %v.", []string{"A=Foo", "C=Delta"}, envs)
	}
	envs = removeEnv(envs, "a")
	if !reflect.DeepEqual(envs, []string{"C=Delta"}) {
		t.Errorf("Expected %v, but got %v.", []string{"C=Delta"}, envs)
	}
	envs = removeEnv(envs, "c")
	if !reflect.DeepEqual(envs, []string{}) {
		t.Errorf("Expected %v, but got %v.", []string{}, envs)
	}
}
