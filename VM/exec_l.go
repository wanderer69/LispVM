package vm

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"io/ioutil"

	. "arkhangelskiy-dv.ru/LispVM/Common"
	. "arkhangelskiy-dv.ru/LispVM/Shared"
	. "arkhangelskiy-dv.ru/LispVM/Shared/common"
	. "arkhangelskiy-dv.ru/LispVM/BinPack"
)

var G_FuncsInt_dict map[int64]InternalFuncItem

func InitFunc(dir_lst []string, load_ext_func bool) {
	_, _, d2 := CreateIntFunc(dir_lst, load_ext_func)
	G_FuncsInt_dict = d2
}

// вызов функции
func FuncCall(env *EnvironmentLVM, f *Func, Func_call_prev_ptr int, proc_id int) int {
	// строим структуру исполнения
	fc := Func_call{}
	fc.Env = env
	fc.Global_func_dict = make(map[string]*Func)
	fc.Var_dict = make(map[string]*Cell)
	fc.Const_dict = make(map[string]*Cell)
	fc.Func_current = f
	fc.PC = 0
	// строим список переменных
	for i := 0; i < len(fc.Func_current.Var_list); i++ {
		fc.Var_value = append(fc.Var_value, Nil)
	}
	env.Call_Item[proc_id].Func_call_array = append(env.Call_Item[proc_id].Func_call_array, &fc)
	return len(env.Call_Item[proc_id].Func_call_array) - 1
}

func CheckFuncArgs(fc *Func_call, f_n int64) ([]*Cell, int) {
	ifi := G_FuncsInt_dict[f_n]
	args := []*Cell{}
	flag := false
	pos := 0
	if len(fc.Stack.Stack) > 0 {
		//fmt.Printf("len(fc.Stack.Stack) %v\r\n", len(fc.Stack.Stack))
		for i := len(fc.Stack.Stack) - 1; i > -1; i-- {
			//fmt.Printf("'%v'\r\n", fc.Stack.Stack[i].String(false))
			if fc.Stack.Stack[i] == Enter {
				//fmt.Printf("i %v\r\n", i)
				pos = i
				flag = true
				break
			}
		}
		//fmt.Printf("flag %v pos %v\r\n", flag, pos)
		if flag {
			// в стеке нашли Enter!!!! значит был применен выход через return
			// по очереди вытаскиваем аргументы
			for i := pos + 1; i < len(fc.Stack.Stack); i++ {
				args = append(args, fc.Stack.Stack[i])
			}
			fc.Stack.Stack = fc.Stack.Stack[:pos]
		} else {
			// в стеке Enter нет!
			return nil, -1
		}
	} else {
		// стек пуст!
		return nil, -2
	}
	//fmt.Printf("ifi.Args %v, len(args) %v, len(ifi.Args) %v, ifi %v\r\n", ifi.Args, len(args), len(ifi.Args), ifi)
	if len(args) < len(ifi.Args) {
		// недостаточно аргументов
		return nil, -3
	}
	return args, 0
}

/*
func (env *EnvironmentLVM) ToStack(current_fc int, args []*Cell) int {
	fc := env.Func_call_array[current_fc]

	for _, cc := range args {
		fc.Stack.Stack = append(fc.Stack.Stack, cc)
	}

	return 0
}
*/

func SetCurrentExec(env *EnvironmentLVM, current_fc int, args []*Cell, proc_id int) int {
	fc := env.Call_Item[proc_id].Func_call_array[current_fc]

	for _, cc := range args {
		fc.Stack.Stack = append(fc.Stack.Stack, cc)
	}
	env.Call_Item[proc_id].CurrentExec = append(env.Call_Item[proc_id].CurrentExec, current_fc)
	ia := Int_array{}
	env.Call_Item[proc_id].FuncStack = append(env.Call_Item[proc_id].FuncStack, &ia)
	return 0
}

func recoveryFunction() {
	if recoveryMessage:=recover(); recoveryMessage != nil {
		fmt.Println(recoveryMessage)
	}
	fmt.Println("This is recovery function...")
}

