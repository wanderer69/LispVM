package common

import (
	//	"errors"
	"fmt"
	//	. "../Common"
)

// указатель дл¤ внешней функции
type Func_ext func(fc *Func_call) (*Cell, error)

type Func struct {
	Name         *Cell       // имя если оно имеет смысл
	Type         int         // тип 0 встроенная 1 внешняя
	Func_ext_ptr *Func_ext   // указатель на встроенную функцию
	Source       []*Cell     // исходник тела внешней функции
	Args         *Cell       // список аргументов
	Code         []Func_code // компилированный список вызовов
	Var_list_    []*Cell     // значений глобальных переменных
	Var_list     []string    // значений глобальных переменных
	Const_list   []*Cell     // значений глобальных констант
}

func (fn *Func) String() string {
	tp := ""
	ss := fmt.Sprintf("func %v", fn.Name.String(false))
	if fn.Type == 0 {
		tp = "int"
		ss = ss + fmt.Sprintf("<%v>", fn.Func_ext_ptr)
	} else if fn.Type == 1 {
		tp = "ext"
		ss = ss + fmt.Sprintf("<%v>", tp)
	}
	ss = ss + fmt.Sprintf(" %v var [", fn.Args.String(false))
	for i, v := range fn.Var_list {
		ss = ss + fmt.Sprintf(" %v %v", i, v)
	}
	ss = ss + fmt.Sprintf("]")

	ss = ss + fmt.Sprintf("const [")
	for i, v := range fn.Const_list {
		ss = ss + fmt.Sprintf(" %v %v", i, v.String(false))
	}
	ss = ss + fmt.Sprintf("] code {")
	for j, c := range fn.Code {
		ss = ss + fmt.Sprintf(" %v %v", j, PrintCode(c))
	}

	//fmt.Printf("func %v %v\r\n", fn.Name, fn.Type)
	return ss
}
