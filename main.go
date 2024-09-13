package main

import (
	"dump_bochs_physical_memory/helper"
	"fmt"
	"log"

	"golang.org/x/sys/windows"
)

var target_process_name = "bochsdbg.exe"

// Bochs-2.6.11
var bochs_image_md5 = "46bb6791c069ca4323ea891b5a370b78"

// // var bochs_dbg_fetch_mem_call_text_section_offset = 0x1231EA
var bochs_dbg_fetch_mem_text_section_offset = 0x1B2620
var bochs_dbg_cpu_pointer_data_section_offset = 0x75FC0

var bochs_dbg_fetch_mem_phy_addr_param = 0xc0000
var bochs_dbg_fetch_mem_address_data_size_param uint8 = 1
var allocated_memory_size = 0x1000

func main() {
	log.Println("*")

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
	return

	target_process_image_text_section_va, err := helper.GetPESectionVirtualAddress(target_process_image_file_path, ".text")
	fatal(err)

	bochs_dbg_fetch_mem_address := uint64(target_process_image_module_handle) + uint64(target_process_image_text_section_va) + uint64(bochs_dbg_fetch_mem_text_section_offset)
	log.Println("*Bochs dbg_fetch_men address:", fmt.Sprintf("%x", bochs_dbg_fetch_mem_address))

	target_process_image_data_section_va, err := helper.GetPESectionVirtualAddress(target_process_image_file_path, ".data")
	fatal(err)

	bochs_dbg_cpu_pointer_address := uint64(target_process_image_module_handle) + uint64(target_process_image_data_section_va) + uint64(bochs_dbg_cpu_pointer_data_section_offset)
	log.Println("*Bochs dbg_cpu pointer address:", fmt.Sprintf("%x", bochs_dbg_cpu_pointer_address))

	allocated_data_memory_address, err := helper.AllocateProcessMemory(target_process_handle, allocated_memory_size, windows.PAGE_READWRITE)
	fatal(err)

	allocated_data2_memory_address, err := helper.AllocateProcessMemory(target_process_handle, allocated_memory_size, windows.PAGE_READWRITE)
	fatal(err)

	inject_instructions, err := helper.GetInjectInstructions(bochs_dbg_fetch_mem_address, bochs_dbg_cpu_pointer_address, uint64(bochs_dbg_fetch_mem_phy_addr_param), bochs_dbg_fetch_mem_address_data_size_param, uint64(allocated_data_memory_address), int64(allocated_memory_size), uint64(allocated_data2_memory_address))
	if err != nil {
		helper.FreeProcessMemory(target_process_handle, allocated_data_memory_address)
		helper.FreeProcessMemory(target_process_handle, allocated_data2_memory_address)
	}
	fatal(err)

	allocated_code_memory_address, err := helper.AllocateProcessMemory(target_process_handle, allocated_memory_size, windows.PAGE_EXECUTE_READ)
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

	_ = thread_handle
	_ = thread_id

	helper.FreeProcessMemory(target_process_handle, allocated_data_memory_address)
	helper.FreeProcessMemory(target_process_handle, allocated_data2_memory_address)
	helper.FreeProcessMemory(target_process_handle, allocated_code_memory_address)

}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
