package helper

import (
	"fmt"
	"log"
	"syscall"

	"golang.org/x/sys/windows"
)

func GetProcessHandle(pid int) (windows.Handle, error) {
	log.Println("Getting process handle for pid", pid)
	target_parent_process_handle, err := windows.OpenProcess(syscall.STANDARD_RIGHTS_REQUIRED|syscall.SYNCHRONIZE|0xfff, false, uint32(pid))
	if err != nil {
		log.Println("Unable to get process handle for pid", pid)
		return 0, err
	}
	log.Println("Returning handle for pid", pid, ":", fmt.Sprintf("%x", target_parent_process_handle))
	return target_parent_process_handle, nil
}
