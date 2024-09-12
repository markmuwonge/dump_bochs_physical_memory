package helper

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"golang.org/x/sys/windows"
)

func GetProcessModuleHandle(process_handle windows.Handle, module_name string) (windows.Handle, error) {
	msg := fmt.Sprintln("*Getting module handle for module", module_name, "using process handle", fmt.Sprintf("%x", process_handle))
	msg = strings.TrimSpace(msg)
	module_handles, err := get_process_module_handles(process_handle)
	if err != nil {
		log.Println(msg, ":", "failed")
		return 0, err
	}

	var msvcrt_module_handle windows.Handle
	for _, h := range module_handles {
		base_name := make([]uint16, windows.MAX_PATH)
		windows.GetModuleBaseName(process_handle, h, &base_name[0], windows.MAX_PATH)

		if strings.Compare(module_name, windows.UTF16ToString(base_name)) == 0 {
			msvcrt_module_handle = h
			break
		}
	}
	if msvcrt_module_handle == 0 {
		return 0, errors.New(fmt.Sprintln(msg, ":", "failed"))
	}

	log.Println(msg, "succeeded:", fmt.Sprintf("%x", msvcrt_module_handle))
	return msvcrt_module_handle, nil
}
