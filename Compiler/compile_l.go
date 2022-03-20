package compiler

import (
	"fmt"

	. "arkhangelskiy-dv.ru/LispVM/Common"
	//	"strconv"
	//	"unicode/utf8"
)

func CreateOp() ([]OperatorItem, map[string]OperatorItem) {
	oid := make(map[string]OperatorItem)
	oia := []OperatorItem{}
	oi := OperatorItem{"setq", OP_setq, 2}
	oia = append(oia, oi)
	oid["setq"] = oi
	oi = OperatorItem{"progn", OP_progn, 1}
	oia = append(oia, oi)
	oid["progn"] = oi
	oi = OperatorItem{"if", OP_if, 3}
	oia = append(oia, oi)
	oid["if"] = oi
	oi = OperatorItem{"lambda", OP_lambda, 2}
	oia = append(oia, oi)
	oid["lambda"] = oi
	oi = OperatorItem{"let", OP_let, 2}
	oia = append(oia, oi)
	oid["let"] = oi
	oi = OperatorItem{"defun", OP_defun, 3}
	oia = append(oia, oi)
	oid["defun"] = oi
	oi = OperatorItem{"apply", OP_apply, 1}
	oia = append(oia, oi)
	oid["apply"] = oi
	oi = OperatorItem{"and", OP_and, 2}
	oia = append(oia, oi)
	oid["and"] = oi
	oi = OperatorItem{"or", OP_or, 2}
	oia = append(oia, oi)
	oid["or"] = oi
	oi = OperatorItem{"quote", OP_quote, 1}
	oia = append(oia, oi)
	oid["quote"] = oi
	oi = OperatorItem{"loop", OP_loop, 1}
	oia = append(oia, oi)
	oid["loop"] = oi
	oi = OperatorItem{"go", OP_go, 2}
	oia = append(oia, oi)
	oid["go"] = oi
	oi = OperatorItem{"->", OP_to_channel, 2}
	oia = append(oia, oi)
	oid["<-"] = oi
	oi = OperatorItem{"->", OP_from_channel, 1}
	oia = append(oia, oi)
	oid["<-"] = oi

	oi = OperatorItem{"list", OP_list, 2}
	oia = append(oia, oi)
	oid["list"] = oi

	oi = OperatorItem{"values", OP_values, 2}
	oia = append(oia, oi)
	oid["values"] = oi

	oi = OperatorItem{"bind", OP_bind, 2}
	oia = append(oia, oi)
	oid["bind"] = oi

	return oia, oid
}

var G_Operators_list []OperatorItem
var G_Operators_dict map[string]OperatorItem
var G_Funcs_list []InternalFuncItem
var G_Funcs_dict map[string]InternalFuncItem

func InitFuncOp(dir_lst []string, ext_func bool) bool {
	l, d := CreateOp()
	G_Operators_list = l
	G_Operators_dict = d
	l1, d1, _ := CreateIntFunc(dir_lst, ext_func)
	G_Funcs_list = l1
	G_Funcs_dict = d1
	//	G_FuncsInt_dict = d2
	return true
}

func CheckInternal(s string) (bool, int) {
	//fmt.Printf("'%v'\r\n", s)
	if _, ok := G_Operators_dict[s]; ok {
		return true, 1
	} else {
		if _, ok := G_Funcs_dict[s]; ok {
			return true, 2
		}
	}
	return false, -1
}

func GetOperator(s string) *OperatorItem {
	if o, ok := G_Operators_dict[s]; ok {
		return &o
	}
	return nil
}

func GetInternal(s string) *InternalFuncItem {
	if ifi, ok := G_Funcs_dict[s]; ok {
		return &ifi
	}
	return nil
}

func InitCompilerEnv() *CompilerEnv {
	var ce CompilerEnv
	ce.VarDict = make(map[string]int64)
	return &ce
}

func CheckVar(n string, ce *CompilerEnv) int64 {
	// проверяем, что переменная есть если ее нет - возвращаем -1
	if p, ok := ce.VarDict[n]; ok {
		return p
	}
	return -1
}

func AddVar(n string, ce *CompilerEnv) int64 {
	// проверяем, что переменная есть если ее нет добавляем
	l := int64(len(ce.VarList))
	ce.VarList = append(ce.VarList, n)
	ce.VarDict[n] = l
	return l
}

