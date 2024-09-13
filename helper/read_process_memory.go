package helper

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"golang.org/x/sys/windows"
)

func ReadProcessMemory(process_handle windows.Handle, address uintptr, number_of_bytes_to_read uint64) ([]byte, error) {
	msg := fmt.Sprintln("Reading", number_of_bytes_to_read, "bytes at", fmt.Sprintf("%x", address), "for process linked to handle", fmt.Sprintf("%x", process_handle))
	msg = strings.TrimSpace(msg)
	log.Println(msg)

	buffer := make([]byte, number_of_bytes_to_read)
	var number_of_bytes_read uintptr
	err_msg := fmt.Sprintln(msg, ":failed")
	err := windows.ReadProcessMemory(process_handle, address, &buffer[0], uintptr(number_of_bytes_to_read), &number_of_bytes_read)
	if err != nil {
		log.Println(err_msg)
		return nil, err
	} else if uint64(number_of_bytes_read) != number_of_bytes_to_read {
		log.Println(err_msg)
		return nil, errors.New(err_msg)
	}
	log.Println(msg, ":succeeded")
	return buffer, nil

}
