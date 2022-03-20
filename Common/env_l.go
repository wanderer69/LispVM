package common

import (
"sync"
)
type CallItem struct {
	Func_call_array []*Func_call     // массив исполняемых вызовов
	FuncStack       []*Int_array     // массив массивов стека вызовов ()
	CurrentExec     []int            // список вызовов
	Result          *[]*Cell         // указательна результат выполнения
	Wg              *sync.WaitGroup  // указатель для окончания 
}

// окружение вызова функции
type EnvironmentLVM struct {
        Current_proc int
        Proc_dict map[int]int            // словарь идентификаторов процессов
        Call_Item []CallItem             // Исполняемая часть - массив процессов
//	Func_call_array []*Func_call     // массив исполняемых вызовов
//	FuncStack       []*Int_array     // массив массивов стека вызовов ()
	Func_dict       map[string]*Func //словарь глобальных функций
	Var_dict        map[string]*Cell // словарь значений глобальных переменных
	Const_dict      map[string]*Cell // словарь значений глобальных констант
	Var_list        []*Cell          // словарь значений глобальных переменных
	Const_list      []*Cell          // словарь значений глобальных констант
//	CurrentExec     []int            // список вызовов
//	Result          *[]*Cell         // указательна результат выполнения
	NoPrint         bool
	ExtContext      interface{}      // Внешний контекст
	Debug           int              // признак отладки
	FileFormat      string           // признак формата файла
}

func InitEnvironment(no_print bool) *EnvironmentLVM {
	var e EnvironmentLVM
	e.Func_dict = make(map[string]*Func)
	e.Var_dict = make(map[string]*Cell)
	e.Const_dict = make(map[string]*Cell)
	e.NoPrint = no_print
        e.Proc_dict = make(map[int]int)
        // e.Call_Item = append(e.Call_Item, CallItem{})
        e.Current_proc = 0
	return &e
}

func (e *EnvironmentLVM) FuncToEnv(f *Func) {
	e.Func_dict[f.Name.String(false)] = f
}

func (e *Env) EvalString(str string) (string, int) {
	c, err1 := e.StringToCell(str)
	if err1 != 0 {
		return "", err1
	}
	r, err2 := e.EvalCell(c)
	if err2 != 0 {
		return "", err1
	}

	s, err3 := e.CellToString(r)
	if err3 != 0 {
		return "", err1
	}
	return s, 0
}

func (e *Env) StringToCell(str string) (*Cell, int) {
	return nil, 0
}

func (e *Env) EvalCell(c *Cell) (*Cell, int) {
	return nil, 0
}

func (e *Env) CellToString(c *Cell) (string, int) {
	return "", 0
}
