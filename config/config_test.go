package config

import (
	"testing"
	"fmt"
	"io/ioutil"
	"syscall"
)

func TestReadPortConfigWithFile(t *testing.T) {
	var testConfig AppConfig

	f, err := ioutil.TempFile("", "temp_port_file")
	if err != nil { panic(err) }
	defer syscall.Unlink(f.Name())
	ioutil.WriteFile(f.Name(), []byte("9000"), 0644)

	ReadPortConfig(f.Name(), &testConfig)

	if testConfig.Port != ":9000" {
		t.Error(fmt.Sprintf("expected %s to be :9000", testConfig.Port))
	}

}

func TestReadPortConfigServerWithoutFile(t *testing.T) {
	var testConfig AppConfig

	testConfig.Mode = "server"

	ReadPortConfig("not_exist", &testConfig)

	if testConfig.Port != ":6800" {
		t.Error(fmt.Sprintf("expected %s to be :6800", testConfig.Port))
	}

}

func TestReadPortConfigClientWithoutFile(t *testing.T) {
	var testConfig AppConfig

	testConfig.Mode = "client"

	ReadPortConfig("not_exist", &testConfig)

	if testConfig.Port != ":6600" {
		t.Error(fmt.Sprintf("expected %s to be :6600", testConfig.Port))
	}

}