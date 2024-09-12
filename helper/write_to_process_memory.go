package helper

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

func WriteToProcessMemory(process_handle windows.Handle, address uintptr, bytes []byte) error {
	msg := fmt.Sprintln("Writing", len(bytes), "bytes to", fmt.Sprintf("%x", address), "for process linked to handle", fmt.Sprintf("%x", process_handle))
	msg = strings.TrimSpace(msg)
	log.Println(msg)

	bytes_size := len(bytes)
	var bytes_written uintptr
	err := windows.WriteProcessMemory(process_handle, address, (*byte)(unsafe.Pointer(&bytes[0])), uintptr(bytes_size), &bytes_written)

	if err != nil {
		log.Println(msg, ":failed")
		return err
	} else if int(bytes_written) != bytes_size {

		return errors.New(fmt.Sprintln(msg, ":failed"))
	}

	log.Println(msg, ":succeeded")
	return nil
}