func Exec(env *EnvironmentLVM, proc_id int, debug bool) int {
        // defer recoveryFunction()
	// var finished chan bool
	for i := 0; i < len(env.Call_Item[proc_id].CurrentExec); i++ {
		current_fc := env.Call_Item[proc_id].CurrentExec[i]
		fc := env.Call_Item[proc_id].Func_call_array[current_fc]
		if debug {
			fmt.Printf("Exec fc.PC %v len(fc.Call_Item[proc_id].Func_current.Code) %v\r\n", fc.PC, len(fc.Func_current.Code))
		}
		// проверяем, что можно что то выполнить или это конец
		if fc.PC > len(fc.Func_current.Code)-1 {
			// закончено выполнение
			// результат лежит в стеке
			// ищем enter в стеке
			result := []*Cell{}
			flag := false
			pos := 0
			if len(fc.Stack.Stack) > 0 {
				// fmt.Printf("Find enter ")
				for i := len(fc.Stack.Stack) - 1; i > -1; i-- {
					// fmt.Printf(" %v", fc.Stack.Stack[i].String(false))
					if fc.Stack.Stack[i] == Enter {
						pos = i
						flag = true
						break
					}
				}
				// fmt.Printf("\r\n")
				// fmt.Printf("flag %v\r\n", flag)
				if flag {
					// в стеке нашли Enter!!!! значит был применен выход через return
					// по очереди вытаскиваем аргументы
					// fmt.Printf("Result 1")
					for i := pos + 1; i < len(fc.Stack.Stack); i++ {
						result = append(result, fc.Stack.Stack[i])
						// fmt.Printf(" %v", fc.Stack.Stack[i].String(false))
					}
					// fmt.Printf("\r\n")
					fc.Stack.Stack = fc.Stack.Stack[:pos]
				} else {
					// в стеке Enter нет! Значит возвращаемое значение - в голове стека
					// fmt.Printf("Result 2")
					vr := fc.Stack.Stack[len(fc.Stack.Stack)-1]
					if vr.Type == Cell_multiple {
						for i := 0; i < int(vr.Value_int); i++ {
							vrv := fc.Stack.Stack[len(fc.Stack.Stack)-2-i]
							result = append(result, vrv)
							//fmt.Printf(" %v", vrv)
						}
					} else {
						result = append(result, vr)
					}
					// fmt.Printf("%v\r\n", result)
				}
			} else {
				result = append(result, Nil)
			}

			// смотрим, есть ли возврат?
			p_ia := env.Call_Item[proc_id].FuncStack[i]
			ia := *p_ia
			// fmt.Printf("len(ia) %v\r\n", len(ia))
			if len(ia) > 0 {
				// стек функций не пуст!
				// удаляем старый
				//fmt.Printf("current_fc %v env.Func_call_array %v\r\n", current_fc, env.Func_call_array)
				if current_fc == 0 {
					env.Call_Item[proc_id].Func_call_array = env.Call_Item[proc_id].Func_call_array[1 : len(env.Call_Item[proc_id].Func_call_array)-1]
				} else if current_fc == len(env.Call_Item[proc_id].Func_call_array)-1 {
					env.Call_Item[proc_id].Func_call_array = env.Call_Item[proc_id].Func_call_array[:len(env.Call_Item[proc_id].Func_call_array)-1]
					//fmt.Printf("current_fc %v env.Func_call_array %v\r\n", current_fc, env.Func_call_array)
				} else {
					env.Call_Item[proc_id].Func_call_array = append(env.Call_Item[proc_id].Func_call_array[:current_fc-1], env.Call_Item[proc_id].Func_call_array[current_fc+1:len(env.Call_Item[proc_id].Func_call_array)-1]...)
				}
				// восстанавливаем
				current_fc_o := ia[len(ia)-1]
				if len(ia) > 1 {
					ia = ia[:len(ia)-2]
				} else {
					ia = ia[:len(ia)-1]
				}
				env.Call_Item[proc_id].FuncStack[i] = &ia
				current_fc = current_fc_o
				//fmt.Printf("i %v env.CurrentExec %v\r\n", i, env.CurrentExec)
				env.Call_Item[proc_id].CurrentExec[i] = current_fc

				fc = env.Call_Item[proc_id].Func_call_array[current_fc]
				if debug {
					fmt.Printf("return to fc.PC %v len(fc.Func_current.Code) %v\r\n", fc.PC, len(fc.Func_current.Code))
				}
				// добавляем в стек результаты возврата
				if debug {
					fmt.Printf("Result")
				}
				for _, v := range result {
					fc.Stack.Stack = append(fc.Stack.Stack, v)
					if debug {
						fmt.Printf(" %v", v.String(false))
					}
				}
				if debug {
					fmt.Printf(" end\r\n")
				}
				// и продолжаем исполнение
				continue
			} else {
				// заканчиваем исполнение
				if len(env.Call_Item[proc_id].Func_call_array) == 1 {
					env.Call_Item[proc_id].Func_call_array = env.Call_Item[proc_id].Func_call_array[:len(env.Call_Item[proc_id].Func_call_array)-1]
				} else {
					if current_fc == 0 {
						//fmt.Printf("current_fc %v, len(env.Func_call_array) %v\r\n", current_fc, len(env.Call_Item[proc_id].Func_call_array))
						env.Call_Item[proc_id].Func_call_array = env.Call_Item[proc_id].Func_call_array[1 : len(env.Call_Item[proc_id].Func_call_array)-1]
					} else if current_fc == len(env.Call_Item[proc_id].Func_call_array)-1 {
						env.Call_Item[proc_id].Func_call_array = env.Call_Item[proc_id].Func_call_array[:len(env.Call_Item[proc_id].Func_call_array)-1]
					} else {
						env.Call_Item[proc_id].Func_call_array = append(env.Call_Item[proc_id].Func_call_array[:current_fc-1], env.Call_Item[proc_id].Func_call_array[current_fc+1:len(env.Call_Item[proc_id].Func_call_array)-1]...)
					}
				}
				if len(env.Call_Item[proc_id].CurrentExec) == 1 {
					env.Call_Item[proc_id].CurrentExec = env.Call_Item[proc_id].CurrentExec[:len(env.Call_Item[proc_id].CurrentExec)-1]
				} else {
					if i == 0 {
						env.Call_Item[proc_id].CurrentExec = env.Call_Item[proc_id].CurrentExec[1 : len(env.Call_Item[proc_id].CurrentExec)-1]
					} else if i == len(env.Call_Item[proc_id].CurrentExec)-1 {
						env.Call_Item[proc_id].CurrentExec = env.Call_Item[proc_id].CurrentExec[:len(env.Call_Item[proc_id].CurrentExec)-1]
					} else {
						env.Call_Item[proc_id].CurrentExec = append(env.Call_Item[proc_id].CurrentExec[:i-1], env.Call_Item[proc_id].CurrentExec[i+1:len(env.Call_Item[proc_id].CurrentExec)-1]...)
					}
				}
				if len(env.Call_Item[proc_id].FuncStack) == 1 {
					env.Call_Item[proc_id].FuncStack = env.Call_Item[proc_id].FuncStack[:len(env.Call_Item[proc_id].FuncStack)-1]
				} else {
					if i == 0 {
						env.Call_Item[proc_id].FuncStack = env.Call_Item[proc_id].FuncStack[1 : len(env.Call_Item[proc_id].FuncStack)-1]
					} else if i == len(env.Call_Item[proc_id].FuncStack)-1 {
						env.Call_Item[proc_id].FuncStack = env.Call_Item[proc_id].FuncStack[:len(env.Call_Item[proc_id].FuncStack)-1]
					} else {
						env.Call_Item[proc_id].FuncStack = append(env.Call_Item[proc_id].FuncStack[:i-1], env.Call_Item[proc_id].FuncStack[i+1:len(env.Call_Item[proc_id].FuncStack)-1]...)
					}
				}
				// останавливаемся!!!! в том числе и что бы вернуть результат
				env.Call_Item[proc_id].Result = &result
				return 0
			}
		}
		// выполняем одну команду
		cc := fc.Func_current.Code[fc.PC]
		if debug {
			sc, ok := G_vmc_to_str[cc.Code]
			if !ok {

			}
			fmt.Printf("Cmd %v Const %v\r\n", sc, cc.Const)
		}
		switch cc.Code {
		case VMC_nop:
			fc.PC = fc.PC + 1
		case VMC_const:
			// загружает в стек по номеру
			// fmt.Printf("fc.Func_current.Const_list[%v] %v\r\n", cc.Const[0], fc.Func_current.Const_list[cc.Const[0]])
			fc.Stack.Stack = append(fc.Stack.Stack, fc.Func_current.Const_list[cc.Const[0]])
			fc.PC = fc.PC + 1
		case VMC_dup:
			// копирует в верхушке стека
			v := fc.Stack.Stack[len(fc.Stack.Stack)-1]
			fc.Stack.Stack = append(fc.Stack.Stack, v)
			fc.PC = fc.PC + 1
		case VMC_branch:
			// по значению аргумента переход
			pc := cc.Const[0]
			//fmt.Printf("pc %v\r\n", pc)
			fc.PC = fc.PC + int(pc) + 1
		case VMC_branch_true:
			// по значению аргумента переход
			v := fc.Stack.Stack[len(fc.Stack.Stack)-1]
			fc.Stack.Stack = fc.Stack.Stack[:len(fc.Stack.Stack)-1]
			if v == True {
				pc := cc.Const[0]
				//fmt.Printf("pc %v\r\n", pc)
				fc.PC = fc.PC + int(pc) + 1
			} else {
				fc.PC = fc.PC + 1
			}
		case VMC_branch_false:
			// по значению аргумента переход
			v := fc.Stack.Stack[len(fc.Stack.Stack)-1]
			fc.Stack.Stack = fc.Stack.Stack[:len(fc.Stack.Stack)-1]
			if v != True {
				pc := cc.Const[0]
				//fmt.Printf("pc %v\r\n", pc)
				fc.PC = fc.PC + int(pc) + 1
			} else {
				fc.PC = fc.PC + 1
			}
		case VMC_push:
			// из переменной значение в стек переменную указывает номер аргумента
			if false {
				fmt.Printf("Stack")
				for i := len(fc.Stack.Stack) - 1; i > -1; i-- {
					fmt.Printf(" '%v'", fc.Stack.Stack[i].String(false))
				}
				fmt.Printf(" !\r\n")
			}
			v := fc.Var_value[cc.Const[0]]
			fc.Stack.Stack = append(fc.Stack.Stack, v)
			fc.PC = fc.PC + 1
		case VMC_pop:
			// из вершины стека в переменную. переменную указывает номер аргумента
			if false {
				fmt.Printf("Stack")
				for i := len(fc.Stack.Stack) - 1; i > -1; i-- {
					fmt.Printf(" '%v'", fc.Stack.Stack[i].String(false))
				}
				fmt.Printf(" !\r\n")
			}
			v := fc.Stack.Stack[len(fc.Stack.Stack)-1]
			fc.Stack.Stack = fc.Stack.Stack[:len(fc.Stack.Stack)-1]
			fc.Var_value[cc.Const[0]] = v
			fc.PC = fc.PC + 1
		case VMC_enter:
			// формирует код начала вызова в стеке
			fc.Stack.Stack = append(fc.Stack.Stack, Enter)
			fc.PC = fc.PC + 1
		case VMC_call:
			// вызывает функцию. имя функции - первый элемент в стеке. строит новый стек вызова функции
			if false {
				fmt.Printf("Stack")
				for i := len(fc.Stack.Stack) - 1; i > -1; i-- {
					fmt.Printf(" '%v'", fc.Stack.Stack[i].String(false))
				}
				fmt.Printf(" !\r\n")
			}
			if cc.Cmd != nil {
				// вызываем запуск функции
				err := cc.Cmd(fc)
				if err != nil {
					// Ошибка!
				}
			} else {
				fn_num := cc.Const[0]
				args, res := CheckFuncArgs(fc, fn_num)
				if res != 0 {
					// ошибка!!!!
					fmt.Printf("Error: %v\r\n", res)
				}
				if false {
					fmt.Printf("fn_num %v\r\n", fn_num)
					for i, _ := range args {
						fmt.Printf("args %v\r\n", args[i].String(false))
					}
				}
				ff, ok := G_IntFuncDict[fn_num]
				if ok {
					err := ff(fc, args)
					if err != nil {
						fmt.Printf("Error %v\r\n", err)
					}
				} else {
					fmt.Printf("Error func id %v not found\r\n", fn_num)
				}
			}
			fc.PC = fc.PC + 1
		case VMC_e_call:
			// вызывает функцию. имя функции - первый элемент в стеке. строит новый стек вызова функции
			if false {
				fmt.Printf("Stack")
				for i := len(fc.Stack.Stack) - 1; i > -1; i-- {
					fmt.Printf(" '%v'", fc.Stack.Stack[i].String(false))
				}
				fmt.Printf(" !\r\n")
			}
			v := fc.Stack.Stack[len(fc.Stack.Stack)-1]
			fc.Stack.Stack = fc.Stack.Stack[:len(fc.Stack.Stack)-1]
			// ищем функцию
			//fmt.Printf("env.Func_dict %#v\r\n", env.Func_dict)
			//fmt.Printf("v.String(false) %v\r\n", v.String(false))
			f, ok := env.Func_dict[v.String(false)]
			if !ok {
				fmt.Printf("Func %v not found\r\n", v.String(false))
				return -1
			}
			// ищем Enter в стеке
			pos := 0
			for i := len(fc.Stack.Stack) - 1; i > -1; i-- {
				if fc.Stack.Stack[i] == Enter {
					pos = i
					break
				}
			}
			// по очереди вытаскиваем аргументы
			args_len := len(fc.Stack.Stack) - (pos + 1)
			args := make([]*Cell, args_len)
			j := args_len - 1
			for i := pos + 1; i < len(fc.Stack.Stack); i++ {
				args[j] = fc.Stack.Stack[i]
				j = j - 1
			}
			//fmt.Printf("%v\r\n", ss)
			fc.Stack.Stack = fc.Stack.Stack[:pos]
			fc.PC = fc.PC + 1
			// предыдущий добавляем в стек функций
			p_ia := env.Call_Item[proc_id].FuncStack[i]
			ia := *p_ia
			ia = append(ia, current_fc)
			env.Call_Item[proc_id].FuncStack[i] = &ia
			// строим новый
			current_fc_n := FuncCall(env, f, -1, proc_id)
			fc_n := env.Call_Item[proc_id].Func_call_array[current_fc_n]
			for _, cc := range args {
				fc_n.Stack.Stack = append(fc_n.Stack.Stack, cc)
			}
			env.Call_Item[proc_id].CurrentExec[i] = current_fc_n
			// далее идет выполнение
		case VMC_g_call:
			// вызывает функцию. имя функции - первый элемент в стеке. строит новый стек вызова функции
			if false {
				fmt.Printf("Stack")
				for i := len(fc.Stack.Stack) - 1; i > -1; i-- {
					fmt.Printf(" '%v'", fc.Stack.Stack[i].String(false))
				}
				fmt.Printf(" !\r\n")
			}
			v := fc.Stack.Stack[len(fc.Stack.Stack)-1]
			fc.Stack.Stack = fc.Stack.Stack[:len(fc.Stack.Stack)-1]
			// ищем функцию
			//fmt.Printf("env.Func_dict %#v\r\n", env.Func_dict)
			//fmt.Printf("v.String(false) %v\r\n", v.String(false))
			f, ok := env.Func_dict[v.String(false)]
			if !ok {
				fmt.Printf("Func %v not found\r\n", v.String(false))
				return -1
			}
			// ищем Enter в стеке
			pos := 0
			for i := len(fc.Stack.Stack) - 1; i > -1; i-- {
				if fc.Stack.Stack[i] == Enter {
					pos = i
					break
				}
			}
			// по очереди вытаскиваем аргументы
			args_len := len(fc.Stack.Stack) - (pos + 1)
			args := make([]*Cell, args_len)
			j := args_len - 1
			for i := pos + 1; i < len(fc.Stack.Stack); i++ {
				args[j] = fc.Stack.Stack[i]
				j = j - 1
			}
			//fmt.Printf("%v\r\n", ss)
			fc.Stack.Stack = fc.Stack.Stack[:pos]
			fc.PC = fc.PC + 1
			// строим новый вызов
			/*
				// предыдущий добавляем в стек функций
				p_ia := env.Call_Item[proc_id].FuncStack[i]
				ia := *p_ia
				ia = append(ia, current_fc)
				env.Call_Item[proc_id].FuncStack[i] = &ia
			*/

			// строим новый
			/*
				current_fc_n := FuncCall(env, f, -1, proc_id)

				fc_n := env.Call_Item[proc_id].Func_call_array[current_fc_n]
				for _, cc := range args {
					fc_n.Stack.Stack = append(fc_n.Stack.Stack, cc)
				}
				env.Call_Item[proc_id].CurrentExec[i] = current_fc_n
			*/
			// далее идет выполнение
			// finished = make(chan bool)
			go ProcExec(env, f, args, env.Call_Item[proc_id].Wg, -1, false, debug)
		}
	}
	// <-finished
	return 1
}

