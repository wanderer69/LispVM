package common

import (
	"fmt"
	//	"strconv"
	//	"unicode/utf8"
)

var G_str_to_vmc map[string]int
var G_vmc_to_str map[int64]string

func init() {
	G_str_to_vmc = make(map[string]int)
	G_vmc_to_str = make(map[int64]string)
	G_str_to_vmc["const"] = VMC_const
	G_str_to_vmc["dup"] = VMC_dup
	G_str_to_vmc["push"] = VMC_push
	G_str_to_vmc["pop"] = VMC_pop
	G_str_to_vmc["enter"] = VMC_enter
	G_str_to_vmc["call"] = VMC_call
	G_str_to_vmc["e_call"] = VMC_e_call
	G_str_to_vmc["branch"] = VMC_branch
	G_str_to_vmc["branch_true"] = VMC_branch_true
	G_str_to_vmc["branch_false"] = VMC_branch_false

	G_vmc_to_str[VMC_const] = "const"
	G_vmc_to_str[VMC_dup] = "dup"
	G_vmc_to_str[VMC_push] = "push"
	G_vmc_to_str[VMC_pop] = "pop"
	G_vmc_to_str[VMC_enter] = "enter"
	G_vmc_to_str[VMC_call] = "call"
	G_vmc_to_str[VMC_e_call] = "e_call"
	G_vmc_to_str[VMC_branch] = "branch"
	G_vmc_to_str[VMC_branch_true] = "branch_true"
	G_vmc_to_str[VMC_branch_false] = "branch_false"
}

func PrintCode(fc Func_code) string {
	/*
		Code  int                                            // код некомпилированной команды
		Cmd   func(env *Environment, current_call int) error // указатель на вызов функции
		Const []int                                          // массив констант
	*/
	sc, ok := G_vmc_to_str[fc.Code]
	if !ok {

	}
	ss := fmt.Sprintf("%v [", sc)
	for _, a := range fc.Const {
		ss = ss + fmt.Sprintf(" %v", a)
	}
	ss = ss + fmt.Sprintf("]")
	return ss
}

/*
	Name         *Cell       // имя если оно имеет смысл
	Type         int         // тип 0 встроенная 1 внешняя
	Func_ext_ptr *Func_ext   // указатель на встроенную функцию
	Source       []*Cell     // исходник тела внешней функции
	Args         *Cell       // список аргументов
	Code         []Func_code // компилированный список вызовов
	Var_list_    []*Cell     // значений глобальных переменных
	Var_list     []string    // значений глобальных переменных
	Const_list   []*Cell     // значений глобальных констант
*/
