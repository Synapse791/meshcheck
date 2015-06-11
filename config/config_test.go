package config

import (
	"testing"
	"fmt"
	"io/ioutil"
	"syscall"
)

func TestFixTrailingSlash(t *testing.T) {
	with    := "/tmp/"
	without := "/tmp"

	check_with      := FixTrailingSlash(with)
	check_without   := FixTrailingSlash(without)

	if check_with != "/tmp/" {
		t.Errorf("started with %s. Expected %s to be /tmp/", with, check_with)
	}

	if check_without != "/tmp/" {
		t.Errorf("started with %s. Expected %s to be /tmp/", without, check_without)
	}

}

func TestReadConnectionConfigWithFile(t *testing.T) {
	var testConfig AppConfig

	f, err := ioutil.TempFile("", "temp_connection_file")
	if err != nil { panic(err) }
	defer syscall.Unlink(f.Name())
	ioutil.WriteFile(f.Name(), []byte("10.100.1.2:80\n192.168.1.1:5000"), 0644)

	appErr := ReadConnectionConfig(f.Name(), &testConfig)

	if appErr != nil {
		t.Fatal(appErr.Error())
	}

	if testConfig.Connections[0].ToAddress != "10.100.1.2" || testConfig.Connections[0].Port != "80" {
		t.Error(fmt.Sprintf("expected %s to be 10.100.1.2 and %s to be 80", testConfig.Connections[0].ToAddress, testConfig.Connections[0].Port))
	}

	if testConfig.Connections[1].ToAddress != "192.168.1.1" || testConfig.Connections[1].Port != "5000" {
		t.Error(fmt.Sprintf("expected %s to be 192.168.1.1 and %s to be 5000", testConfig.Connections[1].ToAddress, testConfig.Connections[1].Port))
	}
}

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

func TestReadPortConfigInvalidMode(t *testing.T) {
	var testConfig AppConfig

	testConfig.Mode = "not_mode"

	err := ReadPortConfig("not_exist", &testConfig)

	if err == nil {
		t.Errorf("expected failure but got success")
	}

	if err != nil && err.Error() != "Unknown mode: not_mode" {
		t.Errorf("expected error message 'Unknown mode: not_mode' got '%s'", err.Error())
	}

	if err != nil && err.Error() == "Unknown mode: not_mode" {
		t.Errorf("expected error message 'Unknown mode: not_mode' got '%s'", err.Error())
	}

}