package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s [command]\n", os.Args[0])
		os.Exit(1)
	}

	cmdName := os.Args[1]
	if filepath.Base(os.Args[1]) == os.Args[1] {
		if lp, err := exec.LookPath(os.Args[1]); err != nil {
			fmt.Println("look path error:", err)
			os.Exit(1)
		} else {
			cmdName = lp
		}
	}

	procAttr := &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("look path error:", err)
		os.Exit(1)
	}
	fmt.Println(cwd)

	start := time.Now()
	process, err := os.StartProcess(cmdName, []string{cwd}, procAttr)
	if err != nil {
		fmt.Println("start process error:", err)
		os.Exit(2)
	}

	processState, err := process.Wait()
	if err != nil {
		fmt.Println("wait error:", err)
		os.Exit(3)
	}

	fmt.Println()
	fmt.Println("real", time.Now().Sub(start))
	fmt.Println("user", processState.UserTime())
	fmt.Println("system", processState.SystemTime())
}
