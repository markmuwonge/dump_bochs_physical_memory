package main

import (
	"dump_bochs_physical_memory/helper"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/alexflint/go-arg"
	"github.com/itchio/ox/syscallex"
	"golang.org/x/sys/windows"
)

var target_process_name = "bochsdbg.exe"

// Bochs-2.6.11
var bochs_image_md5 = "46bb6791c069ca4323ea891b5a370b78"

// // var bochs_dbg_fetch_mem_call_text_section_offset = 0x1231EA
var bochs_dbg_fetch_mem_text_section_offset = 0x1B2620
var bochs_dbg_cpu_pointer_data_section_offset = 0x75FC0

var bochs_dbg_fetch_mem_address_data_size_param uint8 = 1

var allocated_data2_memory_size = 0x1000
var allocated_code_memory_size = 0x1000

// int((math.Pow(2, 20)) * 100)
var thread_timeout_millis = 10000

var args struct {
	PHYADDR      uint   `arg:"required"`
	BYTECOUNT    uint   `arg:"required"`
	DUMPFILEPATH string `arg:"required"`
}

func main() {
	log.Println("*")
	arg.MustParse(&args)

	target_process_id, err := helper.GetProcessId(target_process_name)
	fatal(err)

	target_process_handle, err := helper.GetProcessHandle(target_process_id)
	fatal(err)

	target_process_image_module_handle, err := helper.GetProcessModuleHandle(target_process_handle, target_process_name)
	fatal(err)

	target_process_image_file_path, err := helper.GetModuleImageFilePath(target_process_handle, target_process_image_module_handle)
	fatal(err)

	err = helper.MD5Match(target_process_image_file_path, bochs_image_md5)
	fatal(err)

	target_process_image_text_section_va, err := helper.GetPESectionVirtualAddress(target_process_image_file_path, ".text")
	fatal(err)

	bochs_dbg_fetch_mem_address := uint64(target_process_image_module_handle) + uint64(target_process_image_text_section_va) + uint64(bochs_dbg_fetch_mem_text_section_offset)
	log.Println("*Bochs dbg_fetch_men address:", fmt.Sprintf("%x", bochs_dbg_fetch_mem_address))

	target_process_image_data_section_va, err := helper.GetPESectionVirtualAddress(target_process_image_file_path, ".data")
	fatal(err)

	bochs_dbg_cpu_pointer_address := uint64(target_process_image_module_handle) + uint64(target_process_image_data_section_va) + uint64(bochs_dbg_cpu_pointer_data_section_offset)
	log.Println("*Bochs dbg_cpu pointer address:", fmt.Sprintf("%x", bochs_dbg_cpu_pointer_address))

	allocated_data_memory_address, err := helper.AllocateProcessMemory(target_process_handle, int(args.BYTECOUNT), windows.PAGE_READWRITE)
	fatal(err)

	allocated_data2_memory_address, err := helper.AllocateProcessMemory(target_process_handle, allocated_data2_memory_size, windows.PAGE_READWRITE)
	fatal(err)

	inject_instructions, err := helper.GetInjectInstructions(bochs_dbg_fetch_mem_address, bochs_dbg_cpu_pointer_address, uint64(args.PHYADDR), bochs_dbg_fetch_mem_address_data_size_param, uint64(allocated_data_memory_address), int64(args.BYTECOUNT), uint64(allocated_data2_memory_address))
	if err != nil {
		helper.FreeProcessMemory(target_process_handle, allocated_data_memory_address)
		helper.FreeProcessMemory(target_process_handle, allocated_data2_memory_address)
	}
	fatal(err)

	allocated_code_memory_address, err := helper.AllocateProcessMemory(target_process_handle, allocated_code_memory_size, windows.PAGE_EXECUTE_READ)
	if err != nil {
		helper.FreeProcessMemory(target_process_handle, allocated_data_memory_address)
		helper.FreeProcessMemory(target_process_handle, allocated_data2_memory_address)
	}
	fatal(err)

	err = helper.WriteToProcessMemory(target_process_handle, allocated_code_memory_address, inject_instructions)
	if err != nil {
		helper.FreeProcessMemory(target_process_handle, allocated_data_memory_address)
		helper.FreeProcessMemory(target_process_handle, allocated_data2_memory_address)
		helper.FreeProcessMemory(target_process_handle, allocated_code_memory_address)
	}
	fatal(err)

	thread_handle, thread_id, err := helper.CreateRemoteThread(target_process_handle, allocated_code_memory_address)
	if err != nil {
		helper.FreeProcessMemory(target_process_handle, allocated_data_memory_address)
		helper.FreeProcessMemory(target_process_handle, allocated_data2_memory_address)
		helper.FreeProcessMemory(target_process_handle, allocated_code_memory_address)
	}
	fatal(err)

	_, err = windows.WaitForSingleObject((windows.Handle)(thread_handle), uint32(thread_timeout_millis))
	if err != nil {
		helper.FreeProcessMemory(target_process_handle, allocated_data_memory_address)
		helper.FreeProcessMemory(target_process_handle, allocated_data2_memory_address)
		helper.FreeProcessMemory(target_process_handle, allocated_code_memory_address)
	}
	fatal(err)

	thread_exit_code, err := syscallex.GetExitCodeThread((syscall.Handle)(thread_handle))
	if err != nil {
		helper.FreeProcessMemory(target_process_handle, allocated_data_memory_address)
		helper.FreeProcessMemory(target_process_handle, allocated_data2_memory_address)
		helper.FreeProcessMemory(target_process_handle, allocated_code_memory_address)
	}
	fatal(err)
	log.Println("*Remote thread", thread_id, "exit code:", thread_exit_code)

	read_process_memory, err := helper.ReadProcessMemory(target_process_handle, allocated_data_memory_address, uint64(args.BYTECOUNT))
	if err != nil {
		helper.FreeProcessMemory(target_process_handle, allocated_data_memory_address)
		helper.FreeProcessMemory(target_process_handle, allocated_data2_memory_address)
		helper.FreeProcessMemory(target_process_handle, allocated_code_memory_address)
	}

	helper.FreeProcessMemory(target_process_handle, allocated_data_memory_address)
	helper.FreeProcessMemory(target_process_handle, allocated_data2_memory_address)
	helper.FreeProcessMemory(target_process_handle, allocated_code_memory_address)

	dump_file, err := os.Create(args.DUMPFILEPATH)
	fatal(err)

	_, err = dump_file.Write(read_process_memory)
	fatal(err)

	log.Println("**")
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
