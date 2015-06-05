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