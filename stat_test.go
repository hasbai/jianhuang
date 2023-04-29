package main

import (
	"fmt"
	"testing"
)

func TestStat(t *testing.T) {
	fmt.Println(readableBytes(3 * 1024 * 1024))
}
