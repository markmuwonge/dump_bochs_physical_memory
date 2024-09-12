// func get_module_export_function_address(process_handle windows.Handle, process_module_handles []windows.Handle, module_base_name string, function_name string) (uint64, error) {

// 	module_handle, err := get_target_process_module_handle(process_handle, process_module_handles, module_base_name)
// 	error_msg := fmt.Sprintln("Unable to get module export function address for", module_base_name, ":", function_name)

// 	file_path, err := get_module_image_file_path(process_handle, module_handle)
// 	if err != nil {
// 		log.Println(error_msg, "(0)")
// 		return 0, err
// 	}

// 	module_pe, err := peparser.New(file_path, nil)
// 	if err != nil {
// 		log.Println(error_msg, "(1)")
// 		return 0, err
// 	}
// 	err = module_pe.Parse()
// 	if err != nil {
// 		log.Println(error_msg, "(2)")
// 		return 0, err
// 	}

// 	exit_thread_export_function := funk.Find(module_pe.Export.Functions, func(export_function peparser.ExportFunction) bool {
// 		return strings.Compare(export_function.Name, function_name) == 0
// 	})

// 	if exit_thread_export_function == nil {
// 		log.Println(error_msg, "(3)")
// 		return 0, err
// 	}

// 	export_function_address := uint64(exit_thread_export_function.(peparser.ExportFunction).FunctionRVA) + uint64(module_handle)
// 	log.Println("Export function address for", module_base_name, ":", function_name, ":", fmt.Sprintf("%x", export_function_address))

// 	return export_function_address, nil
// }