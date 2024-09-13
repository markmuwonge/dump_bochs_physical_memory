package helper

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/thoas/go-funk"
)

func MD5Match(file_path string, md5_test string) error {
	msg := fmt.Sprintln("Testing if file at", file_path, "matches md5", md5_test)
	msg = strings.TrimSpace(msg)

	log.Println(msg)

	target_process_image_file_bytes, err := os.ReadFile(file_path)
	err_msg := fmt.Sprintln(msg, ":false")
	if err != nil {
		log.Println(err_msg)
		return err
	}

	target_process_image_file_md5 := funk.Map(md5.Sum(target_process_image_file_bytes), func(b byte) byte {
		return b
	}).([]byte)

	if fmt.Sprintf("%x", target_process_image_file_md5) == md5_test {
		log.Println(msg, ":true")
		return nil
	} else {
		log.Println(err_msg)
		return errors.New(err_msg)
	}

}
