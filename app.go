package main

import (
	"os"
	"os/exec"
	"slices"
	"strings"

	"gobrago/lib"
)

func main() {
	argsWithoutProg := os.Args[1:]
	runIdx := slices.Index(argsWithoutProg, "--run")
	run := runIdx != -1
	if run {
		argsWithoutProg = slices.Delete(argsWithoutProg, runIdx, runIdx+1)
	}
	if len(argsWithoutProg) != 2 {
		panic("Incorrect number of arguments passed. Expected two: the Gobra cfg file and the verification job file")
	}
	installCfgPath := argsWithoutProg[0]
	jobCfgPath := argsWithoutProg[1]
	cmd, err := lib.GenCmd(installCfgPath, jobCfgPath)
	if err != nil {
		panic(err)
	}
	println(cmd)
	if run {
		cmdParts := strings.Split(cmd, " ")
		exeCmd := exec.Command(cmdParts[0], cmdParts[1:]...)
		exeCmd.Stdout = os.Stdout
		exeCmd.Stderr = os.Stderr
		err := exeCmd.Run()
		if err != nil {
			panic(err)
		}
	}
}
