package common

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func CreateIntFunc(dir_lst []string, ext_func bool) ([]InternalFuncItem, map[string]InternalFuncItem, map[int64]InternalFuncItem) {
	ifid := make(map[string]InternalFuncItem)
	ifia := []InternalFuncItem{}
	ifiad := make(map[int64]InternalFuncItem)
	ifi := InternalFuncItem{"car", FC_car, []string{"%c"}}
	ifia = append(ifia, ifi)
	ifid["car"] = ifi
	ifiad[FC_car] = ifi

	ifi = InternalFuncItem{"cdr", FC_cdr, []string{"%c"}}
	ifia = append(ifia, ifi)
	ifid["cdr"] = ifi
	ifiad[FC_cdr] = ifi
	ifi = InternalFuncItem{"cons", FC_cons, []string{"%c", "%c"}}
	ifia = append(ifia, ifi)
	ifid["cons"] = ifi
	ifiad[FC_cons] = ifi
	ifi = InternalFuncItem{"listp", FC_listp, []string{"%v"}}
	ifia = append(ifia, ifi)
	ifid["listp"] = ifi
	ifiad[FC_listp] = ifi
	ifi = InternalFuncItem{"eq", FC_eq, []string{"%v", "%v"}}
	ifia = append(ifia, ifi)
	ifid["eq"] = ifi
	ifiad[FC_eq] = ifi
	ifi = InternalFuncItem{"rplaca", FC_rplaca, []string{"%c", "%v"}}
	ifia = append(ifia, ifi)
	ifid["rplaca"] = ifi
	ifiad[FC_rplaca] = ifi
	ifi = InternalFuncItem{"rplacd", FC_rplacd, []string{"%c", "%v"}}
	ifia = append(ifia, ifi)
	ifid["rplacd"] = ifi
	ifiad[FC_rplacd] = ifi
	// ???
	ifi = InternalFuncItem{"list", FC_list, []string{"%c"}}
	ifia = append(ifia, ifi)
	ifid["list"] = ifi
	ifiad[FC_list] = ifi

	ifi = InternalFuncItem{"print", FC_print, []string{"%c"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["print"] = ifi
	ifiad[FC_print] = ifi
	ifi = InternalFuncItem{"+", FC_add, []string{"%v", "%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["+"] = ifi
	ifiad[FC_add] = ifi
	ifi = InternalFuncItem{"-", FC_sub, []string{"%v", "%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["-"] = ifi
	ifiad[FC_sub] = ifi

	ifi = InternalFuncItem{"*", FC_mul, []string{"%v", "%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["*"] = ifi
	ifiad[FC_mul] = ifi

	ifi = InternalFuncItem{"/", FC_div, []string{"%v", "%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["/"] = ifi
	ifiad[FC_div] = ifi

	ifi = InternalFuncItem{"type", FC_type, []string{"%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["type"] = ifi
	ifiad[FC_type] = ifi

	ifi = InternalFuncItem{"str", FC_2str, []string{"%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["str"] = ifi
	ifiad[FC_2str] = ifi

	ifi = InternalFuncItem{"int", FC_2int, []string{"%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["int"] = ifi
	ifiad[FC_2int] = ifi

	ifi = InternalFuncItem{"float", FC_2float, []string{"%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["float"] = ifi
	ifiad[FC_2float] = ifi

	ifi = InternalFuncItem{"typep", FC_typep, []string{"%v", "%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["typep"] = ifi
	ifiad[FC_typep] = ifi

	ifi = InternalFuncItem{"assert", FC_assert, []string{"%v", "%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["assert"] = ifi
	ifiad[FC_assert] = ifi

	ifi = InternalFuncItem{"make_array", FC_make_array, []string{"%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["make_array"] = ifi
	ifiad[FC_make_array] = ifi

	ifi = InternalFuncItem{"make_dict", FC_make_dict, []string{"%v", "%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["make_dict"] = ifi
	ifiad[FC_make_dict] = ifi

	ifi = InternalFuncItem{"array", FC_array, []string{"%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["array"] = ifi
	ifiad[FC_array] = ifi

	ifi = InternalFuncItem{"dict", FC_dict, []string{"%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["dict"] = ifi
	ifiad[FC_dict] = ifi

	ifi = InternalFuncItem{"item", FC_item, []string{"%v", "%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["item"] = ifi
	ifiad[FC_item] = ifi

	ifi = InternalFuncItem{"slice", FC_slice, []string{"%v", "%v", "%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["slice"] = ifi
	ifiad[FC_slice] = ifi

	ifi = InternalFuncItem{"insert", FC_insert, []string{"%v", "%v", "%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["insert"] = ifi
	ifiad[FC_insert] = ifi

	ifi = InternalFuncItem{"append", FC_append, []string{"%v", "%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["append"] = ifi
	ifiad[FC_append] = ifi

	ifi = InternalFuncItem{"import", FC_import, []string{"%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["import"] = ifi
	ifiad[FC_import] = ifi

	ifi = InternalFuncItem{"strip", FC_strip, []string{"%v", "%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["strip"] = ifi
	ifiad[FC_strip] = ifi

	ifi = InternalFuncItem{"split", FC_split, []string{"%v", "%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["split"] = ifi
	ifiad[FC_split] = ifi

	ifi = InternalFuncItem{"index", FC_index, []string{"%v", "%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["index"] = ifi
	ifiad[FC_index] = ifi

	ifi = InternalFuncItem{"length", FC_length, []string{"%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["length"] = ifi
	ifiad[FC_length] = ifi

	ifi = InternalFuncItem{"join", FC_join, []string{"%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["join"] = ifi
	ifiad[FC_join] = ifi

	ifi = InternalFuncItem{"set_dict", FC_set_dict, []string{"%v", "%v", "%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["set_dict"] = ifi
	ifiad[FC_set_dict] = ifi

	ifi = InternalFuncItem{"load_lib", FC_load_lib, []string{"%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["load_lib"] = ifi
	ifiad[FC_load_lib] = ifi

	ifi = InternalFuncItem{"call_lib_func", FC_call_lib_func, []string{"%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["call_lib_func"] = ifi
	ifiad[FC_call_lib_func] = ifi

	ifi = InternalFuncItem{"unload_lib", FC_unload_lib, []string{"%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["unload_lib"] = ifi
	ifiad[FC_unload_lib] = ifi

	ifi = InternalFuncItem{"princ", FC_princ, []string{"%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["princ"] = ifi
	ifiad[FC_princ] = ifi

	ifi = InternalFuncItem{"format", FC_format, []string{"%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["format"] = ifi
	ifiad[FC_format] = ifi

	ifi = InternalFuncItem{"iterate", FC_iterate, []string{"%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["iterate"] = ifi
	ifiad[FC_iterate] = ifi

	ifi = InternalFuncItem{"iterate_item", FC_get_iter, []string{"%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["iterate_item"] = ifi
	ifiad[FC_get_iter] = ifi

	ifi = InternalFuncItem{"make_iter", FC_make_iter, []string{"%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["make_iter"] = ifi
	ifiad[FC_make_iter] = ifi

	ifi = InternalFuncItem{"not", FC_not, []string{"%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["not"] = ifi
	ifiad[FC_not] = ifi

	ifi = InternalFuncItem{"find_func", FC_find_func, []string{"%v"}} // , "%r"
	ifia = append(ifia, ifi)
	ifid["find_func"] = ifi
	ifiad[FC_find_func] = ifi

	/*
		ifi = InternalFuncItem{"", FC_, []string{"%v"}} // , "%r"
		ifia = append(ifia, ifi)
		ifid[""] = ifi
		ifiad[FC_] = ifi

		ifi = InternalFuncItem{"", FC_, []string{"%v"}} // , "%r"
		ifia = append(ifia, ifi)
		ifid[""] = ifi
		ifiad[FC_] = ifi
	*/
	LoadIntFuncDescrFiles := func(dir string) {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			fmt.Printf("%v\r\n", err)
			return
		}

		for _, file := range files {
			//fmt.Println(file.Name())
			ss := file.Name()
			if strings.Index(ss, "IntFunc_") == 0 {
				iff, err := LoadIntFuncFile(file.Name())
				if err == nil {
					for _, ifd := range iff.Funcs {
						ifi := InternalFuncItem{ifd.Name, ifd.FuncCode, ifd.Args}
						ifia = append(ifia, ifi)
						ifid[ifd.Name] = ifi
						ifiad[ifd.FuncCode] = ifi
					}
				}
			}
		}
	}
	if ext_func {
		for _, dir := range dir_lst {
			LoadIntFuncDescrFiles(dir)
		}
	}
	return ifia, ifid, ifiad
}

/*
func AddIntFunc(ifd IntFuncDescr, ifia *[]InternalFuncItem, ifid map[string]InternalFuncItem, ifiad map[int64]InternalFuncItem) {
	ifi := InternalFuncItem{ifd.Name, ifd.FuncCode, ifd.Args}
	*ifia = append(*ifia, ifi)
	ifid[ifd.Name] = ifi
	ifiad[ifd.FuncCode] = ifi
}
*/
