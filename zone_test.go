package main

import (
	"testing"
)

var (
	err error
)

func TestLoad(t *testing.T) {
	z := NewZone()
	if err = z.LoadFile("zones/test1.com.z"); err != nil {
		t.Fatal("load file failed:", err)
	}
}
