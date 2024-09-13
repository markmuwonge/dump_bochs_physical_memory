package helper

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

func GetConfig(string config_file_path) (gjson.Result, error) {
	msg := fmt.Sprintln("Getting config")
	msg = strings.TrimSpace(msg)
	err_msg := fmt.Sprintln(msg, ":failed")
	config_file_bytes, err := os.ReadFile(config_path)
	if err != nil {
		log.Println(err_msg)
		return nil, err
	}
	cofnig_file_string := string(config_file_bytes)
	if !gjson.Valid(cofnig_file_string) {
		log.Println(err_msg)
		return nil, err
	}

	log.Println(msg, ":succeeded")
	return gjson.Parse(json), nil
}
