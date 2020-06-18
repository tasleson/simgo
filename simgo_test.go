package main

import (
	"fmt"
	"os"
	"testing"
)

// TestMain ...
func TestMain(t *testing.T) {
	fmt.Printf("TestMain %s\n", os.Getenv("LSM_GO_FD"))
	main()
}
