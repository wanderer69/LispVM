rem ..\VM_Main\VM_Main.exe -file_name=iterate.comp -debug=true -args "example_1"
rem ..\VM_Main\VM_Main.exe -file_name=iterate.comp -args "12345"
rem ..\VM_Main\VM_Main.exe -file_name=iterate.comp -args "(1 2 3 4 5)"
rem ..\VM_Main\VM_Main.exe -file_name=iterate.comp -args "[1, 2, 3, 4, 5]"
rem ..\VM_Main\VM_Main.exe -file_name=iterate.comp -args "{\"one\":1, \"two\":2, \"three\":3, \"fourth\":4, \"five\":5}"
..\VM_Main\VM_Main.exe -format=bin -file_name=iterate.bin -args "(1 2 3 4 5)"

