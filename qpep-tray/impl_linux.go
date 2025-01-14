package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/parvit/qpep/shared"
)

const (
	BASEDIR_ENVIRONMENTVAR = "HOME"
	EXENAME                = "qpep"
)

func getClientCommand() *exec.Cmd {
	exeFile := filepath.Join(ExeDir, EXENAME)

	var verboseFlag = ""
	if qpepConfig.Verbose {
		verboseFlag = "--verbose"
	}

	cmd := exec.Command(exeFile,
		"--client",
		fmt.Sprintf("%v", verboseFlag),
		"--acks", fmt.Sprintf("%d", qpepConfig.Acks),
		"--ackDelay", fmt.Sprintf("%d", qpepConfig.AckDelay),
		"--congestion", fmt.Sprintf("%d", qpepConfig.Congestion),
		"--decimate", fmt.Sprintf("%d", qpepConfig.Decimate),
		"--minBeforeDecimation", fmt.Sprintf("%d", qpepConfig.DelayDecimate),
		"--gateway", fmt.Sprintf("%s", qpepConfig.GatewayHost),
		"--port", fmt.Sprintf("%d", qpepConfig.GatewayPort),
		"--listenaddress", fmt.Sprintf("%s", qpepConfig.ListenHost),
		"--listenport", fmt.Sprintf("%d", qpepConfig.ListenPort),
		"--multistream", fmt.Sprintf("%v", qpepConfig.MultiStream),
		"--varAckDelay", fmt.Sprintf("%d", qpepConfig.VarAckDelay),
		"--threads", fmt.Sprintf("%d", qpepConfig.WinDivertThreads))

	if cmd == nil {
		ErrorMsg("Could not create client command")
		return nil
	}
	cmd.Dir, _ = filepath.Abs(ExeDir)

	log.Println(cmd.Path)
	log.Println(cmd.Dir)
	return cmd
}

func getServerCommand() *exec.Cmd {
	exeFile := filepath.Join(ExeDir, EXENAME)

	var verboseFlag = ""
	if qpepConfig.Verbose {
		verboseFlag = "--verbose"
	}

	cmd := exec.Command(exeFile,
		fmt.Sprintf("%v", verboseFlag),
		"--acks", fmt.Sprintf("%d", qpepConfig.Acks),
		"--ackDelay", fmt.Sprintf("%d", qpepConfig.AckDelay),
		"--congestion", fmt.Sprintf("%d", qpepConfig.Congestion),
		"--decimate", fmt.Sprintf("%d", qpepConfig.Decimate),
		"--minBeforeDecimation", fmt.Sprintf("%d", qpepConfig.DelayDecimate),
		"--gateway", fmt.Sprintf("%s", qpepConfig.GatewayHost),
		"--port", fmt.Sprintf("%d", qpepConfig.GatewayPort),
		"--apiport", fmt.Sprintf("%d", qpepConfig.GatewayAPIPort),
		"--listenaddress", fmt.Sprintf("%s", qpepConfig.ListenHost),
		"--listenport", fmt.Sprintf("%d", qpepConfig.ListenPort),
		"--multistream", fmt.Sprintf("%v", qpepConfig.MultiStream),
		"--varAckDelay", fmt.Sprintf("%d", qpepConfig.VarAckDelay),
		"--threads", fmt.Sprintf("%d", qpepConfig.WinDivertThreads))

	if cmd == nil {
		ErrorMsg("Could not create client command")
		return nil
	}
	cmd.Dir, _ = filepath.Abs(ExeDir)

	log.Println(cmd.Path)
	log.Println(cmd.Dir)
	return cmd
}

func stopClientProcess() error {
	return stopProcess(clientCmd.Process.Pid)
}
func stopServerProcess() error {
	return stopProcess(serverCmd.Process.Pid)
}

func stopProcess(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		ErrorMsg("Could not terminate client process: %v", err)
		return shared.ErrFailed
	}

	log.Println("Waiting for client exe to terminate")
	if err := proc.Signal(syscall.SIGINT); err != nil {
		ErrorMsg("Could not terminate client process: %v", err)
		return shared.ErrFailed
	}

	return nil
}
