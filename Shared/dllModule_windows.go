package shared_module

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
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
	Proc_prt *syscall.Proc
}

type ExtLib struct {
	FileName    string
	Path        string
	DLL_ptr     *syscall.DLL
	ExtFuncList []ExtFunc
	ExtFuncDict map[string]ExtFunc
}

func LoadExtLib(file_lib_name string) (*ExtLib, error) {
	el := ExtLib{}
	el.FileName = file_lib_name
	h, err1 := syscall.LoadDLL(file_lib_name)
	if err1 != nil {
		fmt.Println(err1)
		return nil, err1
	}
	el.DLL_ptr = h
	el.ExtFuncDict = make(map[string]ExtFunc)
	return &el, nil
}

func (ex *ExtLib) CallExtFunc(func_name string, data_b []CellShort) (int, *CellShort, error) {
	data_ptr := unsafe.Pointer(&data_b)
	res := 0
	len_in := len(data_b)
	len_out := 0 // unsafe.Pointer(&css_in1)
	ef_, ok := ex.ExtFuncDict[func_name]
	proc_func := ef_.Proc_prt
	if !ok {
		proc_func_, err := ex.DLL_ptr.FindProc(func_name)
		if err != nil {
			fmt.Printf("LoadDataCell %v\r\n", err)
			return -1, nil, err
		}
		proc_func = proc_func_
		ef := ExtFunc{}
		ef.Name = func_name
		ef.Proc_prt = proc_func_
		ex.ExtFuncDict[func_name] = ef
	}
	//     fmt.Printf("proc_func %v\r\n", proc_func)
	r1_2, r2_2, lastErr := proc_func.Call(uintptr(data_ptr), uintptr(len_in), uintptr(unsafe.Pointer(&len_out)), uintptr(unsafe.Pointer(&res)))
	//     fmt.Printf("r1_2 %v r2_2 %v lastErr %v\r\n", r1_2, r2_2, lastErr)
	le := fmt.Sprintf("%v", lastErr)
	if le != "The operation completed successfully." {
		fmt.Printf("Proc %v Call %v %v\r\n", func_name, lastErr, r2_2)
		return -1, nil, errors.New(le)
	}
	//fmt.Printf("len_out %v\r\n", len_out)
	var css_out CellShort
	if len_out > 0 {
		css_out = *(*CellShort)(unsafe.Pointer(r1_2))
		return res, &css_out, nil
	}
	//C.free(unsafe.Pointer(r1_2))
	return res, nil, nil
}

func (ex *ExtLib) Resume() {
	ex.DLL_ptr.Release()
}
