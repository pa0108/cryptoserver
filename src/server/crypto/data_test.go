package crypto

import (
	"crypto/src/crypto"
	"testing"
)

func TestGetData(t *testing.T) {
	client := crypto.NewClient("", "")
	client.GetData()
	if len(outData) == 0 {
		t.Errorf("Failed to receive data")
	}
	return
}

func TestServeHTTP(t *testing.T) {
	return
}
