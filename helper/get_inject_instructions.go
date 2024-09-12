package helper

import (
	"fmt"
	"log"
	"strings"

	asm "github.com/twitchyliquid64/golang-asm"

	"github.com/twitchyliquid64/golang-asm/obj"
	"github.com/twitchyliquid64/golang-asm/obj/x86"
)

func GetInjectInstructions(bochs_dbg_fetch_mem_address uint64, bochs_dbg_cpu_pointer_address uint64, bochs_dbg_fetch_mem_phy_addr_param uint64, bochs_dbg_fetch_mem_address_data_size_param uint8, allocated_data_memory_address uint64, allocated_data_memory_size int64, allocated_data2_memory_address uint64) ([]byte, error) {
	msg := fmt.Sprintln("Getting inject instructions")
	msg = strings.TrimSpace(msg)
	instructions := make([]byte, 0)
	asm_builder, err := asm.NewBuilder("amd64", 64)
	log.Println(msg)
	if err != nil {
		log.Println(msg, ":failed")
		return instructions, err
	}

	instruction_pointers := make([]*obj.Prog, 0)
	registers := []int16{x86.REG_AX, x86.REG_CX, x86.REG_DX, x86.REG_R8, x86.REG_R9, x86.REG_AX, x86.REG_AX}
	addresses := []int64{int64(allocated_data2_memory_address), int64(bochs_dbg_cpu_pointer_address), int64(bochs_dbg_fetch_mem_phy_addr_param), int64(bochs_dbg_fetch_mem_address_data_size_param), int64(allocated_data_memory_address), int64(bochs_dbg_fetch_mem_address), 0}
	loop_counter_address := addresses[0] + 8 //8 byte
	for i := 0; i < len(registers); i++ {
		if i == 0 {
			//call to bochs_dbg_fetch_mem removes the return address from stack that's required when createremote thread exits
			pop_return_address_instruction := asm_builder.NewProg()
			pop_return_address_instruction.As = x86.APOPQ
			pop_return_address_instruction.To.Type = obj.TYPE_MEM
			pop_return_address_instruction.To.Offset = addresses[i]
			instruction_pointers = append(instruction_pointers, pop_return_address_instruction)

			init_loop_counter_instruction := asm_builder.NewProg()
			init_loop_counter_instruction.As = x86.AMOVQ
			init_loop_counter_instruction.To.Type = obj.TYPE_MEM
			init_loop_counter_instruction.To.Offset = loop_counter_address
			init_loop_counter_instruction.From.Type = obj.TYPE_CONST
			init_loop_counter_instruction.From.Offset = 0
			instruction_pointers = append(instruction_pointers, init_loop_counter_instruction)

			compare_counter_instruction := asm_builder.NewProg()
			compare_counter_instruction.As = x86.ACMPQ
			compare_counter_instruction.From.Type = obj.TYPE_MEM
			compare_counter_instruction.From.Offset = loop_counter_address
			compare_counter_instruction.To.Type = obj.TYPE_CONST
			compare_counter_instruction.To.Offset = allocated_data_memory_size
			instruction_pointers = append(instruction_pointers, compare_counter_instruction)

			//je last instruciton - 1
			jmp_if_counter_done_instruction := asm_builder.NewProg()
			jmp_if_counter_done_instruction.As = x86.AJEQ
			jmp_if_counter_done_instruction.To.Type = obj.TYPE_BRANCH
			instruction_pointers = append(instruction_pointers, jmp_if_counter_done_instruction)

		}

		mov_imm_to_reg_instruction := asm_builder.NewProg()
		mov_imm_to_reg_instruction.As = x86.AMOVQ
		mov_imm_to_reg_instruction.To.Type = obj.TYPE_REG
		mov_imm_to_reg_instruction.To.Reg = registers[i]
		mov_imm_to_reg_instruction.From.Type = obj.TYPE_CONST
		mov_imm_to_reg_instruction.From.Offset = addresses[i]
		instruction_pointers = append(instruction_pointers, mov_imm_to_reg_instruction)

		if registers[i] == x86.REG_DX || registers[i] == x86.REG_R9 {
			inc_phy_mem_param_or_dest_buf_instrcution := asm_builder.NewProg()
			inc_phy_mem_param_or_dest_buf_instrcution.As = x86.AADDQ
			inc_phy_mem_param_or_dest_buf_instrcution.To.Type = obj.TYPE_REG
			inc_phy_mem_param_or_dest_buf_instrcution.To.Reg = registers[i]
			inc_phy_mem_param_or_dest_buf_instrcution.From.Type = obj.TYPE_MEM
			inc_phy_mem_param_or_dest_buf_instrcution.From.Offset = loop_counter_address
			instruction_pointers = append(instruction_pointers, inc_phy_mem_param_or_dest_buf_instrcution)
		}

		if i == len(registers)-2 {
			call_rax_instruction := asm_builder.NewProg()
			call_rax_instruction.As = obj.ACALL
			call_rax_instruction.To.Type = obj.TYPE_REG
			call_rax_instruction.To.Reg = x86.REG_AX
			instruction_pointers = append(instruction_pointers, call_rax_instruction)

			increment_instruction := asm_builder.NewProg()
			increment_instruction.As = x86.AINCW
			increment_instruction.To.Type = obj.TYPE_MEM
			increment_instruction.To.Offset = addresses[0] + 8
			instruction_pointers = append(instruction_pointers, increment_instruction)

			jmp_instruction := asm_builder.NewProg()
			jmp_instruction.As = obj.AJMP
			jmp_instruction.To.Type = obj.TYPE_BRANCH
			jmp_instruction.To.Val = instruction_pointers[2]
			instruction_pointers = append(instruction_pointers, jmp_instruction)

			push_return_address_instruction := asm_builder.NewProg()
			push_return_address_instruction.As = x86.APUSHQ
			push_return_address_instruction.From.Type = obj.TYPE_MEM
			push_return_address_instruction.From.Offset = addresses[0]
			instruction_pointers = append(instruction_pointers, push_return_address_instruction)
		}

		if i == len(registers)-1 {
			return_instruction := asm_builder.NewProg()
			return_instruction.As = obj.ARET
			instruction_pointers = append(instruction_pointers, return_instruction)

			//setting destination of jmp_if_counter_done_instruction (destination instruction not present in instruction_pointers array until later)
			instruction_pointers[3].To.SetTarget(instruction_pointers[len(instruction_pointers)-3])
		}
	}

	for _, ip := range instruction_pointers {
		asm_builder.AddInstruction(ip)
	}

	for _, b := range asm_builder.Assemble() {
		instructions = append(instructions, b)
	}

	log.Println(msg, ":succeeded", fmt.Sprintf("%x", instructions))

	return instructions, nil
}
