package launcher

import (
	"fmt"
	"syscall"
)

func LaunchProcess(gamepath string) error {
	cmdline, _ := syscall.UTF16PtrFromString(gamepath)
	startupInfo := new(syscall.StartupInfo)
	processInfo := new(syscall.ProcessInformation)

	err := syscall.CreateProcess(nil, cmdline, nil, nil, false, 0, nil, nil, startupInfo, processInfo)

	if err != nil {
		return fmt.Errorf("error creating process: %s", err.Error())
	}

	pid := processInfo.ProcessId
	fmt.Println("Started Palia with PID:", pid)

	syscall.CloseHandle(processInfo.Process)
	syscall.CloseHandle(processInfo.Thread)

	return nil
}
