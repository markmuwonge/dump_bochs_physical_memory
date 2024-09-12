package helper

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/JamesHovious/w32"
	"golang.org/x/sys/windows"
)

func FreeProcessMemory(process_handle windows.Handle, allocated_data_memory_address uintptr) error {
	msg := fmt.Sprintln("Freeing memory at", fmt.Sprintf("%x", allocated_data_memory_address), "linked to process handle", fmt.Sprintf("%x", process_handle))
	msg = strings.TrimSpace(msg)
	log.Println(msg)

	ret := w32.VirtualFreeEx((w32.HANDLE)(process_handle), allocated_data_memory_address, 0, windows.MEM_RELEASE)

	if ret == false {
		return errors.New(fmt.Sprintln(msg, "failed"))
	}
	log.Println(msg, "succeeded")
	return nil
}
