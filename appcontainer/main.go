// app2.go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	for {
		cmd := exec.Command("ps", "-e", "-o", "pid,comm")
		cmd.Stdout = os.Stdout
		cmd.Run()
		fmt.Println()
		time.Sleep(10 * time.Second)
	}
}
