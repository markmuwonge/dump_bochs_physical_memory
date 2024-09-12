package helper

import (
	"fmt"
	"log"
	"unsafe"

	"golang.org/x/sys/windows"
)

func get_process_module_handles(target_parent_process_handle windows.Handle) ([]windows.Handle, error) {
	log.Println("Getting process module handles for process with handle:", fmt.Sprintf("%x", target_parent_process_handle))
	var module_handle windows.Handle
	var module_handle_bytes_needed uint32
	err := windows.EnumProcessModules(target_parent_process_handle, &module_handle, uint32(unsafe.Sizeof(module_handle)), &module_handle_bytes_needed)
	err_msg := fmt.Sprintln("Unable to get process module handles for process with handle:", fmt.Sprintf("%x", target_parent_process_handle))
	if err != nil {
		log.Println(err_msg)
		return []windows.Handle{}, err
	}

	module_handles := make([]windows.Handle, module_handle_bytes_needed/uint32(unsafe.Sizeof(module_handle)))
	err = windows.EnumProcessModules(target_parent_process_handle, &module_handles[0], uint32(unsafe.Sizeof(module_handle))*module_handle_bytes_needed, &module_handle_bytes_needed)

	if err != nil {
		log.Println(err_msg)
		return []windows.Handle{}, err
	}

	log.Println(len(module_handles), "module handles found for process with handle:", fmt.Sprintf("%x", target_parent_process_handle))
	return module_handles, nil
}