func ProcExec(e *EnvironmentLVM, f *Func, args []*Cell, wg *sync.WaitGroup, proc_id int, no_result bool, debug bool) {
	if len(e.Call_Item) == 0 {
		proc_id = e.Current_proc
		ci := CallItem{}
		e.Call_Item = append(e.Call_Item, ci)
	}
	if proc_id != e.Current_proc {
		// create new
		ci := CallItem{}
		e.Call_Item = append(e.Call_Item, ci)
		proc_id = len(e.Call_Item) - 1
	}
	current_fc := FuncCall(e, f, -1, proc_id)

	SetCurrentExec(e, current_fc, args, proc_id)
	wg.Add(1)
	e.Call_Item[proc_id].Wg = wg
	for {
		res := Exec(e, proc_id, debug)
		if res == 0 {
			if no_result {
				// печатаем результат
				result := *e.Call_Item[proc_id].Result
				for _, v := range result {
					ss := v.String(false)
					fmt.Printf("result %v\r\n", ss)
				}
			}
			break
		} else {
		if res == -1 {
		        fmt.Printf("Error\r\n")
		        break
		}
		}
	}
	// e.Call_Item[proc_id].Finished <- true
	// finished <- true
	wg.Done()
}

type T_IntFunc func(fc *Func_call, args []*Cell) error

var G_IntFuncDict map[int64]T_IntFunc
var G_IntFuncNameDict map[int64]string

