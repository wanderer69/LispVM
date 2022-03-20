package binpack

import (
	"fmt"
	"unsafe"
	"bytes"
	"encoding/binary"

	//	"io/ioutil"
	//	"path/filepath"
	//	"strings"

	//	"unicode/utf8"
	. "arkhangelskiy-dv.ru/LispVM/Common"
)

type Type_lenght_header struct {
	LenFieldName  int32
//	TypeValue     uint16
	LenValue      int32
}

const ID_string = 10
const ID_int32  = 11

func Save_type_lenght_value(field_name string, bb []byte) []byte {
        type_lenght_header := Type_lenght_header{LenFieldName: int32(len(field_name))/*, TypeValue: t_type*/, LenValue: int32(len(bb))}

	len_all := (int)(unsafe.Sizeof(type_lenght_header)) + len(field_name) + len(bb)
	b_in := make([]byte, 0, len_all)
	var buf = bytes.NewBuffer(b_in)
	if err := binary.Write(buf, binary.LittleEndian, &type_lenght_header); err != nil {
		fmt.Printf("Save_type_lenght_value header %v\r\n", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, []byte(field_name)); err != nil {
		fmt.Printf("Save_type_lenght_value name %v\r\n", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, bb); err != nil {
		fmt.Printf("Save_type_lenght_value bb %v\r\n", err)
	}
	return buf.Bytes()
}

type Func_code_store_header struct {
	Code int32 `json:"Code,omitempty"` // код некомпилированной команды
	Const_len int32
}

const ID_func_code_store = 1

type Lenght_header struct {
//	LenFieldName  int32
//	TypeValue     uint16
	LenValue      int32
}

type Cell_store_header struct {
        Len           int32
	Type          int32
	Value_int     int64
	Value_float   float64
/*	
	Value_sym_len     int32
	Value_str_len     int32
	Value_head_len    int32
	Value_last_len    int32
*/

	Value_dict_len    int32
	Value_array_len   int32
/*
	Value_func_len    int32
	Value_obj_len     int32
	Value_ext_len     int32
	Value_channel_len int32
*/
}

//const ID_string = 10
//const ID_int32  = 11

func Save_lenght_value(bb []byte) []byte {
        lenght_header := Lenght_header{LenValue: int32(len(bb))}

	len_all := (int)(unsafe.Sizeof(lenght_header)) + len(bb)
	b_in := make([]byte, 0, len_all)
	var buf = bytes.NewBuffer(b_in)
	if err := binary.Write(buf, binary.LittleEndian, &lenght_header); err != nil {
		fmt.Println(err)
	}
	if err := binary.Write(buf, binary.LittleEndian, bb); err != nil {
		fmt.Println(err)
	}
	return buf.Bytes()
}


const ID_Cell_store = 2

func Save_func_code_store(func_code_store Func_code_store, debug int) ([]byte, error) {
        if debug > 2 {
                fmt.Printf("Save_func_code_store\r\n")
        }
        func_code_store_header := Func_code_store_header{Code: func_code_store.Code, Const_len: int32(len(func_code_store.Const))}

        var bb []byte
        // var bbv []byte

	len_all := (int)(unsafe.Sizeof(func_code_store_header)) // + len(func_code_store.Const)*8 //len(int64)
        for i:=0; i < len(func_code_store.Const); i++ {
                var v int64
                v = int64(func_code_store.Const[i])
                
		b_in := make([]byte, 0, unsafe.Sizeof(v))
		var buf = bytes.NewBuffer(b_in)
		if err := binary.Write(buf, binary.LittleEndian, &v); err != nil {
			fmt.Println(err)
			return []byte{}, err
		}
	        bbl := buf.Bytes()
                bbh := Save_lenght_value([]byte(bbl))
                bb = append(bb, bbh...)
        }
        len_all = len_all + len(bb)

	b_in := make([]byte, 0, len_all)
	var buf = bytes.NewBuffer(b_in)
	if err := binary.Write(buf, binary.LittleEndian, &func_code_store_header); err != nil {
		fmt.Println(err)
		return []byte{}, err
	}

	if err := binary.Write(buf, binary.LittleEndian, &bb); err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
//	fmt.Printf("-> %v\r\n", buf.Bytes())
	return buf.Bytes(), nil
}

func Save_cell_store(cell_store *Cell_store, debug int) ([]byte, error) {
        // для длинных полей отдельная TLV с указанием поля
        if debug > 2 {
        fmt.Printf("Save_cell_store\r\n")
        }
        if debug > 4 {
        fmt.Printf("%#v\r\n", *cell_store)
        }
        cell_store_header := Cell_store_header{Type: cell_store.Type, 
        Value_int: cell_store.Value_int, Value_float: cell_store.Value_float, 
        }
        cell_store_header.Value_dict_len = int32(len(cell_store.Value_dict))
        cell_store_header.Value_array_len = int32(len(cell_store.Value_array))

        var bb []byte
        var bbh []byte
        var bbl []byte
        len_all := (int)(unsafe.Sizeof(cell_store_header))
        switch cell_store.Type {
        case Cell_sym:
        	bb = Save_lenght_value([]byte(cell_store.Value_sym))
        	len_all = len_all + len(bb)
/*
        case Cell_int:
        case Cell_float:
*/
        case Cell_str:
        	bb = Save_lenght_value([]byte(cell_store.Value_str))
        	len_all = len_all + len(bb)
        case Cell_cell:
                csh := cell_store.Value_head[0]
                h, err := Save_cell_store(&csh, debug)
                if err != nil {
			return []byte{}, err
                }
                csl := cell_store.Value_last[0]
                l, err := Save_cell_store(&csl, debug)
                if err != nil {
			return []byte{}, err
                }
        	bbh = Save_lenght_value([]byte(h))
        	bbl = Save_lenght_value([]byte(l)) 
        	len_all = len_all + len(bbh)
        	len_all = len_all + len(bbl)
                bb = append(bb, bbh...)
                bb = append(bb, bbl...)
        	//fmt.Printf("len_all %v\r\n", len_all)
        case Cell_dict:
                for k, v := range cell_store.Value_dict {
                        h, err := Save_cell_store(&v, debug)
                        if err != nil {
				return []byte{}, err
                        }
                        bb_ := Save_type_lenght_value(k, []byte(h))
                        bb = append(bb, bb_...)
                }
        	len_all = len_all + len(bb)
        case Cell_array:
                for i, _ := range cell_store.Value_array {
                        csh := cell_store.Value_array[i]
                        h, err := Save_cell_store(&csh, debug)
                        if err != nil {
				return []byte{}, err
                        }
        		bb_ := Save_lenght_value([]byte(h))
                        bb = append(bb, bb_...)
        	}
        	len_all = len_all + len(bb)
        case Cell_func:
                fmt.Printf("func\r\n")
        case Cell_object:
                fmt.Printf("object\r\n")
        case Cell_channel:
                fmt.Printf("channel\r\n")
        case Cell_dict_n:
                fmt.Printf("dict\r\n")
        case Cell_multiple:
                fmt.Printf("multiple\r\n")

        }
	b_in := make([]byte, 0, len_all)
	var buf = bytes.NewBuffer(b_in)
	cell_store_header.Len = int32(len_all)
	//fmt.Printf("cell_store_header %#v\r\n", cell_store_header)
	if err := binary.Write(buf, binary.LittleEndian, &cell_store_header); err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
	if err := binary.Write(buf, binary.LittleEndian, &bb); err != nil {
		fmt.Println(err)
		return []byte{}, err
	}

	//fmt.Printf("cell_store %v\r\n", buf.Bytes())
	return buf.Bytes(), nil
}

type Func_store_header struct {
        Len		int32
	Type		int32             `json:"type"` // тип 0 встроенная 1 внешняя
	Code_len	int32
	Var_list_len	int32
	Const_list_len	int32
}

func Save_func_store(func_store Func_store, debug int) ([]byte, error) {
        if debug > 2 {
        fmt.Printf("Save_func_store\r\n")
        }
        func_store_header := Func_store_header{Type: func_store.Type}

        var bb []byte

        len_all := (int)(unsafe.Sizeof(func_store_header))

        bb_name := Save_lenght_value([]byte(func_store.Name))
        len_all = len_all + len(bb_name)
        bb = append(bb, bb_name...)

        cs_args := func_store.Args
        //fmt.Printf("cs_args %#v\r\n", cs_args)
        args, err := Save_cell_store(&cs_args, debug)
        if err != nil {
		return []byte{}, err
        }
        //fmt.Printf("args  %#v\r\n", len(args))
        bb_args := Save_lenght_value([]byte(args))
        len_all = len_all + len(bb_args)
        bb = append(bb, bb_args...)

        for i, _ := range func_store.Code {
                fcs, err := Save_func_code_store(func_store.Code[i], debug)
                if err != nil {
			return []byte{}, err
                }

                bb_fcs := Save_lenght_value([]byte(fcs))
                len_all = len_all + len(bb_fcs)
                bb = append(bb, bb_fcs...)
        }

        for i, _ := range func_store.Var_list {
                bb_vl := Save_lenght_value([]byte(func_store.Var_list[i]))
                len_all = len_all + len(bb_vl)
                bb = append(bb, bb_vl...)
        }

        for i, _ := range func_store.Const_list {
                csh := func_store.Const_list[i]
                h, err := Save_cell_store(&csh, debug)
                if err != nil {
			return []byte{}, err
                }
        	bb_ := Save_lenght_value([]byte(h))
                len_all = len_all + len(bb_)
                bb = append(bb, bb_...)
        }

	func_store_header.Code_len        = int32(len(func_store.Code))
	func_store_header.Var_list_len    = int32(len(func_store.Var_list))
	func_store_header.Const_list_len  = int32(len(func_store.Const_list))

	b_in := make([]byte, 0, len_all)
	func_store_header.Len = int32(len_all)
	var buf = bytes.NewBuffer(b_in)
	if err := binary.Write(buf, binary.LittleEndian, &func_store_header); err != nil {
		fmt.Println(err)
	}
	if err := binary.Write(buf, binary.LittleEndian, &bb); err != nil {
		fmt.Println(err)
	}

	return buf.Bytes(), nil
}

type Module_store_header struct {
        Len		int32
        Type		int32
        Var_list_len	int32
        Const_list_len	int32
        Func_list_len	int32
}

func Save_module_store(module_store Module_store, debug int) ([]byte, error) {
        if debug > 0 {
        fmt.Printf("Save_module_store\r\n")
        }
        module_store_header := Module_store_header{Type: 0}

        var bb []byte

        len_all := (int)(unsafe.Sizeof(module_store_header))

        bb_name := Save_lenght_value([]byte(module_store.Name))
        len_all = len_all + len(bb_name)
        bb = append(bb, bb_name...)

        module_store_header.Var_list_len = int32(len(module_store.Var_list))
        module_store_header.Const_list_len = int32(len(module_store.Const_list))
        module_store_header.Func_list_len = int32(len(module_store.Func_list))

        for i, _ := range module_store.Var_list {
                csh := module_store.Var_list[i]
                h, err := Save_cell_store(&csh, debug)
                if err != nil {
			return []byte{}, err
                }
                bb_vl := Save_lenght_value([]byte(h))
                len_all = len_all + len(bb_vl)
                bb = append(bb, bb_vl...)
        }

        for i, _ := range module_store.Const_list {
                csh := module_store.Const_list[i]
                h, err := Save_cell_store(&csh, debug)
                if err != nil {
			return []byte{}, err
                }
        	bb_ := Save_lenght_value([]byte(h))
                len_all = len_all + len(bb_)
                bb = append(bb, bb_...)
        }

        for i, _ := range module_store.Func_list{
                bb_fl, err := Save_func_store(module_store.Func_list[i], debug)
                if err != nil {
			return []byte{}, err
                }
        	bb_ := Save_lenght_value([]byte(bb_fl))
                len_all = len_all + len(bb_)
                bb = append(bb, bb_...)
        }

	b_in := make([]byte, 0, len_all)
	var buf = bytes.NewBuffer(b_in)
	if err := binary.Write(buf, binary.LittleEndian, &module_store_header); err != nil {
		fmt.Println(err)
	}
	if err := binary.Write(buf, binary.LittleEndian, &bb); err != nil {
		fmt.Println(err)
	}

	return buf.Bytes(), nil
}

func Load_module_store(bb []byte, debug int) (*Module_store, []byte, error) {
        pos := 0
        if debug > 0 {
        fmt.Printf("Load_module_store\r\n")
        }
        var module_store Module_store
        var module_store_header Module_store_header
        // len_module_store_header := (int)(unsafe.Sizeof(module_store_header))

	var buf = bytes.NewBuffer(make([]byte, 0, len(bb)))
	if err := binary.Write(buf, binary.BigEndian, &bb); err != nil {
		fmt.Println(err)
		return nil, []byte{}, err
	}

	if err := binary.Read(buf, binary.LittleEndian, &module_store_header); err != nil {
		fmt.Println(err)
		return nil, []byte{}, err
	}

        var lenght_header_ Lenght_header

	if err := binary.Read(buf, binary.LittleEndian, &lenght_header_); err != nil {
		fmt.Println(err)
		return nil, []byte{}, err
	}
	bbs := make([]byte, lenght_header_.LenValue)
	if err := binary.Read(buf, binary.LittleEndian, &bbs); err != nil {
		fmt.Println(err)
		return nil, []byte{}, err
	}
        module_store.Name = string(bbs)
//        fmt.Printf("40\r\n")

/*
         = module_store.Var_list
        module_store_header.Const_list_len = module_store.Const_list
        module_store_header.Func_list_len = module_store.Func_list
*/
        pos = pos + int(unsafe.Sizeof(module_store_header))
        if module_store_header.Var_list_len > 0 {
                for i := 0; i < int(module_store_header.Var_list_len); i++ {
                        var lenght_header Lenght_header
                
			if err := binary.Read(buf, binary.LittleEndian, &lenght_header); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
			bbs := make([]byte, lenght_header.LenValue)
			if err := binary.Read(buf, binary.LittleEndian, &bbs); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
	        
                        // pos = pos + int(unsafe.Sizeof(lenght_header))
                        cs, _, err:= Load_cell_store(bbs, debug) // bb[pos:pos+int(lenght_header.LenValue)]
                        if err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
                        module_store.Var_list = append(module_store.Var_list, *cs)
                }
        }
//        fmt.Printf("60\r\n")

        if module_store_header.Const_list_len > 0 {
                for i := 0; i < int(module_store_header.Const_list_len); i++ {
                        var lenght_header Lenght_header
                
			if err := binary.Read(buf, binary.LittleEndian, &lenght_header); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
			bbs := make([]byte, lenght_header.LenValue)
			if err := binary.Read(buf, binary.LittleEndian, &bbs); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
	        
                        // pos = pos + int(unsafe.Sizeof(lenght_header))
                        cs, _, err:= Load_cell_store(bbs, debug) // bb[pos:pos+int(lenght_header.LenValue)]
                        if err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
                        module_store.Const_list = append(module_store.Const_list, *cs)
                }
        }

//        fmt.Printf("80 module_store_header.Func_list_len %v\r\n", module_store_header.Func_list_len)

        if module_store_header.Func_list_len > 0 {
                for i := 0; i < int(module_store_header.Func_list_len); i++ {
                        var lenght_header Lenght_header
                
			if err := binary.Read(buf, binary.LittleEndian, &lenght_header); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
//        fmt.Printf("81 lenght_header.LenValue %v\r\n", lenght_header.LenValue)

			bbs := make([]byte, lenght_header.LenValue)
			if err := binary.Read(buf, binary.LittleEndian, &bbs); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
//        fmt.Printf("82\r\n")
	        
                        // pos = pos + int(unsafe.Sizeof(lenght_header))
                        cs, _, err:= Load_func_store(bbs, debug) // bb[pos:pos+int(lenght_header.LenValue)]
                        if err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
                        module_store.Func_list = append(module_store.Func_list, *cs)
                }
        }

//        fmt.Printf("100\r\n")

        bb_r := buf.Next(buf.Len())
        return &module_store, bb_r, nil
}

func Load_cell_store(bb []byte, debug int) (*Cell_store, []byte, error) {
        //pos := 0
        if debug > 2 {
        fmt.Printf("Load_cell_store\r\n")
        }
        var cell_store Cell_store
        var cell_store_header Cell_store_header
        // len_module_store_header := (int)(unsafe.Sizeof(module_store_header))

	var buf = bytes.NewBuffer(make([]byte, 0, len(bb)))
	if err := binary.Write(buf, binary.BigEndian, &bb); err != nil {
		fmt.Println(err)
		return nil, []byte{}, err
	}

	if err := binary.Read(buf, binary.LittleEndian, &cell_store_header); err != nil {
		fmt.Println(err)
		return nil, []byte{}, err
	}
        cell_store.Type = cell_store_header.Type 
        cell_store.Value_int = cell_store_header.Value_int 
        cell_store.Value_float = cell_store_header.Value_float

        // fmt.Printf("cell_store_header %#v\r\n", cell_store_header)
        switch cell_store.Type {
        case Cell_sym:
                var lenght_header Lenght_header

		if err := binary.Read(buf, binary.LittleEndian, &lenght_header); err != nil {
			fmt.Println(err)
			return nil, []byte{}, err
		}
		ss := make([]byte, lenght_header.LenValue)
		if err := binary.Read(buf, binary.LittleEndian, &ss); err != nil {
			fmt.Println(err)
			return nil, []byte{}, err
		}
                cell_store.Value_sym = string(ss)
        case Cell_str:
                var lenght_header Lenght_header

		if err := binary.Read(buf, binary.LittleEndian, &lenght_header); err != nil {
			fmt.Println(err)
			return nil, []byte{}, err
		}
		ss := make([]byte, lenght_header.LenValue)
		if err := binary.Read(buf, binary.LittleEndian, &ss); err != nil {
			fmt.Println(err)
			return nil, []byte{}, err
		}
                cell_store.Value_str = string(ss)
        case Cell_cell:
                var lenght_header Lenght_header

		if err := binary.Read(buf, binary.LittleEndian, &lenght_header); err != nil {
			fmt.Println(err)
			return nil, []byte{}, err
		}
		ssh := make([]byte, lenght_header.LenValue)
		if err := binary.Read(buf, binary.LittleEndian, &ssh); err != nil {
			fmt.Println(err)
			return nil, []byte{}, err
		}

                csh, _, err:= Load_cell_store(ssh, debug) // bb[pos:pos+int(lenght_header.LenValue)]
                if err != nil {
			fmt.Println(err)
			return nil, []byte{}, err
		}
                cell_store.Value_head = append(cell_store.Value_head, *csh)

		if err := binary.Read(buf, binary.LittleEndian, &lenght_header); err != nil {
			fmt.Println(err)
			return nil, []byte{}, err
		}
		ssl := make([]byte, lenght_header.LenValue)
		if err := binary.Read(buf, binary.LittleEndian, &ssl); err != nil {
			fmt.Println(err)
			return nil, []byte{}, err
		}

                csl, _, err:= Load_cell_store(ssl, debug) // bb[pos:pos+int(lenght_header.LenValue)]
                if err != nil {
			fmt.Println(err)
			return nil, []byte{}, err
		}
                cell_store.Value_last = append(cell_store.Value_last, *csl)
        case Cell_dict:
                cell_store.Value_dict = make(Cell_store_arr)
                for i:=0;i < int(cell_store_header.Value_dict_len); i++ {
                        var type_lenght_header Type_lenght_header
                        
			if err := binary.Read(buf, binary.LittleEndian, &type_lenght_header); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
			ssv := make([]byte, type_lenght_header.LenFieldName)
			if err := binary.Read(buf, binary.LittleEndian, &ssv); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
			ssh := make([]byte, type_lenght_header.LenValue)
			if err := binary.Read(buf, binary.LittleEndian, &ssh); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
		        fname := string(ssv)
                        csh, _, err:= Load_cell_store(ssh, debug) // bb[pos:pos+int(lenght_header.LenValue)]
                        if err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
                        cell_store.Value_dict[fname] = *csh
        	}
        case Cell_array:
                for i:=0; i < int(cell_store_header.Value_array_len); i++ {
                        var lenght_header Lenght_header
                        
			if err := binary.Read(buf, binary.LittleEndian, &lenght_header); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
			ssh := make([]byte, lenght_header.LenValue)
			if err := binary.Read(buf, binary.LittleEndian, &ssh); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
		        
                        csh, _, err:= Load_cell_store(ssh, debug) // bb[pos:pos+int(lenght_header.LenValue)]
                        if err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
                        cell_store.Value_array = append(cell_store.Value_array, *csh)
        	}
/*
        case Cell_func:
        case Cell_object:
        case Cell_channel:
        case Cell_dict_n:
        case Cell_multiple:
*/
        }
        bb_r := buf.Next(buf.Len())
	return &cell_store, bb_r, nil
}

func Load_func_store(bb []byte, debug int) (*Func_store, []byte, error) {
//        pos := 0
        if debug > 2 {
        fmt.Printf("Load_func_store\r\n")
        }
        var func_store Func_store
        var func_store_header Func_store_header
        // len_module_store_header := (int)(unsafe.Sizeof(module_store_header))

	var buf = bytes.NewBuffer(make([]byte, 0, len(bb)))
	if err := binary.Write(buf, binary.BigEndian, &bb); err != nil {
		fmt.Printf("%v\r\n", err)
		return nil, []byte{}, err
	}

	if err := binary.Read(buf, binary.LittleEndian, &func_store_header); err != nil {
		fmt.Printf("%v\r\n", err)
		return nil, []byte{}, err
	}

        func_store.Type =  func_store_header.Type

        var lenght_header_ Lenght_header

	if err := binary.Read(buf, binary.LittleEndian, &lenght_header_); err != nil {
		fmt.Printf("%v\r\n", err)
		return nil, []byte{}, err
	}
	bbs := make([]byte, lenght_header_.LenValue)
	if err := binary.Read(buf, binary.LittleEndian, &bbs); err != nil {
		fmt.Printf("%v\r\n", err)
		return nil, []byte{}, err
	}
        func_store.Name = string(bbs)

//        fmt.Printf("func_store %#v\r\n", func_store)

        var lenght_header_a Lenght_header
        
	if err := binary.Read(buf, binary.LittleEndian, &lenght_header_a); err != nil {
		fmt.Printf("%v\r\n", err)
		return nil, []byte{}, err
	}
	bba := make([]byte, lenght_header_a.LenValue)
	if err := binary.Read(buf, binary.LittleEndian, &bba); err != nil {
		fmt.Printf("%v\r\n", err)
		return nil, []byte{}, err
	}
//	fmt.Printf("lenght_header_a %v\r\n", lenght_header_a)
//	fmt.Printf("bba %v\r\n", bba)
        // pos = pos + int(unsafe.Sizeof(lenght_header))
        cs, _, err:= Load_cell_store(bba, debug) // bb[pos:pos+int(lenght_header.LenValue)]
        if err != nil {
		fmt.Printf("%v\r\n", err)
		return nil, []byte{}, err
	}
        func_store.Args = *cs
//        fmt.Printf("func_store %#v\r\n", func_store)
//        fmt.Printf("func_store_header %#v\r\n", func_store_header)
        if func_store_header.Code_len > 0 {
                for i := 0; i < int(func_store_header.Code_len); i++ {
                        var lenght_header Lenght_header
                
			if err := binary.Read(buf, binary.LittleEndian, &lenght_header); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
//                        fmt.Printf("lenght_header %#v\r\n", lenght_header)

			bbs := make([]byte, lenght_header.LenValue)
			if err := binary.Read(buf, binary.LittleEndian, &bbs); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}

//                        fmt.Printf("bbs %v\r\n", bbs)
	        
                        // pos = pos + int(unsafe.Sizeof(lenght_header))
                        cs, _, err:= Load_func_code_store(bbs, debug) // bb[pos:pos+int(lenght_header.LenValue)]
                        if err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
                        func_store.Code = append(func_store.Code, *cs)
                }
        }
        if func_store_header.Var_list_len > 0 {
                for i := 0; i < int(func_store_header.Var_list_len); i++ {
                        var lenght_header Lenght_header
                
			if err := binary.Read(buf, binary.LittleEndian, &lenght_header); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
			bbs := make([]byte, lenght_header.LenValue)
			if err := binary.Read(buf, binary.LittleEndian, &bbs); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
                        func_store.Var_list = append(func_store.Var_list, string(bbs))
                }
        }
        if func_store_header.Const_list_len > 0 {
                for i := 0; i < int(func_store_header.Const_list_len); i++ {
                        var lenght_header Lenght_header
                
			if err := binary.Read(buf, binary.LittleEndian, &lenght_header); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
			bbs := make([]byte, lenght_header.LenValue)
			if err := binary.Read(buf, binary.LittleEndian, &bbs); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
	        
                        // pos = pos + int(unsafe.Sizeof(lenght_header))
                        cs, _, err:= Load_cell_store(bbs, debug) // bb[pos:pos+int(lenght_header.LenValue)]
                        if err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
                        func_store.Const_list = append(func_store.Const_list, *cs)
                }
        }

        bb_r := buf.Next(buf.Len())
        return &func_store, bb_r, nil
}

func Load_func_code_store(bb []byte, debug int) (*Func_code_store, []byte, error) {
//        pos := 0
        if debug > 2 {
        fmt.Printf("Load_func_code_store %v\r\n", bb)
        }
        var func_code_store Func_code_store
        var func_code_store_header Func_code_store_header
        // len_module_store_header := (int)(unsafe.Sizeof(module_store_header))

	var buf = bytes.NewBuffer(make([]byte, 0, len(bb)))
	if err := binary.Write(buf, binary.BigEndian, &bb); err != nil {
		fmt.Println(err)
		return nil, []byte{}, err
	}

	if err := binary.Read(buf, binary.LittleEndian, &func_code_store_header); err != nil {
		fmt.Println(err)
		return nil, []byte{}, err
	}
        func_code_store.Code = func_code_store_header.Code
//	fmt.Printf("func_code_store_header %#v\r\n", func_code_store_header)
        if func_code_store_header.Const_len > 0 {
                var lenght_header Lenght_header
                
		if err := binary.Read(buf, binary.LittleEndian, &lenght_header); err != nil {
			fmt.Println(err)
			return nil, []byte{}, err
		}

                for i := 0; i < int(func_code_store_header.Const_len); i++ {
			// bbs := make([]byte, lenght_header.LenValue)
			var val int64
			if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
				fmt.Println(err)
				return nil, []byte{}, err
			}
                        func_code_store.Const = append(func_code_store.Const, val)
                }
        }
        bb_r := buf.Next(buf.Len())
        return &func_code_store, bb_r, nil
}

