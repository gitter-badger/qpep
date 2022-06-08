package main

import (
	"log"
	"os/exec"
)

var serverCmd *exec.Cmd

func startServer() error {
	if serverCmd != nil {
		log.Println("ERROR: Cannot start an already running server, first stop it")
		return ErrFailed
	}

	serverCmd = getServerCommand()

	if err := serverCmd.Start(); err != nil {
		ErrorMsg("Could not start server program: %v", err)
		serverCmd = nil
		return ErrCommandNotStarted
	}
	InfoMsg("Server started")

	return nil
}

func stopServer() error {
	if serverCmd == nil {
		log.Println("ERROR: Cannot stop an already server client, first start it")
		return nil
	}

	if err := stopServerProcess(); err != nil {
		if ok := ConfirmMsg("Could not stop process gracefully (%v), do you want to force-terminate it?", err); !ok {
			return err
		}
		if err := serverCmd.Process.Kill(); err != nil {
			ErrorMsg("Could not force-terminate process")
			return err
		}
	}

	serverCmd.Wait()
	serverCmd = nil
	InfoMsg("Server stopped")
	return nil
}
