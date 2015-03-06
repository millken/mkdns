package main

import (
	"testing"
)

var (
	err error
)

func TestLoad(t *testing.T) {
	z := NewZone("test")
	if err = z.LoadFile("zones/test1.com.z"); err != nil {
		t.Fatal("load file failed:", err)
	}
}
