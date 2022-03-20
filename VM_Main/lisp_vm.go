package main

import (
	"flag"
	"fmt"
	"sync"
	"io/ioutil"

	//	"io/ioutil"
	//	"path/filepath"
	//	"strings"

	//	"unicode/utf8"
	. "arkhangelskiy-dv.ru/LispVM/Common"
	. "arkhangelskiy-dv.ru/LispVM/VM"
	. "arkhangelskiy-dv.ru/LispVM/BinPack"
)

func Str2Arg(s string) (*Cell, int) {
	s_len := len(s)

	s_c, _, _, err := Load_list(0, s_len, s_len, s, 0, "", 0)
	if err != 0 {
		fmt.Printf("err %v\r\n", err)
		return nil, err
	}
	/*
		ss := s_c.String(false)
		fmt.Printf("s_c %v\r\n", ss)
	*/
	return s_c, 0
}

/*
func LoadModuleBin(file_name string) (Module_store, error) {
	var cp Module_store
	data, err := ioutil.ReadFile(file_name)
	if err != nil {
		fmt.Print(err)
		return cp, err
	}

	err = json.Unmarshal(data, &cp)
	if err != nil {
		fmt.Println("error:", err)
		return cp, err
	}
	return cp, nil
}
*/

func main() {
	var file_name = flag.String("file_name", "", "file name source lisp file")
	//	var mode = flag.String("mode", "", "compile,execute")
	var debug = flag.String("debug", "", "enable debugging true|false")
	var no_result = flag.String("noresult", "", "no print result true|false")
	var arguments = flag.String("args", "", "arguments list")
	var format = flag.String("format", "", "file format")
	var view = flag.String("view", "", "view print true|false")
	flag.Parse()

	G_no_result := true
	G_debug := false
	if *debug == "true" {
		G_debug = true
	}
	if *no_result == "false" {
		G_no_result = false
	}
	G_fformat := "comp"
	if *format == "comp" {
	} else {
		if *format == "bin" {
                       G_fformat = "bin"
		} else {
		}
	}

	// InitFunc([]string{"."})
	InitIntFunc()
	no_print := true
	if *view == "false" {
		no_print = false
	}
	e := InitEnvironment(no_print)

	var ms Module_store

	// считывание
	if G_fformat == "comp" {
		ms_, err := LoadModule(*file_name)
		if err != nil {
			fmt.Printf("%v\r\n", err)
		}
		ms = ms_
	} else {
		if G_fformat == "bin" {
			data, err := ioutil.ReadFile(*file_name)
			if err != nil {
				fmt.Print(err)
				return
			}

                        ms_, bb_n, err := Load_module_store(data, 0)
                        if err != nil {
                                fmt.Printf("Error Load_module_store %v\r\n", err)
                        }
                        ms = *ms_
                        //fmt.Printf("%#v\r\n", *ms_)
                        if len(bb_n) > 0 {
				fmt.Printf("Error in file - remainder %v\r\n", bb_n)
				return
                        }
                        //fmt.Printf("bb_n %v\r\n", bb_n)
		} else {
		}
	}

	for _, ff := range ms.Func_list {
		// fmt.Printf("%v\r\n", ff)
		f2 := Func_store2Func(ff)
		// fmt.Printf("%v\r\n", f2)
		e.FuncToEnv(&f2)
	}
	if false {
		for _, ff := range e.Func_dict {
			// fmt.Printf("-> %#v\r\n", ff)
			fns := ff.String()
			fmt.Printf("%v\r\n", fns)
		}
	}
	e.FileFormat = G_fformat
	// исполнение
	var wg sync.WaitGroup
	// init
	f_init, ok := e.Func_dict["init"]
	if !ok {
		fmt.Printf("Func init not found\r\n")
	} else {
		args := []*Cell{Nil}
		ProcExec(e, f_init, args, &wg, e.Current_proc, G_no_result, G_debug)
		// not wait all processes
		//wg.Wait()
	}
	// main
	f, ok := e.Func_dict["main"]
	if !ok {
		fmt.Printf("Func main not found. Execution cancelled.\r\n")
		return
	}
	arg1_c, err1 := Str2Arg(*arguments)
	if err1 != 0 {
		fmt.Printf("Error in Str2Arg\r\n")
	}
	args := []*Cell{arg1_c}

	if false {
		fc := FuncCall(e, f, -1, e.Current_proc)
		SetCurrentExec(e, fc, args, e.Current_proc)
		for {
			res := Exec(e, e.Current_proc, G_debug)
			if res == 0 {
				if G_no_result {
					// печатаем результат
					result := *e.Call_Item[e.Current_proc].Result
					for _, v := range result {
						ss := v.String(false)
						fmt.Printf("result %v\r\n", ss)
					}
				}
				break
			}
		}
	}
	// var wg sync.WaitGroup
	ProcExec(e, f, args, &wg, e.Current_proc, G_no_result, G_debug)
	// wait stop all process
	wg.Wait()
}
