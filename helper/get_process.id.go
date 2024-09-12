package helper

import (
	"errors"
	"fmt"
	"log"

	go_ps "github.com/mitchellh/go-ps"
)

func GetProcessId(base_name string) (int, error) {
	log.Println("*Geting process id for", base_name)

	err_msg := fmt.Sprintln("Unable to get process id for", base_name)
	processes, err := go_ps.Processes()
	if err != nil {
		log.Println(err_msg)
		return 0, err
	}
	if len(processes) == 0 {
		log.Println(err_msg)
		return 0, errors.New(err_msg)
	}
	var target_process go_ps.Process
	var target_process_id int
	for _, p := range processes {
		if p.Executable() == base_name {
			target_process_id = p.Pid()
			log.Println("Process id found for", p.Executable(), ":", target_process_id)
			target_process = p
			break
		}
	}
	if target_process == nil {
		return 0, errors.New(err_msg)
	}
	return target_process.Pid(), nil
}
