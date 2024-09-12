package helper

import (
	"fmt"
	"log"
	"strings"

	"github.com/JamesHovious/w32"
	"golang.org/x/sys/windows"
)

func CreateRemoteThread(process_handle windows.Handle, start_address uintptr) (w32.HANDLE, uint32, error) {
	msg := fmt.Sprintln("Creating remote thread using process handle", fmt.Sprintf("%x", process_handle), "with start address", fmt.Sprintf("%x", start_address))
	msg = strings.TrimSpace(msg)
	log.Println(msg)
	thread_handle, thread_id, err := w32.CreateRemoteThread((w32.HANDLE)(process_handle), nil, 0, uint32(start_address), 0, 0)
	if err != nil {
		log.Println(msg, ":failed")
		return 0, 0, err
	}
	log.Println(msg, ":succeeded with thread handle", fmt.Sprintf("%x", thread_handle), "and thread id", thread_id)
	return thread_handle, thread_id, err

}
