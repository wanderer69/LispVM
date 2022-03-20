package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	. "arkhangelskiy-dv.ru/LispVM/Common"
	. "arkhangelskiy-dv.ru/LispVM/Compiler"
	. "arkhangelskiy-dv.ru/LispVM/BinPack"

	//	"unicode/utf8"
	"github.com/satori/go.uuid"
)

func CompileList(s string, ce *CompilerEnv, debug bool) (*Func, int) {
	l, _, _, err1 := Load_list(0, len(s), len(s), s, 0, "", 0)
	if err1 != 0 {
		fmt.Printf("err %v\r\n", err1)
		return nil, err1
	}

	if debug {
		ss := l.String(false)
		fmt.Printf("l %v\r\n", ss)
	}

	// ce := InitCompilerEnv()
	_, res := Compile(l, ce, debug)
	if !res {
		// найдена ошибка
		return nil, -1
	}
	if len(ce.FuncList) > 0 {
		return ce.FuncList[0], 0
	}
	return nil, -2
}

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

func main() {
	// var file_name = flag.String("file_name", "", "file name source lisp file")
	var verbose = flag.String("v", "", "verbose mode")
	var debug = flag.String("debug", "", "enable debugging")
	var format = flag.String("format", "", "format output file")
	var external = flag.String("external", "", "use external func description")
	var path = flag.String("path", ".", "path to file with external func description")

	// var arguments = flag.String("args", "", "arguments list")
	flag.Parse()
	tail := flag.Args()
	Gformat := "comp"
	if *format == "bin" {
                Gformat = "bin"
        } else {
		if *format == "comp" {
                        Gformat = "comp"
                } else {
                
                }
        }
        ext_func := false
	if *external == "true" {
                ext_func = true
        }
        ext_path := *path

	if false {
		GetID := func() (uuid.UUID, int64) {
		        /*
			u, err := uuid.NewV4()
			if err != nil {
				fmt.Printf("Error %v\r\n", err)
				return u, 0
			}
			*/
			u := uuid.NewV4()
			// fmt.Printf("%s\n", u)
			u1 := binary.BigEndian.Uint64(u[0:8])
			u2 := binary.BigEndian.Uint64(u[8:16])
			fc := int64(u1 + u2)
			return u, fc
		}
		_, fc1 := GetID()
		ifd1 := IntFuncDescr{fc1, "test1", []string{"%v", "%v"}}
		_, fc2 := GetID()
		ifd2 := IntFuncDescr{fc2, "test2", []string{"%v"}}
		id, _ := GetID()
		iff := IntFuncFile{[]byte(id[0:16]), "test_file", []IntFuncDescr{ifd1, ifd2}}
		SaveIntFuncFile("IntFunc_test.json", iff)
		return
	}

	InitFuncOp([]string{ext_path}, ext_func)
	G_debug := false
	// e := InitEnvironment()
	if *debug != "" {
		G_debug = true
	}
	gce := InitGlobalCompilerEnv()
	for i, _ := range tail {
		file_name := tail[i]
		data, err := ioutil.ReadFile(file_name)
		if err != nil {
			fmt.Print(err)
			return
		}
		sl := []string{}
		s := string(data)
		ll := strings.Split(s, "\r\n")
		ll_n := []string{}
		for j, _ := range ll {
			ls := strings.Trim(ll[j], " \r\n\t")
			if len(ls) > 0 {
				if ls[0] == ';' && ls[1] == ';' {
					// это комментарий
				} else {
					ll_n = append(ll_n, ls)
				}
			}
		}
		s = strings.Join(ll_n, "\r\n")
		for {
			//fmt.Printf("'%v'", s)
			l_1, pos_beg, pos_end, err1 := Load_list(0, len(s), len(s), s, 0, "", 0)
			if err1 != 0 {
				fmt.Printf("err %v\r\n", err1)
				break
			}
			//fmt.Printf("l %v %v %v\r\n", l_1, pos_beg, pos_end)
			/*
				if pos_beg == pos_end {
					break
				}
			*/
			if l_1 != nil {
				ss := l_1.String(false)
				//fmt.Printf("l %v\r\n", ss)
				sl = append(sl, ss)
				// fmt.Printf("pos_beg %v, pos_end %v '%v'\r\n", pos_beg, pos_end, s[pos_beg:])
				if pos_beg >= pos_end {
					break
				}
				s = s[pos_beg:]
				s = strings.Trim(s, " \r\n\t")
			} else {
				break
			}
		}
		for _, s := range sl {
			ce := InitCompilerEnv()
			f1, err1 := CompileList(s, ce, G_debug)
			if err1 != 0 {
				// найдена ошибка
				fmt.Printf("Error compile %v\r\n", s)
				return
			}
			gce.FuncToEnv(f1)
		}
		if *verbose != "" {
			for _, ff := range gce.Func_dict {
				fns := ff.String()
				fmt.Printf("%v\r\n", fns)
			}
		}
		_, file := filepath.Split(file_name)
		//fmt.Printf("file %v\r\n", file)
		extension := filepath.Ext(file)
		//fmt.Printf("extension %v\r\n", extension)
		name := file_name[:len(file_name)-len(extension)] //strings.TrimRight(file, extension)
		//fmt.Printf("name %v\r\n", name)
		ms := CreateModule_store(name)
		for _, ff := range gce.Func_dict {
			fns, res := Func2Func_store(ff)
			if res != 0 {
				fmt.Printf("Error compilation %v\r\n", res)
				return
			}
			//fmt.Printf("fns %v\r\n", fns)
			ms.AddFunc(fns)
		}
		if Gformat == "comp" {
		        SaveModule(name+".comp", ms)
		} else {
		        if Gformat == "bin" {
                                bb, err := Save_module_store(ms, 0)
                                if err != nil {
                                        fmt.Printf("Error Save_module_store %v\r\n", err)
                                }
                                
                                //permissions := 0644 // or whatever you need
                                err1 := ioutil.WriteFile(name + ".bin", bb, 0644) // permissions
                                if err1 != nil { 
                		    // handle error
					fmt.Printf("Error compilation write %v\r\n", err1)
		                	return
                                }
                                if false {
                                        SaveModule(name+".comp", ms)
                                        ms__, bb_n, err := Load_module_store(bb, 0)
                                        if err != nil {
                                                fmt.Printf("Error Load_module_store %v\r\n", err)
                                        }
                                        
                                        //fmt.Printf("%#v\r\n", *ms_)
                                        if len(bb_n) > 0 {
						fmt.Printf("Error in file - remainder %v\r\n", bb_n)
						return
                                        }
                                        SaveModule(name+"_.comp", *ms__)
                                }
                        }
                }
	}
}
