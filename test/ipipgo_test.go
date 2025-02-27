package test

import (
	"crypto/md5"
	"testing"
)

func TestIpipgo(t *testing.T) {
	hash := md5.New()
	hash.Write([]byte("test"))
	// ipipgo.CreateAccount(hex.EncodeToString(hash.Sum(nil)))
	t.Log("TestIpipgo")
}
