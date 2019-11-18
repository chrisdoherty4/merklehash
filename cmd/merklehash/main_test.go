package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

var tempDir string
var testDir []string

func TestMain(m *testing.M) {
	// Generate test directories
	tempDir, err := ioutil.TempDir("", "merklehash-")
	if err != nil {
		log.Fatal(err)
	}

	testDir = make([]string, 3)
	content := []byte("Hello World!")

	for i := 1; i <= 3; i++ {
		testDir[i] = filepath.Join(tempDir, "test-"+strconv.Itoa(i))

		if err := os.Mkdir(testDir[i], os.ModeDir); err != nil {
			os.RemoveAll(tempDir)
			log.Fatal(err)
		}

		for j := 1; j <= i; j++ {
			if err := ioutil.WriteFile(filepath.Join(testDir[i], "test-"+strconv.Itoa(j)+".txt"), content, 0755); err != nil {
				os.RemoveAll(tempDir)
				log.Fatal(err)
			}
		}
	}

	os.Exit(m.Run())

	// /os.RemoveAll(tempDir)
}

func TestCli(t *testing.T) {
	fmt.Println("Hello world!")
}
