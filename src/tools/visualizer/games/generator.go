package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var (
	wg sync.WaitGroup
)

func worker(i int, url string, seed int64) {
	c, err := exec.Command("../../../distribution/ic20_linux", "-u", url, "-s", fmt.Sprint(seed)).Output()
	if err != nil {
		log.Fatalf("failed to run command line tool: %s", err)
	}
	bs := strings.Split(string(c), "\n")
	of, err := os.Create(fmt.Sprintf("%d.json", i))
	if err != nil {
		log.Fatalf("failed to create output file: %s", err)
	}
	defer of.Close()
	of.WriteString("[")
	for i := 0; i < len(bs); i += 2 {
		of.WriteString(bs[i])
		if i != len(bs)-2 {
			of.WriteString(",")
		}
	}
	of.WriteString("]")
	wg.Done()
}

func main() {
	us := []string{
		"http://localhost:50123",
		"http://localhost:50123",
		"http://localhost:50123",
		"http://localhost:50123",
	}

	for i, u := range us {
		wg.Add(1)
		go worker(i, u, 1)
	}

	wg.Wait()
}
