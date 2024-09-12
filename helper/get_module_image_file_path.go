package helper

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/sys/windows"
)

func GetModuleImageFilePath(process_handle windows.Handle, module_handle windows.Handle) (string, error) {
	msg := fmt.Sprintln("Getting module image file path of module", fmt.Sprintf("%x", module_handle), "linked to process handle", fmt.Sprintf("%x", process_handle))
	msg = strings.TrimSpace(msg)
	log.Println(msg)
	file_name := make([]uint16, windows.MAX_PATH)
	err := windows.GetModuleFileNameEx(process_handle, module_handle, &file_name[0], windows.MAX_PATH)
	if err != nil {
		log.Println(msg, "failed")
		return "", err
	}

	file_path := windows.UTF16ToString(file_name)
	log.Println(msg, "succeeded:", file_path)
	return file_path, nil
}