func InitIntFunc() {
	G_IntFuncDict = make(map[int64]T_IntFunc)
	G_IntFuncNameDict = make(map[int64]string)

	G_IntFuncDict[FC_car] = IF_Car
	G_IntFuncDict[FC_cdr] = IF_Cdr
	G_IntFuncDict[FC_cons] = IF_Cons
	G_IntFuncDict[FC_type] = IF_Type
	G_IntFuncDict[FC_typep] = IF_TypeP
	G_IntFuncDict[FC_listp] = nil
	G_IntFuncDict[FC_eq] = IF_Eq
	G_IntFuncDict[FC_rplaca] = IF_Rplaca
	G_IntFuncDict[FC_rplacd] = IF_Rplacd
	G_IntFuncDict[FC_list] = nil
	G_IntFuncDict[FC_print] = IF_Print
	G_IntFuncDict[FC_add] = IF_Add
	G_IntFuncDict[FC_sub] = IF_Sub
	G_IntFuncDict[FC_mul] = IF_Mul
	G_IntFuncDict[FC_div] = IF_Div
	G_IntFuncDict[FC_2str] = IF_ToStr
	G_IntFuncDict[FC_2int] = IF_ToInt
	G_IntFuncDict[FC_2float] = IF_ToFloat
	G_IntFuncDict[FC_assert] = IF_Assert
	G_IntFuncDict[FC_make_array] = IF_MakeArray
	G_IntFuncDict[FC_make_dict] = IF_MakeDict
	G_IntFuncDict[FC_append] = IF_Append
	G_IntFuncDict[FC_slice] = IF_SliceArray
	//	G_IntFuncDict[FC_insert] = IF_InsertArray
	G_IntFuncDict[FC_item] = IF_Item
	G_IntFuncDict[FC_import] = IF_Import
	G_IntFuncDict[FC_strip] = IF_Strip
	G_IntFuncDict[FC_split] = IF_Split
	G_IntFuncDict[FC_index] = IF_Index
	G_IntFuncDict[FC_length] = IF_Length
	G_IntFuncDict[FC_join] = IF_Join
	G_IntFuncDict[FC_set_dict] = IF_SetDict
	G_IntFuncDict[FC_load_lib] = IF_LoadLib
	G_IntFuncDict[FC_call_lib_func] = IF_CallLibFunc
	G_IntFuncDict[FC_princ] = IF_Princ
	G_IntFuncDict[FC_format] = IF_Format
	G_IntFuncDict[FC_iterate] = IF_Iterate
	G_IntFuncDict[FC_get_iter] = IF_GetIter
	G_IntFuncDict[FC_make_iter] = IF_MakeIter
	G_IntFuncDict[FC_not] = IF_Not
	G_IntFuncDict[FC_find_func] = IF_FindFunc
	//
	G_IntFuncNameDict[FC_car] = "Car"
	G_IntFuncNameDict[FC_cdr] = "Cdr"
	G_IntFuncNameDict[FC_cons] = "Cons"
	G_IntFuncNameDict[FC_type] = "Type"
	G_IntFuncNameDict[FC_typep] = "TypeP"
	G_IntFuncNameDict[FC_listp] = ""
	G_IntFuncNameDict[FC_eq] = "Eq"
	G_IntFuncNameDict[FC_rplaca] = "Rplaca"
	G_IntFuncNameDict[FC_rplacd] = "Rplacd"
	G_IntFuncNameDict[FC_list] = ""
	G_IntFuncNameDict[FC_print] = "Print"
	G_IntFuncNameDict[FC_add] = "Add"
	G_IntFuncNameDict[FC_sub] = "Sub"
	G_IntFuncNameDict[FC_mul] = "Mul"
	G_IntFuncNameDict[FC_div] = "Div"
	G_IntFuncNameDict[FC_2str] = "ToStr"
	G_IntFuncNameDict[FC_2int] = "ToInt"
	G_IntFuncNameDict[FC_2float] = "ToFloat"
	G_IntFuncNameDict[FC_assert] = "Assert"
	G_IntFuncNameDict[FC_make_array] = "MakeArray"
	G_IntFuncNameDict[FC_make_dict] = "MakeDict"
	G_IntFuncNameDict[FC_append] = "Append"
	G_IntFuncNameDict[FC_slice] = "SliceArray"
	//	G_IntFuncNameDict[FC_insert] = "InsertArray"
	G_IntFuncNameDict[FC_item] = "Item"
	G_IntFuncNameDict[FC_import] = "Import"
	G_IntFuncNameDict[FC_strip] = "Strip"
	G_IntFuncNameDict[FC_split] = "Split"
	G_IntFuncNameDict[FC_index] = "Index"
	G_IntFuncNameDict[FC_length] = "Length"
	G_IntFuncNameDict[FC_join] = "Join"
	G_IntFuncNameDict[FC_set_dict] = "SetDict"
	G_IntFuncNameDict[FC_load_lib] = "LoadLib"
	G_IntFuncNameDict[FC_call_lib_func] = "CallLibFunc"
	G_IntFuncNameDict[FC_princ] = "Princ"
	G_IntFuncNameDict[FC_format] = "Format"
	G_IntFuncNameDict[FC_iterate] = "Iterate"
	G_IntFuncNameDict[FC_get_iter] = "GetIter"
	G_IntFuncNameDict[FC_make_iter] = "MakeIter"
	G_IntFuncNameDict[FC_not] = "Not"
	G_IntFuncNameDict[FC_find_func] = "FindFunc"
}

func AddIntFunc(name string, id int64, ff T_IntFunc) error {
	// надо добавить проверку что такой функции нет
	if _, ok := G_IntFuncDict[id]; ok {
		return errors.New("function exist")
	} else {
		G_IntFuncDict[id] = ff
		G_IntFuncNameDict[id] = name
	}
	return nil
}

