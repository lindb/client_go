package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lindb/lindb/pkg/fileutil"
)

func TestLoadConfig(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "client_config.toml")

	var err error

	// not exists
	if fileutil.Exist(tmpFile) {
		if err := os.Remove(tmpFile); err != nil {
			t.Fatal(err)
		}
	}

	_, err = LoadConfig(tmpFile)
	if err == nil {
		t.Fatal("should be error")
	}

	// empty file
	err = fileutil.EncodeToml(tmpFile, &ClientConfig{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = LoadConfig(tmpFile)
	if err == nil {
		t.Fatal("should be error")
	}

	// no broker url
	err = fileutil.EncodeToml(tmpFile, &ClientConfig{BrokerURL: ""})
	if err != nil {
		t.Fatal(err)
	}

	_, err = LoadConfig(tmpFile)
	if err == nil {
		t.Fatal("should be error")
	}

	// normal
	err = fileutil.EncodeToml(tmpFile, &ClientConfig{BrokerURL: "http://127.0.0.1:9000"})
	if err != nil {
		t.Fatal(err)
	}

	_, err = LoadConfig(tmpFile)
	if err != nil {
		t.Fatal("should be error")
	}
}
