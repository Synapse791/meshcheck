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
	var testCases = []struct {
		Ip      string
		Port    string
	}{
		{"10.100.1.2",  "80"},
		{"192.168.1.1", "5000"},
		{"10.100.2.3",  "22"},
	}

	testString := fmt.Sprintf("%s:%s\n%s:%s\n%s:%s", testCases[0].Ip, testCases[0].Port, testCases[1].Ip, testCases[1].Port, testCases[2].Ip, testCases[2].Port)

	var testConfig AppConfig

	f, err := ioutil.TempFile("", "temp_connection_file")
	if err != nil { panic(err) }
	defer syscall.Unlink(f.Name())
	ioutil.WriteFile(f.Name(), []byte(testString), 0644)

	appErr := ReadConnectionConfig(f.Name(), &testConfig)

	if appErr != nil {
		t.Fatal(appErr.Error())
	}

	for caseNum, tc := range testCases {
		if testConfig.Connections[caseNum].ToAddress != tc.Ip {
			t.Errorf("expected %s to be %s", testConfig.Connections[caseNum].ToAddress, tc.Ip)
		}
		if testConfig.Connections[caseNum].Port != tc.Port {
			t.Errorf("expected %s to be %s", testConfig.Connections[caseNum].Port, tc.Port)
		}
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

}