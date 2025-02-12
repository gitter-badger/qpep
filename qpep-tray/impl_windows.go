package main

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/parvit/qpep/shared"
)

const (
	BASEDIR_ENVIRONMENTVAR = "APPDATA"
	EXENAME                = "qpep.exe"
)

func getClientCommand() *exec.Cmd {
	exeFile, _ := filepath.Abs(filepath.Join(ExeDir, EXENAME))
	//handle, _ := syscall.GetCurrentProcess()

	var verboseFlag = ""
	if qpepConfig.Verbose {
		verboseFlag = "-verbose"
	}

	attr := &syscall.SysProcAttr{
		HideWindow: true,
		CmdLine: fmt.Sprintf(
			exeFile+` `+
				`-client `+
				`%s `+
				`-threads %d `+
				`-gateway "%s" `+
				`-port %d `+
				`-apiport %d `+
				`-listenaddress "%s" `+
				`-listenport %d `+
				`-acks %d `+
				`-ackDelay %d `+
				`-congestion %d `+
				`-decimate %d `+
				`-minBeforeDecimation %d `+
				`-multistream %v `+
				`-varAckDelay %d `,
			verboseFlag,
			qpepConfig.WinDivertThreads,
			qpepConfig.GatewayHost,
			qpepConfig.GatewayPort,
			qpepConfig.GatewayAPIPort,
			qpepConfig.ListenHost,
			qpepConfig.ListenPort,
			qpepConfig.Acks,
			qpepConfig.AckDelay,
			qpepConfig.Congestion,
			qpepConfig.Decimate,
			qpepConfig.DelayDecimate,
			qpepConfig.MultiStream,
			qpepConfig.VarAckDelay),
	}

	cmd := exec.Command(exeFile)
	if cmd == nil {
		ErrorMsg("Could not create client command")
		return nil
	}
	cmd.Dir, _ = filepath.Abs(ExeDir)
	cmd.SysProcAttr = attr

	log.Println(cmd.Path)
	log.Println(cmd.Dir)
	log.Println(cmd.SysProcAttr.CmdLine)
	return cmd
}

func getServerCommand() *exec.Cmd {
	exeFile := filepath.Join(ExeDir, EXENAME)
	//handle, _ := syscall.GetCurrentProcess()

	var verboseFlag = ""
	if qpepConfig.Verbose {
		verboseFlag = "-verbose"
	}

	attr := &syscall.SysProcAttr{
		HideWindow: true,
		CmdLine: fmt.Sprintf(
			exeFile+` `+
				`%s `+
				`-threads %d `+
				`-gateway "%s" `+
				`-port %d `+
				`-apiport %d `+
				`-listenaddress "%s" `+
				`-listenport %d `+
				`-acks %d `+
				`-ackDelay %d `+
				`-congestion %d `+
				`-decimate %d `+
				`-minBeforeDecimation %d `+
				`-multistream %v `+
				`-varAckDelay %d `,
			verboseFlag,
			qpepConfig.WinDivertThreads,
			qpepConfig.GatewayHost,
			qpepConfig.GatewayPort,
			qpepConfig.GatewayAPIPort,
			qpepConfig.ListenHost,
			qpepConfig.ListenPort,
			qpepConfig.Acks,
			qpepConfig.AckDelay,
			qpepConfig.Congestion,
			qpepConfig.Decimate,
			qpepConfig.DelayDecimate,
			qpepConfig.MultiStream,
			qpepConfig.VarAckDelay),
	}

	cmd := exec.Command(exeFile)
	if cmd == nil {
		ErrorMsg("Could not create client command")
		return nil
	}
	cmd.Dir, _ = filepath.Abs(ExeDir)
	cmd.SysProcAttr = attr

	log.Println(cmd.Path)
	log.Println(cmd.Dir)
	log.Println(cmd.SysProcAttr.CmdLine)
	return cmd
}

func stopClientProcess() error {
	return stopProcess(clientCmd.Process.Pid)
}
func stopServerProcess() error {
	return stopProcess(serverCmd.Process.Pid)
}

func stopProcess(pid int) error {
	d, e := syscall.LoadDLL("kernel32.dll")
	if e != nil {
		return shared.ErrFailed
	}
	p, e := d.FindProc("GenerateConsoleCtrlEvent")
	if e != nil {
		return shared.ErrFailed
	}
	r, _, e := p.Call(syscall.CTRL_BREAK_EVENT, uintptr(pid))
	if r == 0 {
		return shared.ErrFailed
	}

	return nil
}
