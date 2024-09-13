package helper

import (
	"fmt"
	"log"
	"strings"

	"github.com/JamesHovious/w32"
	"golang.org/x/sys/windows"
)

func AllocateProcessMemory(process_handle windows.Handle, allocated_memory_size int, protection_type int) (uintptr, error) {
	msg := fmt.Sprintln("Allocating", fmt.Sprintf("%x", allocated_memory_size), "of memory in process linked to handle", fmt.Sprintf("%x", process_handle), "(", fmt.Sprintf("With protection:%x", protection_type), ")")
	msg = strings.TrimSpace(msg)
	log.Println(msg)
	allocated_data_memory, err := w32.VirtualAllocEx((w32.HANDLE)(process_handle), 0, allocated_memory_size, windows.MEM_COMMIT|windows.MEM_RESERVE, protection_type)
	if err != nil {
		log.Println(msg, ": failed")
		return allocated_data_memory, err
	} else {

		log.Println(msg, ": succeeded at address", fmt.Sprintf("%x", allocated_data_memory))
		return allocated_data_memory, nil
	}

}