func AddConst(n *Cell, ce *CompilerEnv) int64 {
	// проверяем, что переменная есть если ее нет добавляем
	//fmt.Printf("n %#v\r\n", n)
	//fmt.Printf("ce.ConstList %v\r\n", ce.ConstList)
	l := int64(len(ce.ConstList))
	ce.ConstList = append(ce.ConstList, n)
	return l
}

func Compile(c *Cell, ce *CompilerEnv, debug bool) ([]Func_code, bool) {
	fca := []Func_code{}
	// result := ""
	if debug {
		fmt.Printf("c.Type %v\r\n", c.Type)
	}
	switch c.Type {
	case Cell_sym:
		// пытаемся вычислить ЗНАЧЕНИЕ этого символа и загрузить в стек
		if (c == Nil) || (c == True) || (c == False) ||
			(c == T_int) || (c == T_float) ||
			(c == T_sym) || (c == T_str) || (c == T_cell) ||
			(c == T_dict) || (c == T_array) || (c == T_func) ||
			(c == T_ext) || (c == T_channel) || (c == T_next) {
			l := AddConst(c, ce)
			fc := Func_code{VMC_const, nil, []int64{l}}
			fca = append(fca, fc)
		} else {
			n := CheckVar(c.Value_sym, ce)
			if n == -1 {
				// надо ли добавлять? Если не было до Этого объявления
				n = AddVar(c.Value_sym, ce)
			}
			fc := Func_code{VMC_push, nil, []int64{n}}
			//fmt.Printf("Cell_sym %#v\r\n", fc)
			fca = append(fca, fc)
		}
	case Cell_int:
		// загружаем константу в стек
		// fmt.Printf("c %#v\r\n", c)
		l := AddConst(c, ce)
		// fmt.Printf("l %v\r\n", l)
		fc := Func_code{VMC_const, nil, []int64{l}}
		fca = append(fca, fc)
		// result = fmt.Sprintf("%v", c.Value_int)
	case Cell_str:
		// загружаем константу в стек
		l := AddConst(c, ce)
		fc := Func_code{VMC_const, nil, []int64{l}}
		//fmt.Printf("Cell_str %#v\r\n", fc)
		fca = append(fca, fc)
		// result = fmt.Sprintf("'%s'", c.Value_str)
	case Cell_float:
		// загружаем константу в стек
		// fmt.Printf("c %#v\r\n", c)
		l := AddConst(c, ce)
		// fmt.Printf("l %v\r\n", l)
		fc := Func_code{VMC_const, nil, []int64{l}}
		fca = append(fca, fc)
		// result = fmt.Sprintf("%v", c.Value_int)
	case Cell_cell:
		// вычисляем вызов функции (или замыкания)
		// надо проверить, в каком уровне
		// h := ""
		tl := 0
		// есть встроенные функции, которые имеют известное поведение
		if c.Value_head != nil {
			// проверяем, что там не список
			th := c.Value_head.Type
			if th == Cell_cell {
				// пока не знаю что с этим делать
			} else {
				if debug {
					fmt.Printf("th %v\r\n", th)
				}
				// проверяем что функция встроенная
				if th == Cell_sym {
					ok, r := CheckInternal(c.Value_head.Value_sym)
					if ok {
						// компилируем функцию
						if debug {
							fmt.Printf("r %v\r\n", r)
						}
						if r == 1 {
							op := GetOperator(c.Value_head.Value_sym)
							if debug {
								fmt.Printf("op %#v\r\n", op)
							}
							switch op.Code {
							case OP_setq:
								// разбираем аргументы первый аргумент должен быть символом и его не вычисляем. второй аргумент вычистяется
								f := c.Value_last
								if f.Value_head.Type != Cell_sym {
									// ошибка!!! заканчиваем трансляцию
									fmt.Printf("Error! Type not symbol\r\n")
									return fca, false
								}

								n := f.Value_last
								if n.Type != Cell_cell {
									// ошибка !!!
									fmt.Printf("Error! type not cell\r\n")
									return fca, false
								}
								// вычисляем аргумент
								fca_, res := Compile(n.Value_head, ce, debug)
								if !res {
									// найдена ошибка
									fmt.Printf("Error! Compile error\r\n")
									return fca, false
								}
								fca = append(fca, fca_...)
								// в стеке результат вычисления связываем его
								var_ptr := CheckVar(f.Value_head.Value_sym, ce)
								if var_ptr >= 0 {
									// переменная локальная есть
									fc := Func_code{VMC_pop, nil, []int64{var_ptr}}
									fca = append(fca, fc)
								} else {
									// новая локальная переменная!!! Добавляем ее!!
									l := AddVar(f.Value_head.Value_sym, ce)
									fc := Func_code{VMC_pop, nil, []int64{l}}
									fca = append(fca, fc)
								}
							case OP_progn:
								// вычисляем тело как последовательность списков
								if c.Value_last != nil {
									cc := c.Value_last
									tl = cc.Value_last.Type
									if tl == Cell_cell {
										for {
											// цикл по уровню пока не Nil хвост
											if cc.Value_head != nil {
												// вычисляем аргумент
												fca_, res := Compile(cc.Value_head, ce, debug)
												if !res {
													// найдена ошибка
													fmt.Printf("Error! Compile error\r\n")
													return fca, false
												}
												// fmt.Printf("> fca_ %#v\r\n", fca_)
												fca = append(fca, fca_...)
											}
											if cc.Value_last != nil {
												tl = cc.Value_last.Type
												if tl == Cell_cell {
													cc = cc.Value_last
												} else {
													if cc.Value_last == Nil {
														break
													} else {
														// странная ситуация
														break
													}
												}
											}
										}
									} else {
										// это непонятно что
										fmt.Printf("Error! Not understand expression\r\n")
										return fca, false
									}
								}
							case OP_if:
								// первый аргумент - условие второй - исполняемое тело если условие истинно третий испольняемое тело если условие ложно
								fca_t := []Func_code{}
								fca_n := []Func_code{}
								//ce_n := InitCompilerEnv()
								// вычисляем первый аргумент
								if debug {
									fmt.Printf("собираем условие\r\n")
								}
								eq := c.Value_last
								if eq.Value_head != nil {
									// вычисляем аргумент
									fca_, res := Compile(eq.Value_head, ce, debug)
									if !res {
										// найдена ошибка
										fmt.Printf("error compile eq.Value_head res %v\r\n", res)
										return fca, false
									}
									fca = append(fca, fca_...)
								}

								ev_t := eq.Value_last
								if debug {
									fmt.Printf("собираем если истинно\r\n")
								}
								s := ev_t.Value_head
								if s.Type != Cell_cell {
									// ошибка !!!
									fmt.Printf("error s.Type %v Cell_cell %v\r\n", s.Type, Cell_cell)
									return fca, false
								}
								cc := s
								if debug {
									fmt.Printf("cc true %v\r\n", cc.String(false))
								}
								// вычисляем аргумент
								fca_, res := Compile(cc, ce, debug)
								if !res {
									// найдена ошибка
									fmt.Printf("error compile cc (1) res %v\r\n", res)
									return fca, false
								}
								fca_t = append(fca_t, fca_...)
								ev_n := ev_t.Value_last
								if debug {
									fmt.Printf("собираем если ложно\r\n")
								}
								// проверяем что ветка иначе есть
								if ev_n != Nil {
									s = ev_n.Value_head
									if s.Type != Cell_cell {
										// ошибка !!!
										fmt.Printf("error s.Type %v Cell_cell %v\r\n", s.Type, Cell_cell)
										return fca, false
									}
									cc = s
									if debug {
										fmt.Printf("cc false %v\r\n", cc.String(false))
									}
									// вычисляем аргумент
									fca__, res_ := Compile(cc, ce, debug)
									if !res_ {
										// найдена ошибка
										fmt.Printf("error compile cc(2) res %v\r\n", res_)
										return fca, false
									}
									fca_n = append(fca_n, fca__...)
								}
								// собираем команду
								if debug {
									fmt.Printf("собираем команду\r\n")
								}
								// берем из топа и переход если не истина
								l := int64(len(fca_t) + 1)
								if debug {
									fmt.Printf("VMC_branch_false l = %v\r\n", l)
								}
								fc := Func_code{VMC_branch_false, nil, []int64{l}}
								fca = append(fca, fc)
								fca = append(fca, fca_t...)
								l = int64(len(fca_n))
								if debug {
									fmt.Printf("VMC_branch l = %v\r\n", l)
								}
								fc = Func_code{VMC_branch, nil, []int64{l}}
								fca = append(fca, fc)
								fca = append(fca, fca_n...)
							case OP_lambda:
								// функция без имени, которая просто грузится в стек
								// первый аргумент - список аргументов, второй - список выражений по сути перечень дя трансляции
								// отправляем в стек текущую структуру
								ce_n := InitCompilerEnv()
								fca_n := []Func_code{}
								n := c.Value_last
								if n.Type != Cell_cell {
									// ошибка !!!
									fmt.Printf("Error! n.Type %v Cell_cell %v\r\n", n.Type, Cell_cell)
									return fca, false
								}
								// проходим по списку переменных. добавляем их в список переменных
								cc := n.Value_head
								for {
									// цикл по уровню пока не Nil хвост
									if cc.Value_head != nil {
										// проверяем тип - должен быть символ
										if cc.Value_head.Type != Cell_sym {
											// ошибка!!! заканчиваем трансляцию
											fmt.Printf("Error! cc.Value_head.Type %v Cell_sym %v\r\n", cc.Value_head.Type, Cell_sym)
											return fca, false
										}
										if debug {
											fmt.Printf("cc.Value_head.Value_sym %v, ce_n %v\r\n", cc.Value_head.Value_sym, ce_n)
										}
										var_ptr := AddVar(cc.Value_head.Value_sym, ce_n)
										// добавляем код который возьмет значение из стека и добавит в переменную аргумента
										fc := Func_code{VMC_pop, nil, []int64{var_ptr}}
										fca_n = append(fca_n, fc)
									}
									if cc.Value_last != nil {
										tl = cc.Value_last.Type
										if tl == Cell_cell {
											cc = cc.Value_last
										} else {
											if cc.Value_last == Nil {
												break
											} else {
												// странная ситуация
												break
											}
										}
									}
								}
								s := n.Value_last
								if s.Type != Cell_cell {
									// ошибка !!!
									fmt.Printf("Error! s.Type %v Cell_cell %v\r\n", s.Type, Cell_cell)
									return fca, false
								}
								cc = s
								if debug {
									fmt.Printf("cc %v\r\n", cc.String(false))
								}
								for {
									// цикл по уровню пока не Nil хвост
									if cc.Value_head != nil {
										// вычисляем аргумент
										fca_, res := Compile(cc.Value_head, ce_n, debug)
										if !res {
											// найдена ошибка
											fmt.Printf("Error! res %v\r\n", res)
											return fca, false
										}
										fca_n = append(fca_n, fca_...)
									}
									if cc.Value_last != nil {
										tl = cc.Value_last.Type
										if tl == Cell_cell {
											cc = cc.Value_last
										} else {
											if cc.Value_last == Nil {
												break
											} else {
												// странная ситуация
												break
											}
										}
									}
								}
								if debug {
									fmt.Printf("ce_n %#v\r\n", ce_n)
								}
								// создаем функцию
								f_f := Func{}
								// генерируем имя
								f_f.Name = CellCreateFromString(GetMID()) //f.Value_head
								f_f.Type = 1
								// надо добавить выгрузку из стека в значения аргументов !!!!
								f_f.Code = fca_n
								f_f.Args = n.Value_head
								f_f.Var_list = ce_n.VarList
								f_f.Const_list = ce_n.ConstList
								ce.FuncList = append(ce.FuncList, &f_f)
								if debug {
									fmt.Printf("-> %v\r\n", f_f)
								}
								res_ := &Cell{Type: Cell_func, Value_func: &f_f}
								l := AddConst(res_, ce)
								fc := Func_code{VMC_const, nil, []int64{l}}
								fca = append(fca, fc)
							case OP_loop:
								// первый аргумент - условие, второй аргумент - тело цикла
								fca_c := []Func_code{}
								fca_l := []Func_code{}
								// вычисляем первый аргумент
								if debug {
									fmt.Printf("собираем условие\r\n")
								}
								eq := c.Value_last
								if eq.Value_head != nil {
									// вычисляем аргумент
									fca_, res := Compile(eq.Value_head, ce, debug)
									if !res {
										// найдена ошибка
										fmt.Printf("error compile eq.Value_head res %v\r\n", res)
										return fca, false
									}
									fca_c = append(fca_c, fca_...)
								}
								/*
									ev_l := eq.Value_last
									if debug {
										fmt.Printf("собираем тело\r\n")
									}
								*/
								s := eq.Value_last
								if s.Type != Cell_cell {
									// ошибка !!!
									fmt.Printf("error s.Type %v Cell_cell %v\r\n", s.Type, Cell_cell)
									return fca, false
								}
								cc := s
								// вычисляем тело
								for {
									if debug {
										fmt.Printf("cc true %v\r\n", cc.String(false))
									}
									// цикл по уровню пока не Nil хвост
									if cc.Value_head != nil {
										// fmt.Printf("cc true %v\r\n", cc.Value_head.String(false))
										// вычисляем аргумент
										fca_, res := Compile(cc.Value_head, ce, debug)
										if !res {
											// найдена ошибка
											fmt.Printf("Error! res %v\r\n", res)
											return fca, false
										}
										fca_l = append(fca_l, fca_...)
									}
									if cc.Value_last != nil {
										tl = cc.Value_last.Type
										if tl == Cell_cell {
											cc = cc.Value_last
										} else {
											if cc.Value_last == Nil {
												break
											} else {
												// странная ситуация
												break
											}
										}
									}
								}
								lb := int64(len(fca_c))
								//fmt.Printf("lb %v\r\n", lb)

								fca = append(fca, fca_c...)
								// берем из топа и переход если не истина
								l := int64(len(fca_l) + 1)
								if debug {
									fmt.Printf("VMC_branch_false l = %v\r\n", l)
								}
								//fmt.Printf("l %v\r\n", l)
								fc := Func_code{VMC_branch_true, nil, []int64{l}}
								fca = append(fca, fc)
								fca = append(fca, fca_l...)
								ll := int64(l + 1 + lb)
								//fmt.Printf("ll %v\r\n", ll)
								fc = Func_code{VMC_branch, nil, []int64{-ll}}
								fca = append(fca, fc)
								fc = Func_code{VMC_nop, nil, []int64{}}
								fca = append(fca, fc)
							case OP_let:
							case OP_defun:
								// первый аргумент - наименование, второй - список аргументов, третий - список выражений по сути перечень дя трансляции
								fca_n := []Func_code{}
								f := c.Value_last
								flag := false
								ss := f.Value_last
								if f.Value_head.Type != Cell_sym {
									if f.Value_head.Type == Cell_cell {
										// a la scheme
										f = f.Value_head
										if f.Value_head.Type != Cell_sym {
											fmt.Printf("Error! f.Value_head.Type %v Cell_sym %v\r\n", f.Value_head.Type, Cell_sym)
											return fca, false
										}
										flag = true
									} else {
										// ошибка!!! заканчиваем трансляцию
										fmt.Printf("Error! f.Value_head.Type %v Cell_sym %v\r\n", f.Value_head.Type, Cell_sym)
										return fca, false
									}
								}
								// fmt.Printf("> f %v\r\n", f.String(false))
								// отправляем в стек текущую структуру
								ce_n := InitCompilerEnv()
								n := f.Value_last
								if n.Type != Cell_cell {
									// ошибка !!!
									fmt.Printf("Error! n.Type %v Cell_cell %v\r\n", n.Type, Cell_cell)
									return fca, false
								}
								cc := n.Value_head
								if flag {
									cc = n
								}
								// проходим по списку переменных. добавляем их в список переменных
								for {
									// цикл по уровню пока не Nil хвост
									if cc.Value_head != nil {
										// проверяем тип - должен быть символ
										if cc.Value_head.Type != Cell_sym {
											// ошибка!!! заканчиваем трансляцию
											fmt.Printf("Error! cc.Value_head.Type %v Cell_sym %v\r\n", cc.Value_head.Type, Cell_sym)
											return fca, false
										}
										if debug {
											fmt.Printf("cc.Value_head.Value_sym %v, ce_n %v\r\n", cc.Value_head.Value_sym, ce_n)
										}
										var_ptr := AddVar(cc.Value_head.Value_sym, ce_n)
										// добавляем код который возьмет значение из стека и добавит в переменную аргумента
										fc := Func_code{VMC_pop, nil, []int64{var_ptr}}
										fca_n = append(fca_n, fc)
									}
									if cc.Value_last != nil {
										tl = cc.Value_last.Type
										if tl == Cell_cell {
											cc = cc.Value_last
										} else {
											if cc.Value_last == Nil {
												break
											} else {
												// странная ситуация
												break
											}
										}
									}
								}
								s := n.Value_last
								// fmt.Printf("> ss %v\r\n", ss.String(false))
								if flag {
									s = ss
								}
								if s.Type != Cell_cell {
									// ошибка !!!
									fmt.Printf("Error! s.Type %v Cell_cell %v\r\n", s.Type, Cell_cell)
									return fca, false
								}
								cc = s
								if debug {
									fmt.Printf("cc %v\r\n", cc.String(false))
								}
								for {
									// цикл по уровню пока не Nil хвост
									if cc.Value_head != nil {
										// вычисляем аргумент
										fca_, res := Compile(cc.Value_head, ce_n, debug)
										if !res {
											// найдена ошибка
											fmt.Printf("Error! res %v\r\n", res)
											return fca, false
										}
										fca_n = append(fca_n, fca_...)
									}
									if cc.Value_last != nil {
										tl = cc.Value_last.Type
										if tl == Cell_cell {
											cc = cc.Value_last
										} else {
											if cc.Value_last == Nil {
												break
											} else {
												// странная ситуация
												break
											}
										}
									}
								}
								if debug {
									fmt.Printf("ce_n %#v\r\n", ce_n)
								}
								// создаем функцию
								f_f := Func{}
								f_f.Name = f.Value_head
								f_f.Type = 1
								// надо добавить выгрузку из стека в значения аргументов !!!!
								f_f.Code = fca_n
								f_f.Args = n.Value_head
								f_f.Var_list = ce_n.VarList
								f_f.Const_list = ce_n.ConstList
								ce.FuncList = append(ce.FuncList, &f_f)
								if debug {
									fmt.Printf("-> %v\r\n", f_f)
								}
							case OP_apply:
							case OP_and:
							case OP_or:
							case OP_quote:
								// аргумент этой функции один и он не вычисляется его просто записываем в константу
								f := c.Value_last
								// собираем команду
								//fmt.Printf("-> %v\r\n", f.Value_head)
								l := AddConst(f.Value_head, ce)
								fc := Func_code{VMC_const, nil, []int64{l}}
								fca = append(fca, fc)
							case OP_list:
								f := c.Value_last
								for {
									if f.Value_head != nil {
										// вычисляем аргумент
										fca_, res := Compile(f.Value_head, ce, debug)
										if !res {
											// найдена ошибка
											fmt.Printf("Error! Compile error!\r\n")
											return fca, false
										}
										fca = append(fca, fca_...)
									}
									if f.Value_last == Nil {
										break
									} else {
										f = f.Value_last
									}
								}
							case OP_values:
								// разбираем и вычисляем аргументы
								// надо проверить, что есть признак возврата множественного значений. если нет - ошибка!
								cc := c.Value_last
								n := 0
								if cc.Type == Cell_cell {
									for {
										// цикл по уровню пока не Nil хвост
										if cc.Value_head != nil {
											// вычисляем аргумент
											fca_, res := Compile(cc.Value_head, ce, debug)
											if !res {
												// найдена ошибка
												fmt.Printf("Error! Compile error!\r\n")
												return fca, false
											}
											fca = append(fca, fca_...)
											n = n + 1
										}
										if cc.Value_last != nil {
											if cc.Value_last.Type == Cell_cell {
												cc = cc.Value_last
											} else {
												if cc.Value_last == Nil {
													break
												}
											}
										}
									}
									//
									c := &Cell{Type: Cell_multiple, Value_int: int64(n)}
									l := AddConst(c, ce)
									fc := Func_code{VMC_const, nil, []int64{l}}
									fca = append(fca, fc)
								} else {
									// это непонятно что
									fmt.Printf("Error! Not understand !\r\n")
									return fca, false
								}
							case OP_bind:
								// идет список переменных, а потом вызов функции
								cc := c.Value_last
								if cc.Type != Cell_cell {
									// найдена ошибка
									fmt.Printf("Error! Compile error!\r\n")
									return fca, false
								}
								// список переменных
								vl := cc.Value_head
								if vl.Type != Cell_cell {
									// найдена ошибка
									fmt.Printf("Error! Compile error!\r\n")
									return fca, false
								}
								cc = cc.Value_last
								if cc.Type != Cell_cell {
									// найдена ошибка
									fmt.Printf("Error! Compile error!\r\n")
									return fca, false
								}
								if cc.Value_head != nil {
									// вычисляем аргумент
									fca_, res := Compile(cc.Value_head, ce, debug)
									if !res {
										// найдена ошибка
										fmt.Printf("Error! Compile error!\r\n")
										return fca, false
									}
									fca = append(fca, fca_...)
								}
								for {
									// цикл по уровню пока не Nil хвост
									v := vl.Value_head
									if v != nil {
										if v.Type != Cell_sym {
											// ошибка!!! заканчиваем трансляцию
											fmt.Printf("Error! Type not symbol\r\n")
											return fca, false
										}
										// в стеке результат вычисления связываем его
										var_ptr := CheckVar(v.Value_sym, ce)
										if var_ptr >= 0 {
											// переменная локальная есть
											fc := Func_code{VMC_pop, nil, []int64{var_ptr}}
											fca = append(fca, fc)
										} else {
											// новая локальная переменная!!! Добавляем ее!!
											l := AddVar(v.Value_sym, ce)
											fc := Func_code{VMC_pop, nil, []int64{l}}
											fca = append(fca, fc)
										}

									}
									if vl.Value_last != nil {
										if vl.Value_last.Type == Cell_cell {
											vl = vl.Value_last
										} else {
											if vl.Value_last == Nil {
												break
											}
										}
									}
								}
							case OP_go:
								// первый аргумент = функция, второй - список аргументов
								cc := c.Value_last
								if cc.Value_head.Type != Cell_sym {
									// ошибка!!! заканчиваем трансляцию
									fmt.Printf("Error! Go arg error!\r\n")
									return fca, false
								}
								// fmt.Printf("cc %v\r\n", cc.String(false))
								f := cc.Value_head
								// fmt.Printf("f %v\r\n", f.String(false))

								cc = cc.Value_last
								cc = cc.Value_head
								// fmt.Printf("cc %v\r\n", cc.String(false))
								if cc.Type == Cell_cell {
									// формируем команду начала стека
									fc := Func_code{VMC_enter, nil, []int64{}}
									fca = append(fca, fc)
									// добавляем аргументы
									for {
										// цикл по уровню пока не Nil хвост
										if cc.Value_head != nil {
											// вычисляем аргумент
											fca_, res := Compile(cc.Value_head, ce, debug)
											if !res {
												// найдена ошибка
												fmt.Printf("Error! Compile error!\r\n")
												return fca, false
											}
											fca = append(fca, fca_...)
										}
										if cc.Value_last != nil {
											if cc.Value_last.Type == Cell_cell {
												cc = cc.Value_last
											} else {
												if cc.Value_last == Nil {
													break
												}
											}
										}
									}
								} else {
									// это непонятно что
									fmt.Printf("Error! Not understand !\r\n")
									return fca, false
								}
								// собираем команду вызова встроенной функции
								if debug {
								}
								l := AddConst(f, ce)
								if debug {
									fmt.Printf("l %v\r\n", l)
								}
								fc := Func_code{VMC_const, nil, []int64{l}}
								fca = append(fca, fc)
								fc = Func_code{VMC_g_call, nil, []int64{}}
								fca = append(fca, fc)
							}
						} else {
							if r == 2 {
								fn := GetInternal(c.Value_head.Value_sym)
								// смотрим, что список имеет необходимое число аргументов
								if c.Value_last == Nil {
									// ошибка!
									fmt.Printf("Error! List args too short!\r\n")
									return fca, false
								}
								if debug {
									fmt.Printf("fn %v\r\n", fn)
									fmt.Printf("G_Funcs_dict %v\r\n", G_Funcs_dict)
								}
								cc := c.Value_last
								//fmt.Printf("cc %#v\r\n", cc)
								n_arg := 0
								args_len := len(fn.Args)
								flag_r := -1
								for {
									// цикл по уровню пока не Nil хвост
									if cc.Value_head != nil {
										// увеличиваем счетчик аргументов
										n_arg = n_arg + 1
										// проверяем тип
										if args_len > n_arg {
											if fn.Args[n_arg-1] == "%r" {
												flag_r = n_arg - 1
												break
											}
										} else {
											// Ошибка!
										}
									}
									//fmt.Printf("cc.Value_head %v\r\n", cc.Value_head.String(false))
									tl = cc.Value_last.Type
									if cc.Value_last.Type == Cell_cell {
										cc = cc.Value_last
									} else {
										if cc.Value_last == Nil {
											break
										}
									}
								}
								// вычисляем аргументы и собираем их в стек
								cc = c.Value_last
								n_arg = 0
								if debug {
									fmt.Printf("Int func\r\n")
								}
								// формируем команду начала стека
								fc := Func_code{VMC_enter, nil, []int64{}}
								fca = append(fca, fc)
								// добавляем аргументы
								for {
									// цикл по уровню пока не Nil хвост
									if cc.Value_head != nil {
										// увеличиваем номер
										n_arg = n_arg + 1
										if flag_r < n_arg {
											// остальные аргументы надо после вычисления в отдельный список и в стек.
										}
										// вычисляем аргумент
										fca_, res := Compile(cc.Value_head, ce, debug)
										if !res {
											// найдена ошибка
											fmt.Printf("Error! Compile error!\r\n")
											return fca, false
										}
										if debug {
											fmt.Printf("f fca_ %#v\r\n", fca_)
										}
										fca = append(fca, fca_...)
									}
									if cc.Value_last != nil {
										if cc.Value_last.Type == Cell_cell {
											cc = cc.Value_last
										} else {
											if cc.Value_last == Nil {
												break
											}
										}
									}
								}
								// }
								// собираем команду вызова встроенной функции
								fc = Func_code{VMC_call, nil, []int64{fn.Cmd}}
								fca = append(fca, fc)
								if debug {
									fmt.Printf("Int func end\r\n")
								}
							}
						}
					} else {
						if debug {
							fmt.Printf("-> %v\r\n", ok)
						}
						// строим вызов внешней функции
						// вычисляем аргументы и собираем их в стек
						cc := c.Value_last
						if cc.Type == Cell_cell {
							// формируем команду начала стека
							fc := Func_code{VMC_enter, nil, []int64{}}
							fca = append(fca, fc)
							// добавляем аргументы
							for {
								// цикл по уровню пока не Nil хвост
								if cc.Value_head != nil {
									// вычисляем аргумент
									fca_, res := Compile(cc.Value_head, ce, debug)
									if !res {
										// найдена ошибка
										fmt.Printf("Error! Compile error!\r\n")
										return fca, false
									}
									fca = append(fca, fca_...)
								}
								if cc.Value_last != nil {
									if cc.Value_last.Type == Cell_cell {
										cc = cc.Value_last
									} else {
										if cc.Value_last == Nil {
											break
										}
									}
								}
							}
						} else {
							// это непонятно что

							return fca, false
						}
						// собираем команду вызова встроенной функции
						if debug {
							fmt.Printf("c.Value_head %#v\r\n", c.Value_head)
						}
						l := AddConst(c.Value_head, ce)
						if debug {
							fmt.Printf("l %v\r\n", l)
						}
						fc := Func_code{VMC_const, nil, []int64{l}}
						fca = append(fca, fc)
						fc = Func_code{VMC_e_call, nil, []int64{}}
						fca = append(fca, fc)
					}
				} else {
					// это ошибка
					fmt.Printf("Error! Compile error!\r\n")
					return fca, false
				}
			}
		}
	case Cell_dict:
		// загружаем константу в стек
		l := AddConst(c, ce)
		fc := Func_code{VMC_const, nil, []int64{l}}
		fca = append(fca, fc)
		/*
			for k, v := range c.Value_dict {
				if len(result) > 0 {
					result = result + " " + fmt.Sprintf("%s:%s", k, v)
					//				result = result + " " + fmt.Sprintf("%s:%s", k, v.Print(level, debug))
				} else {
					//				result = fmt.Sprintf("%s:%s", k, v.Print(level, debug))
				}
			}
		*/
	case Cell_array:
		// загружаем константу в стек
		l := AddConst(c, ce)
		fc := Func_code{VMC_const, nil, []int64{l}}
		fca = append(fca, fc)
		/*
			for _, v := range c.Value_array {
				if len(result) > 0 {
					result = result + " " + fmt.Sprintf("%s", v)
					//				result = result + " " + fmt.Sprintf("%s", v.Print(level, debug))
				} else {
					//				result = fmt.Sprintf("%s", v.Print(level, debug))
				}
			}
		*/
	case Cell_func:
		l := AddConst(c, ce)
		fc := Func_code{VMC_const, nil, []int64{l}}
		fca = append(fca, fc)
		// result = fmt.Sprintf("%s", c.Value_func)
	}
	return fca, true
}
