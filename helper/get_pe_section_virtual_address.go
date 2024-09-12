package helper

import (
	"debug/pe"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/thoas/go-funk"
)

func GetPESectionVirtualAddress(file_path string, section_name string) (uint32, error) {
	msg := fmt.Sprintln("Getting PE section", section_name, "virutal address from", file_path)
	msg = strings.TrimSpace(msg)
	pe_file, err := pe.Open(file_path)
	if err != nil {
		log.Println(msg, ":", "failed")
		return 0, err
	}

	section := funk.Find(pe_file.Sections, func(section *pe.Section) bool {
		if strings.Compare(section.Name, section_name) == 0 {
			return true
		} else {
			return false
		}
	})

	if section == nil {
		return 0, errors.New(fmt.Sprintln(msg, ":", "failed"))
	}
	va := section.(*pe.Section).VirtualAddress
	log.Println(msg, "succeeded:", fmt.Sprintf("%x", va))

	return va, nil
}
