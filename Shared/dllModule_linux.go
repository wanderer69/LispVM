package shared_module

import (
	"errors"
	"fmt"
//	"syscall"
	"plugin"

//	"unsafe"
)

/*
#include <stdlib.h>
*/
import "C"

import (
	. "arkhangelskiy-dv.ru/LispVM/Shared/common"
)
/*
type CellShort struct {
	Type        int32
	Value_int   int64
	Value_float float64
	Value_sym   string
	Value_str   string
	Value_head  *CellShort
	Value_last  *CellShort
	Value_dict  map[string]*CellShort
	Value_array []*CellShort
}
*/
type ExtFunc struct {
	Name     string
	Proc_prt TExtFunc
}

type ExtLib struct {
	FileName    string
	Path        string
	Plug        *plugin.Plugin
	ExtFuncList []ExtFunc
	ExtFuncDict map[string]ExtFunc
}

func LoadExtLib(file_lib_name string) (*ExtLib, error) {
	el := ExtLib{}
	el.FileName = file_lib_name
	// load module
	// 1. open the so file to load the symbols
	plug, err := plugin.Open(el.FileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	el.Plug = plug
	el.ExtFuncDict = make(map[string]ExtFunc)
	return &el, nil
}

func (ex *ExtLib) CallExtFunc(func_name string, data_b []CellShort) (int, *CellShort, error) {
	// data_ptr := unsafe.Pointer(&data_b)
	res := 0
	ef_, ok := ex.ExtFuncDict[func_name]
	proc_func := ef_.Proc_prt
	if !ok {
		// 2. look up a symbol (an exported function or variable)
		// in this case, variable Greeter
		sym, err := ex.Plug.Lookup(func_name)
		if err != nil {
			fmt.Printf("LoadDataCell lookup %v\r\n", err)
			return -1, nil, err
		}
		// 3. Assert that loaded symbol is of a desired type
		// in this case interface type Greeter (defined above)
		//var proc_func_ TExtFunc
		proc_func_, ok := sym.(func (pa []CellShort, res *int) *CellShort)
		if !ok {
			fmt.Printf("LoadDataCell unexpected type from module symbol %v\r\n", func_name)
			return -1, nil, err
		}
                proc_func = proc_func_
		ef := ExtFunc{}
		ef.Name = func_name
		ef.Proc_prt = proc_func_
		ex.ExtFuncDict[func_name] = ef
	}

	// 4. use the module
	cs_ret := proc_func(data_b, &res)
	if res != 0 {
		le := fmt.Sprintf("Proc %v Call %v %v", func_name, res)
		fmt.Printf("%v\r\n", le)
		return -1, nil, errors.New(le)
	}
	return res, cs_ret, nil
}

func (ex *ExtLib) Resume() {
//	ex.DLL_ptr.Release()
}