func IF_Print(fc *Func_call, args []*Cell) error {
	// по очереди вытаскиваем аргументы
	ss := ""
	var result *Cell
	for _, v := range args {
		if len(ss) == 0 {
			ss = v.String(false)
		} else {
			ss = ss + " " + v.String(false)
		}
		result = v
	}
	if fc.Env.NoPrint {
		fmt.Printf("%v\r\n", ss)
	}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_Princ(fc *Func_call, args []*Cell) error {
	// по очереди вытаскиваем аргументы
	ss := ""
	var result *Cell
	for _, v := range args {
		if len(ss) == 0 {
			ss = v.String(false)
		} else {
			ss = ss + " " + v.String(false)
		}
		result = v
	}
	if fc.Env.NoPrint {
		fmt.Printf("%v", ss)
	}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

// variable
func IF_Format(fc *Func_call, args []*Cell) error {
	// по очереди вытаскиваем аргументы
	if len(args) < 1 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	v := args[0]
	form := ""
	switch v.Type {
	case Cell_str:
		form = v.Value_str
	default:
		// error
	}

	Format := func(a ...interface{}) string {
		ss := fmt.Sprintf(form, a)
		return ss
	}

	fun := reflect.ValueOf(Format)
	in := make([]reflect.Value, len(args)-1)

	var result *Cell
	for i, v := range args {
		if i > 0 {
			param := v.String(false)
			in[i-1] = reflect.ValueOf(param)
		}
	}
	r := fun.Call(in)
	sss := ""
	for _, v := range r {
		sss = sss + fmt.Sprintf("%v", v)
	}

	result = &Cell{Type: Cell_str, Value_str: sss}

	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_Add(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов не менее 2
	//fmt.Printf("len(args) %v\r\n", len(args))
	if len(args) < 2 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	ss_sym := ""
	ss_str := ""
	ss_int := (int64)(0)
	ss_float := 0.0
	var result *Cell
	t := -1
	for _, v := range args {
		//fmt.Printf("t %v v.Type %v\r\n", t, v.Type)
		if t != -1 {
			if t != Cell_sym {

			} else {
				if t != v.Type {
					var ErrTypeError = errors.New("type error")
					return ErrTypeError
				}
			}
		} else {
			t = v.Type
		}
		switch v.Type {
		case Cell_sym:
			switch v.Type {
			case Cell_sym:
				ss_sym = ss_sym + v.Value_sym
			case Cell_int:
				ss_sym = ss_sym + fmt.Sprintf("%v", v.Value_int)
			case Cell_str:
				ss_sym = ss_sym + v.Value_str
			case Cell_float:
				ss_sym = ss_sym + fmt.Sprintf("%f", v.Value_float)
			}
		case Cell_int:
			ss_int = ss_int + v.Value_int
		case Cell_str:
			ss_str = ss_str + v.Value_str
		case Cell_float:
			ss_float = ss_float + v.Value_float
		}
	}
	//fmt.Printf("t %#v\r\n", t)
	switch t {
	case Cell_sym:
		result = &Cell{Type: Cell_sym, Value_sym: ss_sym}
	case Cell_int:
		result = &Cell{Type: Cell_int, Value_int: ss_int}
	case Cell_str:
		result = &Cell{Type: Cell_str, Value_str: ss_str}
	case Cell_float:
		result = &Cell{Type: Cell_float, Value_float: ss_float}
	}
	//fmt.Printf("result %#v\r\n", result)
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_Sub(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов не менее 2
	if len(args) < 2 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	//ss_str := ""
	ss_int := (int64)(0)
	ss_float := 0.0
	var result *Cell
	t := -1
	for i, v := range args {
		if t != -1 {
			if t != v.Type {
				var ErrTypeError = errors.New("type error")
				return ErrTypeError
			}
		} else {
			t = v.Type
		}
		if i == 0 {
			switch v.Type {
			case Cell_int:
				ss_int = v.Value_int
				/*
					case Cell_str:
						ss_str = v.Value_str
				*/
			case Cell_float:
				ss_float = v.Value_float
			}
		} else {
			switch v.Type {
			case Cell_int:
				ss_int = ss_int - v.Value_int
				/*
					case Cell_str:
						ss_str = ss_str - v.Value_str
				*/
			case Cell_float:
				ss_float = ss_float - v.Value_float
			}
		}
	}
	// fmt.Printf("ss_int %v\r\n", ss_int)
	switch t {
	case Cell_int:
		result = &Cell{Type: Cell_int, Value_int: ss_int}
		/*
			case Cell_str:
				result = &Cell{Cell_str, 0, 0, "", ss_str, nil, nil, nil, nil, nil}
		*/
	case Cell_float:
		result = &Cell{Type: Cell_float, Value_float: ss_float}
	}
	//fmt.Printf("%v\r\n", ss)
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_Mul(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов не менее 2
	if len(args) < 2 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	//	ss_str := ""
	ss_int := (int64)(0)
	ss_float := 0.0
	var result *Cell
	t := -1
	for i, v := range args {
		if t != -1 {
			if t != v.Type {
				var ErrTypeError = errors.New("type error")
				return ErrTypeError
			}
		} else {
			t = v.Type
		}
		if i == 0 {
			switch v.Type {
			case Cell_int:
				ss_int = v.Value_int
				/*
					case Cell_str:
						ss_str = ss_str * v.Value_str
				*/
			case Cell_float:
				ss_float = v.Value_float
			}
		} else {
			switch v.Type {
			case Cell_int:
				ss_int = ss_int * v.Value_int
				/*
					case Cell_str:
						ss_str = ss_str * v.Value_str
				*/
			case Cell_float:
				ss_float = ss_float * v.Value_float
			}
		}
	}
	switch t {
	case Cell_int:
		result = &Cell{Type: Cell_int, Value_int: ss_int}
		/*
			case Cell_str:
				result = &Cell{Type:Cell_str, Value_str:ss_str}
		*/
	case Cell_float:
		result = &Cell{Type: Cell_float, Value_float: ss_float}
	}
	//fmt.Printf("%v\r\n", ss)
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_Div(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов не менее 2
	if len(args) < 2 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	//	ss_str := ""
	ss_int := (int64)(0)
	ss_float := 0.0
	var result *Cell
	t := -1
	for i, v := range args {
		if t != -1 {
			if t != v.Type {
				var ErrTypeError = errors.New("type error")
				return ErrTypeError
			}
		} else {
			t = v.Type
		}
		if i == 0 {
			switch v.Type {
			case Cell_int:
				ss_int = v.Value_int
				/*
					case Cell_str:
						ss_str = v.Value_str
				*/
			case Cell_float:
				ss_float = v.Value_float
			}
		} else {
			switch v.Type {
			case Cell_int:
				ss_int = ss_int / v.Value_int
				/*
					case Cell_str:
						ss_str = ss_str / v.Value_str
				*/
			case Cell_float:
				ss_float = ss_float / v.Value_float
			}
		}
	}
	switch t {
	case Cell_int:
		result = &Cell{Type: Cell_int, Value_int: ss_int}
		/*
			case Cell_str:
				result = &Cell{Type:Cell_str, Value_str:ss_str}
		*/
	case Cell_float:
		result = &Cell{Type: Cell_float, Value_float: ss_float}
	}
	//fmt.Printf("%v\r\n", ss)
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_Type(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов 1
	if len(args) < 1 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 1 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	v := args[0]
	//	ss_str := ""
	var result *Cell
	switch v.Type {
	case Cell_int:
		result = T_int
	case Cell_float:
		result = T_float
	case Cell_sym:
		result = T_sym
	case Cell_str:
		result = T_str
	case Cell_cell:
		result = T_cell
	case Cell_dict:
		result = T_dict
	case Cell_array:
		result = T_array
	case Cell_func:
		result = T_func
	case Cell_ext:
		result = T_ext
	case Cell_channel:
		result = T_channel
	case Cell_next:
		result = T_next
	}
	//fmt.Printf("result type %v\r\n", result.String(false))
	// result = &Cell{Type: Cell_str, Value_str: ss_str}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_TypeP(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов 2
	if len(args) < 2 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 2 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	v := args[0]
	t := args[1]
	// fmt.Printf("-> %v %v\r\n", v.String(false), t.String(false))
	ss_str := ""
	switch t.Type {
	case Cell_str:
		ss_str = t.Value_str
	case Cell_sym:
		ss_str = t.Value_sym
	}
	// fmt.Printf("--> %v\r\n", ss_str)
	result := False
	switch v.Type {
	case Cell_int:
		if ss_str == "int" {
			result = True
		}
	case Cell_float:
		if ss_str == "float" {
			result = True
		}
	case Cell_sym:
		if ss_str == "sym" {
			result = True
		}
	case Cell_str:
		if ss_str == "str" {
			result = True
		}
	case Cell_cell:
		if ss_str == "cell" {
			result = True
		}
	case Cell_dict:
		if ss_str == "dict" {
			result = True
		}
	case Cell_array:
		if ss_str == "array" {
			result = True
		}
	case Cell_func:
		if ss_str == "func" {
			result = True
		}
	case Cell_ext:
		if ss_str == "ext" {
			result = True
		}
	case Cell_channel:
		if ss_str == "channel" {
			result = True
		}
	case Cell_next:
		if ss_str == "new_ext" {
			result = True
		}
	}
	// fmt.Printf("result %v\r\n", result)
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_ToStr(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов 1
	if len(args) < 1 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 1 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	v := args[0]
	ss_str := ""
	//ss_int := (int64)(0)
	//ss_float := 0.0
	var result *Cell
	switch v.Type {
	case Cell_int:
		ss_str = fmt.Sprintf("%v", v.Value_int)
	case Cell_str:
		ss_str = v.Value_str
	case Cell_float:
		ss_str = fmt.Sprintf("%v", v.Value_float)
	}
	result = &Cell{Type: Cell_str, Value_str: ss_str}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_ToInt(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов 1
	if len(args) < 1 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 1 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	v := args[0]
	//ss_str := ""
	ss_int := (int64)(0)
	//ss_float := 0.0
	var result *Cell
	switch v.Type {
	case Cell_int:
		ss_int = v.Value_int
	case Cell_str:
		ss_int, _ = strconv.ParseInt(v.Value_str, 0, 64)
	case Cell_float:
		ss_int = int64(v.Value_float)
	}
	result = &Cell{Type: Cell_int, Value_int: ss_int}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_ToFloat(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов 1
	if len(args) < 1 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 1 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	v := args[0]
	//ss_str := ""
	//ss_int := (int64)(0)
	ss_float := 0.0
	var result *Cell
	switch v.Type {
	case Cell_int:
		ss_float = float64(v.Value_int)
	case Cell_str:
		ss_float, _ = strconv.ParseFloat(v.Value_str, 64)
	case Cell_float:
		ss_float = v.Value_float
	}
	result = &Cell{Type: Cell_float, Value_float: ss_float}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_Car(fc *Func_call, args []*Cell) error {
	// по очереди вытаскиваем аргументы
	var result *Cell
	v := args[0]
	switch v.Type {
	case Cell_cell:
		result = v.Value_head
	}

	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_Cdr(fc *Func_call, args []*Cell) error {
	// по очереди вытаскиваем аргументы
	var result *Cell
	v := args[0]
	switch v.Type {
	case Cell_cell:
		result = v.Value_last
	}

	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_Cons(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов - 2
	if len(args) < 2 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 2 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	var result *Cell
	v_h := args[0]
	v_l := args[1]
	result = &Cell{Type: Cell_cell, Value_head: v_h, Value_last: v_l}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_Eq(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов - 2
	if len(args) < 2 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 2 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	a1 := args[0]
	a2 := args[1]
	//fmt.Printf("a1.Type %v a2.Type %v\r\n", a1.Type, a2.Type)
	if a1.Type != a2.Type {
		fc.Stack.Stack = append(fc.Stack.Stack, False)
	} else {
		switch a1.Type {
		case Cell_sym:
			if a1.Value_sym != a2.Value_sym {
				fc.Stack.Stack = append(fc.Stack.Stack, False)
			} else {
				fc.Stack.Stack = append(fc.Stack.Stack, True)
			}
		case Cell_int:
			if a1.Value_int != a2.Value_int {
				fc.Stack.Stack = append(fc.Stack.Stack, False)
			} else {
				fc.Stack.Stack = append(fc.Stack.Stack, True)
			}
		case Cell_str:
			//fmt.Printf("a1.Value_str %v a2.Value_str %v\r\b", a1.Value_str, a2.Value_str)
			if a1.Value_str != a2.Value_str {
				fc.Stack.Stack = append(fc.Stack.Stack, False)
			} else {
				fc.Stack.Stack = append(fc.Stack.Stack, True)
			}
		case Cell_float:
			if a1.Value_float != a2.Value_float {
				fc.Stack.Stack = append(fc.Stack.Stack, False)
			} else {
				fc.Stack.Stack = append(fc.Stack.Stack, True)
			}
		}

	}
	return nil
}

func IF_Not(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов - 1
	if len(args) < 1 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 1 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	a1 := args[0]
	// a2 := args[1]
	/*
		if a1.Type != a2.Type {
			fc.Stack.Stack = append(fc.Stack.Stack, False)
		} else {
	*/
	//ss := a1.String(false)
	//fmt.Printf(">> a1 %v\r\n", ss)
	if a1 == Nil {
		fc.Stack.Stack = append(fc.Stack.Stack, True)
	} else {
		if a1 != False {
			fc.Stack.Stack = append(fc.Stack.Stack, False)
		} else {
			fc.Stack.Stack = append(fc.Stack.Stack, True)
		}
	}
	/*
		}
	*/
	return nil
}

func IF_Rplaca(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов - 2
	if len(args) < 2 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 2 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	c := args[0]
	v := args[1]
	switch c.Type {
	case Cell_cell:
		c.Value_head = v
	}
	fc.Stack.Stack = append(fc.Stack.Stack, c)

	return nil
}

func IF_Rplacd(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов - 2
	if len(args) < 2 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 2 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	c := args[0]
	v := args[1]
	switch c.Type {
	case Cell_cell:
		c.Value_last = v
	}
	fc.Stack.Stack = append(fc.Stack.Stack, c)

	return nil
}

func IF_Assert(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов 2
	if len(args) < 2 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 2 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	v := args[0]
	t := args[1]
	// проверяем, что результат совпадает
	// fmt.Printf("t.Type %v v.Type %v\r\n", t.Type, v.Type)
	ss_int := int64(0)
	ss_str := ""
	ss_sym := ""
	ss_float := 0.0
	switch t.Type {
	case Cell_int:
		ss_int = t.Value_int
	case Cell_float:
		ss_float = t.Value_float
	case Cell_sym:
		ss_sym = t.Value_sym
	case Cell_str:
		ss_str = t.Value_str
	}

	result := False
	switch v.Type {
	case Cell_int:
		//fmt.Printf("ss_int %v v.Value_int %v\r\n", ss_int, v.Value_int)
		if ss_int == v.Value_int {
			result = True
		}
	case Cell_float:
		//fmt.Printf("ss_float %v v.Value_float %v\r\n", ss_float, v.Value_float)
		if ss_float == v.Value_float {
			result = True
		}
	case Cell_sym:
		//fmt.Printf("ss_sym %v v.Value_sym %v\r\n", ss_sym, v.Value_sym)
		if ss_sym == v.Value_sym {
			result = True
		}
	case Cell_str:
		//fmt.Printf("ss_str %v v.Value_str %v\r\n", ss_str, v.Value_str)
		if ss_str == v.Value_str {
			result = True
		}
	}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_MakeArray(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов 1
	if len(args) < 1 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 1 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	v := args[0]
	a := []*Cell{}
	var result *Cell
	switch v.Type {
	case Cell_int:
		if v.Value_int > 0 {
			for i := 0; i < int(v.Value_int); i++ {
				a = append(a, Nil)
			}
		}
	}

	result = &Cell{Type: Cell_array, Value_array: a}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_MakeDict(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов 1
	if len(args) < 2 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 2 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	var result *Cell
	m := make(map[string]*Cell)
	/*
		v := args[0]
		k := args[1]
		vt := false
		switch v.Type {
		case Cell_cell:
			vt = true
		}
		kt := false
		switch k.Type {
		case Cell_cell:
			kt = true
		}
	*/
	result = &Cell{Type: Cell_dict, Value_dict: m}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_Append(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов 1
	if len(args) < 2 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 2 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	var result *Cell
	v := args[0]
	k := args[1]
	switch v.Type {
	case Cell_array:
		v.Value_array = append(v.Value_array, k)
		result = v
	default:
		result = Nil
	}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_SliceArray(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов 3
	if len(args) < 3 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 3 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	var result *Cell
	v := args[0]
	b := args[1]
	e := args[2]
	bi := 0
	ei := 0
	switch b.Type {
	case Cell_int:
		bi = int(b.Value_int)
	}
	switch e.Type {
	case Cell_int:
		ei = int(e.Value_int)
	}
	switch v.Type {
	case Cell_sym:
		if bi < 0 {
			bi = 0
		}
		if ei == 0 {
			ei = len(v.Value_sym)
		}
		if ei > len(v.Value_sym) {
			ei = len(v.Value_sym)
		}
		result = &Cell{Type: Cell_str, Value_sym: v.Value_sym[bi:ei]}
	case Cell_str:
		if bi < 0 {
			bi = 0
		}
		if ei == 0 {
			ei = len(v.Value_str)
		}
		if ei > len(v.Value_str) {
			ei = len(v.Value_str)
		}
		result = &Cell{Type: Cell_str, Value_str: v.Value_str[bi:ei]}
	case Cell_array:
		if bi < 0 {
			bi = 0
		}
		if ei == 0 {
			ei = len(v.Value_array)
		}
		if ei > len(v.Value_array) {
			ei = len(v.Value_array)
		}
		fmt.Printf("len(v.Value_array) %v, bi %v, ei %v\r\n", len(v.Value_array), bi, ei)
		result = &Cell{Type: Cell_array, Value_array: v.Value_array[bi:ei]}
	default:
		result = Nil
	}

	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_Item(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов 2
	if len(args) < 2 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 2 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	var result *Cell
	v := args[0]
	i := args[1]
	switch v.Type {
	case Cell_array:
		ii := 0
		switch i.Type {
		case Cell_int:
			ii = int(i.Value_int)
		}
		result = v.Value_array[ii]
	case Cell_dict:
		ii := ""
		switch i.Type {
		case Cell_str:
			ii = i.Value_str
		case Cell_sym:
			ii = i.Value_sym
		}
		result = v.Value_dict[ii]
	default:
		result = Nil
	}

	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_Length(fc *Func_call, args []*Cell) error {
	// проверяем, что аргумент 1
	if len(args) < 1 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 1 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	var result *Cell
	v := args[0]
	ii := 0
	switch v.Type {
	case Cell_str:
		ii = len(v.Value_str)
	case Cell_array:
		ii = len(v.Value_array)
	case Cell_dict:
		ii = len(v.Value_dict)
	default:
		result = Nil
	}
	if result != Nil {
		result = &Cell{Type: Cell_int, Value_int: int64(ii)}
	}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_Import(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов 1
	if len(args) < 1 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 1 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	v := args[0]
	var result *Cell = False
	switch v.Type {
	case Cell_str:
		filename := v.Value_str
		fmt.Printf("filename %v\r\n", filename)
		// как грузим?
		// смотрим файл с расширением comp
		// есди нету - смотрим с расширением lisp
		// считывание
		ff_ext := fc.Env.FileFormat
		fn := filename + "." + ff_ext
		fmt.Printf("fn %v\r\n", fn)
		var ms Module_store
		//ms, err := LoadModule(fn) // ".comp"
		// считывание
		if ff_ext == "comp" {
			ms_, err := LoadModule(fn)
			if err != nil {
				fmt.Printf("%v\r\n", err)
			}
			ms = ms_
		} else {
			if ff_ext == "bin" {
				data, err := ioutil.ReadFile(fn)
				if err != nil {
					fmt.Print(err)
					return nil
				}
	        
                                ms_, bb_n, err := Load_module_store(data, 0)
                                if err != nil {
                                        fmt.Printf("Error Load_module_store %v\r\n", err)
                                }
                                ms = *ms_
                                //fmt.Printf("%#v\r\n", *ms_)
                                if len(bb_n) > 0 {
					fmt.Printf("Error in file - remainder %v\r\n", bb_n)
					return nil
                                }
                                //fmt.Printf("bb_n %v\r\n", bb_n)
			} else {
			}
		}

/*
		if err != nil {
			fmt.Printf("%v\r\n", err)
		} else {
*/
			for _, ff := range ms.Func_list {
				// fmt.Printf("%v\r\n", ff)
				f2 := Func_store2Func(ff)
				// fmt.Printf("%v\r\n", f2)
				fc.Env.FuncToEnv(&f2)
			}
//		}

		result = True
	default:
		result = Nil
	}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_Strip(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов 2
	if len(args) < 2 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 2 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	var result *Cell
	s := args[0]
	t := args[1]
	// проверяем, что результат совпадает
	//fmt.Printf("t.Type %v v.Type %v\r\n", t.Type, v.Type)
	ss_str := ""
	ss_str_t := ""
	switch s.Type {
	case Cell_str:
		ss_str = s.Value_str
	}
	switch t.Type {
	case Cell_str:
		ss_str_t = t.Value_str
	}
	//fmt.Printf("ss_str '%v', ss_str_t '%v'\r\n", ss_str, ss_str_t)
	ss_str = strings.Trim(ss_str, ss_str_t)
	//fmt.Printf("ss_str '%v'\r\n", ss_str)
	result = &Cell{Type: Cell_str, Value_str: ss_str}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_Split(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов 2
	if len(args) < 2 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 2 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	var result *Cell
	s := args[0]
	t := args[1]
	// проверяем, что результат совпадает
	//fmt.Printf("t.Type %v v.Type %v\r\n", t.Type, v.Type)
	ss_str := ""
	ss_str_t := ""
	switch s.Type {
	case Cell_str:
		ss_str = s.Value_str
	}
	switch t.Type {
	case Cell_str:
		ss_str_t = t.Value_str
	}
	str_lst := strings.Split(ss_str, ss_str_t)
	var ra []*Cell
	for i, _ := range str_lst {
		sc := &Cell{Type: Cell_str, Value_str: str_lst[i]}
		ra = append(ra, sc)
	}
	result = &Cell{Type: Cell_array, Value_array: ra}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_Index(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов 2
	if len(args) < 2 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 2 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	var result *Cell
	s := args[0]
	t := args[1]
	// проверяем, что результат совпадает
	//fmt.Printf("t.Type %v v.Type %v\r\n", t.Type, v.Type)
	ss_str := ""
	ss_str_t := ""
	switch s.Type {
	case Cell_str:
		ss_str = s.Value_str
	}
	switch t.Type {
	case Cell_str:
		ss_str_t = t.Value_str
	}
	inx := strings.Index(ss_str, ss_str_t)
	result = &Cell{Type: Cell_int, Value_int: int64(inx)}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_Join(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов 2
	if len(args) < 2 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 2 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	var result *Cell
	v := args[0]
	t := args[1]
	// проверяем, что результат совпадает
	//fmt.Printf("t.Type %v v.Type %v\r\n", t.Type, v.Type)
	ss_str := ""
	ss_str_t := ""
	switch t.Type {
	case Cell_str:
		ss_str_t = t.Value_str
	}
	switch v.Type {
	case Cell_array:
		for i, _ := range v.Value_array {
			ss := ""
			vi := v.Value_array[i]
			switch vi.Type {
			case Cell_str:
				ss = vi.Value_str
			default:
				ss = vi.String(false)
			}
			if i > 0 {
				ss_str = ss_str + ss_str_t + ss
			} else {
				ss_str = ss
			}
		}
	}
	result = &Cell{Type: Cell_str, Value_str: ss_str}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_SetDict(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов 3
	if len(args) < 3 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 3 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	var result *Cell
	d := args[0]
	i := args[1]
	v := args[2]
	switch d.Type {
	case Cell_dict:
		ii := ""
		switch i.Type {
		case Cell_str:
			ii = i.Value_str
		case Cell_sym:
			ii = i.Value_sym
		}
		d.Value_dict[ii] = v
		result = d
	default:
		result = Nil
	}

	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

// !!!!
var LibDict map[string]*ExtLib

func init() {
	LibDict = make(map[string]*ExtLib)
}

func IF_LoadLib(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов 3
	if len(args) < 1 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 1 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	var result *Cell
	d := args[0]
	module_name := ""
	switch d.Type {
	case Cell_str:
		module_name = d.Value_str
	case Cell_sym:
		module_name = d.Value_sym
	}
	if module_name != "" {
		el, err1 := LoadExtLib(module_name)
		if err1 != nil {
			fmt.Printf("LoadExtLib %v\r\n", err1)
			result = Nil
		} else {
			result = d
			LibDict[module_name] = el
			/*
			   fd := fmt.Sprintf("%v", el)
			   id := []byte{}
			   result = CellCreateExt(id, "дескриптор_файла", []string{fd})
			*/
		}
	}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func Cell2Short(c *Cell) *CellShort {
	s := CellShort{}
	s.Type = int32(c.Type)
	s.Value_int = c.Value_int
	s.Value_float = c.Value_float
	s.Value_sym = c.Value_sym
	s.Value_str = c.Value_str
	//s.Value_head = c.Value_head
	//s.Value_last = c.Value_last
	//s.Value_dict = c.Value_dict
	//s.Value_array= c.Value_array
	return &s
}

func Short2Cell(s *CellShort) *Cell {
	c := Cell{}
	c.Type = int(s.Type)
	c.Value_int = s.Value_int
	c.Value_float = s.Value_float
	c.Value_sym = s.Value_sym
	c.Value_str = s.Value_str
	//c.Value_head = s.Value_head
	//c.Value_last = s.Value_last
	//c.Value_dict = s.Value_dict
	//c.Value_array= s.Value_array
	return &c
}

func IF_CallLibFunc(fc *Func_call, args []*Cell) error {
	// проверяем, что аргументов 3
	if len(args) < 3 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 3 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	var result *Cell
	d := args[0]
	f := args[1]
	al := args[2]
	//var fd int64
	func_name := ""
	data_b := []CellShort{}
	module_name := ""
	var g_el *ExtLib
	switch d.Type {
	case Cell_str:
		module_name = d.Value_str
	case Cell_sym:
		module_name = d.Value_sym
	}
	if module_name != "" {
		//      fmt.Printf("module_name %v\r\n", module_name)
		el, ok := LibDict[module_name]
		if !ok {
			errors.New(fmt.Sprintf("lib %v not loaded", module_name))
		}
		g_el = el
	} else {
		return errors.New("wrong lib name")
	}
	/*
	   	switch d.Type {
	   		case Cell_ext:
	   			_, name, value, b := d.GetCellExtValue()
	   			if !b {
	   				return errors.New("ext type value error")
	   			}
	   			if name == "дескриптор_файла" {
	   				fd, _ = strconv.ParseInt(value[0], 0, 64)
	   			} else {
	   				return errors.New("wrong type - expected 'дескриптор_файла'")
	   			}
	           		cs_in := CellShort{}
	           		cs_in.Type = 2
	           		cs_in.Value_int = fd
	                           data_b = append(data_b, cs_in)
	   	}
	*/
	switch f.Type {
	case Cell_str:
		func_name = f.Value_str
	case Cell_sym:
		func_name = f.Value_sym
	}
	//fmt.Printf("func_name %v\r\n", func_name)
	//fmt.Printf("al.Type %v\r\n", al.Type)
	switch al.Type {
	case Cell_array:
		//        fmt.Printf("al.Value_array %#v %v\r\n", al.Value_array, len(al.Value_array))
		for i, _ := range al.Value_array {
			vi := al.Value_array[i]
			//		fmt.Printf("i %v vi %v\r\n", i, vi)
			db := Cell2Short(vi)
			data_b = append(data_b, *db)
		}
	}
	if func_name != "" {
		//    fmt.Printf("data_b %#v\r\n", data_b)
		res, cs_out, err := g_el.CallExtFunc(func_name, data_b)
		if err != nil {
			fmt.Printf("CallExtFunc %v\r\n", err)
			return err
		}
		if res == 0 {
			if cs_out != nil {
				result = Short2Cell(cs_out)
			} else {
				result = Nil
			}
		} else {
			return errors.New("error call func")
		}
	}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_Iterate(fc *Func_call, args []*Cell) error {
	// проверяем, что аргумент 1
	if len(args) < 1 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 1 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	var result *Cell
	it := args[0]
	switch it.Type {
	case Cell_cell:
		v := it.Value_head
		cnt := it.Value_last
		var pos int64
		switch v.Type {
		case Cell_int:
			switch cnt.Type {
			case Cell_int:
				pos = cnt.Value_int
				//val := v.Value_int
				pos = pos << 1
				if pos == 0 {
					result = Nil
				} else {
					cnt.Value_int = pos
					result = it
				}
			default:
				result = Nil
			}
		case Cell_str:
			switch cnt.Type {
			case Cell_int:
				pos = cnt.Value_int
				val := v.Value_str
				pos = pos + 1
				if pos == int64(len(val)) {
					result = Nil
				} else {
					cnt.Value_int = pos
					result = it
				}
			default:
				result = Nil
			}
		case Cell_array:
			switch cnt.Type {
			case Cell_int:
				pos = cnt.Value_int
				val := v.Value_array
				pos = pos + 1
				if pos == int64(len(val)) {
					result = Nil
				} else {
					cnt.Value_int = pos
					result = it
				}
			default:
				result = Nil
			}
		case Cell_dict:
			switch cnt.Type {
			case Cell_int:
				pos = cnt.Value_int
				val := v.Value_dict
				pos = pos + 1
				// keys := reflect.ValueOf(val).MapKeys()
				if pos == int64(len(val)) {
					result = Nil
				} else {
					cnt.Value_int = pos
					result = it
				}
			default:
				result = Nil
			}
		case Cell_cell:
			//ss := cnt.String(false)
			//fmt.Printf("ss %v\r\n", ss)
			switch cnt.Type {
			case Cell_sym:
				fmt.Printf("Cell_sym\r\n")
				if cnt == Nil {
					result = Nil
				}
			case Cell_cell:
				if cnt == Nil {
					result = Nil
				} else {
					it.Value_last = cnt.Value_last
					if it.Value_last == Nil {
						result = Nil
					} else {
						result = it
					}
				}
			default:
				result = Nil
			}
		case Cell_sym:
			fmt.Printf("Cell_sym\r\n")
			if cnt == Nil {
				result = Nil
			}
		default:
			result = Nil
		}
	}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_GetIter(fc *Func_call, args []*Cell) error {
	// проверяем, что аргумент 1
	if len(args) < 1 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 1 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	var result *Cell
	it := args[0]
	switch it.Type {
	case Cell_cell:
		v := it.Value_head
		cnt := it.Value_last
		var pos int64
		switch v.Type {
		case Cell_int:
			switch cnt.Type {
			case Cell_int:
				pos = int64(cnt.Value_int)
				val := v.Value_int
				val = val & pos
				result = &Cell{Type: Cell_int, Value_int: int64(val)}
			default:
				result = Nil
			}
		case Cell_str:
			switch cnt.Type {
			case Cell_int:
				pos = cnt.Value_int
				val := v.Value_str
				val = val[pos : pos+1]
				result = &Cell{Type: Cell_int, Value_int: 0}
			default:
				result = Nil
			}
		case Cell_array:
			switch cnt.Type {
			case Cell_int:
				pos = cnt.Value_int
				val := v.Value_array
				result = val[pos]
			default:
				result = Nil
			}
		case Cell_dict:
			switch cnt.Type {
			case Cell_int:
				pos = cnt.Value_int
				val := v.Value_array
				// fmt.Printf("keys %v pos %v\r\n", keys, pos)
				result = val[pos]
				// &Cell{Type: Cell_str, Value_str: keys[pos].Interface().(string)}
			default:
				result = Nil
			}
		case Cell_cell:
			switch cnt.Type {
			case Cell_cell:
				result = cnt.Value_head
			default:
				result = Nil
			}
		default:
			result = Nil
		}

	}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_MakeIter(fc *Func_call, args []*Cell) error {
	// проверяем, что аргумент 1
	if len(args) < 1 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 1 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// по очереди вытаскиваем аргументы
	var result *Cell
	v := args[0]
	var pos int64
	switch v.Type {
	case Cell_int:
		pos = 1
		l := &Cell{Type: Cell_int, Value_int: pos}
		result = &Cell{Type: Cell_cell, Value_head: v, Value_last: l}
	case Cell_str:
		pos = 0
		l := &Cell{Type: Cell_int, Value_int: pos}
		result = &Cell{Type: Cell_cell, Value_head: v, Value_last: l}
	case Cell_array:
		pos = 0
		l := &Cell{Type: Cell_int, Value_int: pos}
		result = &Cell{Type: Cell_cell, Value_head: v, Value_last: l}
	case Cell_dict:
		pos = 0
		l := &Cell{Type: Cell_int, Value_int: pos}
		keys := reflect.ValueOf(v.Value_dict).MapKeys()
		vv := &Cell{Type: Cell_array, Value_array: []*Cell{}}
		for i := 0; i < len(keys); i++ {
			s := keys[i].Interface().(string)
			vs := &Cell{Type: Cell_str, Value_str: s}
			vv.Value_array = append(vv.Value_array, vs)
		}
		result = &Cell{Type: Cell_cell, Value_head: vv, Value_last: l}
	case Cell_cell:
		result = &Cell{Type: Cell_cell, Value_head: v, Value_last: v}
	default:
		result = Nil
	}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

func IF_FindFunc(fc *Func_call, args []*Cell) error {
	// проверяем, что аргумент 1
	if len(args) < 1 {
		var ErrTypeError = errors.New("too few args")
		return ErrTypeError
	}
	if len(args) > 1 {
		var ErrTypeError = errors.New("too many args")
		return ErrTypeError
	}
	// тип аргумента строка или символ
	var result *Cell
	v := args[0]
	result = Nil
	func_name := ""
	switch v.Type {
	case Cell_sym:
                func_name = v.Value_sym
	case Cell_str:
                func_name = v.Value_str
	}
	if func_name != "" {
		_, ok := fc.Env.Func_dict[func_name]
		if !ok {
			// fmt.Printf("Func %v not found\r\n", func_name)
		} else {
                        result = True
		}
	}
	fc.Stack.Stack = append(fc.Stack.Stack, result)

	return nil
}

