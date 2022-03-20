package common

//	"errors"
//	"fmt"

//	. "../Common"

// вызванна¤ функция - по сути это замыкание или генератор существующий в отдельной реальности
type Func_call struct {
	Env              *EnvironmentLVM     // окружение
	Func_call_prev   *Func_call       // односторонняя связь с данными предка
	Global_func_dict map[string]*Func // словарь глобальных функций
	Stack            Stack_data       // стек данных
	Func_current     *Func            // текущая функция
	PC               int              // текущий указатель на команду
	Var_dict         map[string]*Cell // локальный словарь значений переменных
	Const_dict       map[string]*Cell // локальный словарь значений констант
	Var_value        []*Cell          // локальный массив значений переменных
	// Const_value      []*Cell       // локальный массив значений констант
	//Debug            int              // признак отладки
}

type Func_code struct {
	Code int64 // код некомпилированной команды
	//	CmdName string                    // имя функции
	Cmd   func(fc *Func_call) error // указатель на вызов функции
	Const []int64                   // массив констант
}

func Func_code2Func_code_store(fc Func_code) Func_code_store {
	fcs := Func_code_store{}
	fcs.Code = int32(fc.Code)
	//	fcs.Cmd = fc.CmdName
	for i, _ := range fc.Const {
		fcs.Const = append(fcs.Const, int64(fc.Const[i]))
	}
	return fcs
}

func Func_code_store2Func_code(fcs Func_code_store) Func_code {
	fc := Func_code{}
	fc.Code = int64(fcs.Code)
	//	fc.Cmd = int(fcs.CmdName)
	for i, _ := range fcs.Const {
		fc.Const = append(fc.Const, int64(fcs.Const[i]))
	}
	return fc
}
