package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//	"strings"
	//	"unicode/utf8"
)

const Cell_sym = 1
const Cell_int = 2
const Cell_float = 3
const Cell_str = 4
const Cell_cell = 5
const Cell_dict = 6
const Cell_array = 7
const Cell_func = 8
const Cell_object = 9
const Cell_channel = 10
const Cell_dict_n = 11
const Cell_multiple = 12

const Cell_ext = 100      // признак расширенного типа
const Cell_next = 1000000 // признак расширенного типа

const FC_car = 1
const FC_cdr = 2
const FC_cons = 3
const FC_listp = 4
const FC_eq = 5
const FC_rplaca = 6
const FC_rplacd = 7
const FC_list = 8
const FC_print = 9
const FC_add = 10
const FC_sub = 11
const FC_mul = 12
const FC_div = 13
const FC_type = 14
const FC_2int = 15
const FC_2str = 16
const FC_2float = 17
const FC_typep = 18
const FC_assert = 19
const FC_format = 20
const FC_make_array = 21
const FC_make_dict = 22
const FC_array = 23
const FC_dict = 24
const FC_item = 25
const FC_slice = 26
const FC_insert = 27
const FC_append = 28
const FC_import = 29
const FC_strip = 30
const FC_split = 31
const FC_index = 32
const FC_length = 33
const FC_join = 34
const FC_set_dict = 35
const FC_load_lib = 36
const FC_call_lib_func = 37
const FC_unload_lib = 38
const FC_princ = 39
const FC_iterate = 40
const FC_get_iter = 41
const FC_make_iter = 42
const FC_not = 43
const FC_find_func = 44

const VMC_nop = 0           // пустая команда
const VMC_const = 1         // загружает в стек по номеру
const VMC_dup = 2           // копирует в верхушке стека
const VMC_push = 3          // из переменной значение в стек переменную указывает номер аргумента
const VMC_pop = 4           // из вершины стека в переменную. переменную указывает номер аргумента
const VMC_enter = 5         // формирует код начала вызова в стеке
const VMC_call = 6          // вызывает функцию. имя функции - первый элемент в стеке. строит новый стек вызова функции
const VMC_e_call = 7        // вызывает функцию. имя функции - первый элемент в стеке. строит новый стек вызова функции
const VMC_branch = 8        // переход по относительному адресу из верхушки стека
const VMC_branch_true = 9   // переход по относительному адресу из верхушки стека если в вершине -1 стека истина
const VMC_branch_false = 10 // переход по относительному адресу из верхушки стека если в вершине -1 стека истина
const VMC_g_call = 11       // вызывает функцию в отдельном потоке gorutine. имя функции - первый элемент в стеке. строит новый стек вызова функции

const OP_setq = 1
const OP_progn = 2
const OP_if = 3
const OP_lambda = 4
const OP_let = 5
const OP_defun = 6
const OP_apply = 7
const OP_and = 8
const OP_or = 9
const OP_quote = 10
const OP_loop = 11
const OP_for_each = 12
const OP_go = 13
const OP_to_channel = 14
const OP_from_channel = 15
const OP_list = 16
const OP_values = 17
const OP_bind = 18

type InternalFuncItem struct {
	Name string
	Cmd  int64
	Args []string // список масок аргументов
}

type OperatorItem struct {
	Name string
	Code int
	Args int
}

type Int_array []int

// представление исполняемого модуля в виде файла
type Func_code_store struct {
	Code int32 `json:"Code,omitempty"` // код некомпилированной команды
	//	Cmd   string  `json:"Cmd,omitempty"`   // указатель на вызов функции
	Const []int64 `json:"Const,omitempty"` // массив констант
}

type Cell_store struct {
	Type          int32          `json:"type,omitempty"`
	Value_int     int64          `json:"int,omitempty"`
	Value_float   float64        `json:"float,omitempty"`
	Value_sym     string         `json:"sym,omitempty"`
	Value_str     string         `json:"str,omitempty"`
	Value_head    Cell_store_int `json:"head,omitempty"`
	Value_last    Cell_store_int `json:"last,omitempty"`
	Value_dict    Cell_store_arr `json:"dict,omitempty"`
	Value_array   Cell_store_int `json:"array,omitempty"`
	Value_func    Func_store_int `json:"func,omitempty"`
	Value_obj     string         `json:"object,omitempty"`
	Value_ext     string         `json:"ext,omitempty"`
	Value_channel string         `json:"channel,omitempty"`
}

type Ext_Type_store struct {
	ID    []byte   `json:"ID,omitempty"`    // уникальный идентификатор
	Name  string   `json:"name,omitempty"`  // имя если оно имеет смысл
	Value []string `json:"value,omitempty"` // значение в виде списка
}

type Func_store_int []Func_store
type Cell_store_int []Cell_store
type Cell_store_arr map[string]Cell_store

type Module_store struct {
	Name       string       `json:"name,omitempty"`   // имя если оно имеет смысл
	Var_list   []Cell_store `json:"vars,omitempty"`   // словарь значений глобальных переменных
	Const_list []Cell_store `json:"consts,omitempty"` // словарь значений глобальных констант
	Func_list  []Func_store `json:"funcs,omitempty"`
}

/*
func Load_func() {
	// читаем файл
	//
}
*/

func Cell2Cell_store(c *Cell) (Cell_store, int) {
	cs := Cell_store{}
	cs.Type = int32(c.Type)
	cs.Value_int = c.Value_int
	cs.Value_float = c.Value_float
	cs.Value_sym = c.Value_sym
	cs.Value_str = c.Value_str
	if c.Value_head != nil {
		vh, res := Cell2Cell_store(c.Value_head)
		if res != 0 {
			return cs, res << 1
		}
		cs.Value_head = Cell_store_int{vh}
	}
	if c.Value_last != nil {
		vl, res := Cell2Cell_store(c.Value_last)
		if res != 0 {
			return cs, res << 1
		}
		cs.Value_last = Cell_store_int{vl}
	}
	// d := make(map[string]Cell_store)
	cs.Value_dict = make(map[string]Cell_store)
	for k, v := range c.Value_dict {
		vd, res := Cell2Cell_store(v)
		if res != 0 {
			return cs, res << 1
		}
		cs.Value_dict[k] = vd
	}
	for i, _ := range c.Value_array {
		va, res := Cell2Cell_store(c.Value_array[i])
		if res != 0 {
			return cs, res << 1
		}
		cs.Value_array = append(cs.Value_array, va)
	}
	if c.Value_func != nil {
		fcs, res := Func2Func_store(c.Value_func)
		if res != 0 {
			return cs, res
		}
		cs.Value_func = Func_store_int{fcs}
	}
	if c.Value_ext != "" {
		//
		cs.Value_ext = c.Value_ext
	}
	return cs, 0
}

func Cell_store2Cell(cs Cell_store) *Cell {
	switch cs.Value_sym {
	case "Nil":
		return Nil
	case "True":
		return True
	case "False":
		return False
	case "nil":
		return Nil
	case "true":
		return True
	case "false":
		return False
	case "int":
		return T_int
	case "float":
		return T_float
	case "sym":
		return T_sym
	case "str":
		return T_str
	case "cell":
		return T_cell
	case "dict":
		return T_dict
	case "array":
		return T_array
	case "func":
		return T_func
	case "ext":
		return T_ext
	case "channel":
		return T_channel
	case "new_ext":
		return T_next
	default:
		c := Cell{}
		c.Type = int(cs.Type)
		c.Value_int = cs.Value_int
		c.Value_float = cs.Value_float
		// это символ
		c.Value_sym = cs.Value_sym
		c.Value_str = cs.Value_str
		if len(cs.Value_head) > 0 {
			vh := Cell_store2Cell(cs.Value_head[0])
			c.Value_head = vh
		}
		if len(cs.Value_last) > 0 {
			vl := Cell_store2Cell(cs.Value_last[0])
			c.Value_last = vl
		}
		c.Value_dict = make(map[string]*Cell)
		for k, v := range cs.Value_dict {
			vv := Cell_store2Cell(v)
			c.Value_dict[k] = vv
		}
		for i, _ := range cs.Value_array {
			va := Cell_store2Cell(cs.Value_array[i])
			c.Value_array = append(c.Value_array, va)
		}
		if len(cs.Value_func) > 0 {
			vf := Func_store2Func(cs.Value_func[0])
			c.Value_func = &vf
		}
		c.Value_ext = cs.Value_ext
		return &c
	}
}

type Func_store struct {
	Name       string            `json:"name"` // имя если оно имеет смысл
	Type       int32             `json:"type"` // тип 0 встроенная 1 внешняя
	Args       Cell_store        `json:"args"` // список аргументов
	Code       []Func_code_store `json:"const"`
	Var_list   []string          `json:"var_list"`
	Const_list []Cell_store      `json:"const_list"`
}

func Func2Func_store(fc *Func) (Func_store, int) {
	fcs := Func_store{}
	for i, _ := range fc.Code {
		fcs_i := Func_code2Func_code_store(fc.Code[i])
		fcs.Code = append(fcs.Code, fcs_i)
	}
	fcs.Name = fc.Name.Value_sym
	fcs.Type = int32(fc.Type)
	if fc.Args != nil {
		fca, res := Cell2Cell_store(fc.Args)
		if res != 0 {
			fmt.Print("Error Args %#v\r\n", fc.Args)
			return fcs, res
		}
		fcs.Args = fca
	}
	for i, _ := range fc.Const_list {
		if fc.Const_list[i] != nil {
			fcc, res := Cell2Cell_store(fc.Const_list[i])
			if res != 0 {
				fmt.Print("Error fc.Const_list[%v] %#v\r\n", i, fc.Const_list[i])
				return fcs, res
			}
			fcs.Const_list = append(fcs.Const_list, fcc)
		} else {
			fmt.Print("Error fc.Const_list[%v] %#v\r\n", i, fc.Const_list[i])
			return fcs, 1
		}
	}
	fcs.Var_list = fc.Var_list
	/*
		//	fcs.Cmd = fc.CmdName
		for i, _ := range fc.Const {
			fcs.Const = append(fcs.Const, int32(fc.Const[i]))
		}
	*/
	return fcs, 0
}

func Func_store2Func(fcs Func_store) Func {
	fc := Func{}
	for i, _ := range fcs.Code {
		fc_i := Func_code_store2Func_code(fcs.Code[i])
		fc.Code = append(fc.Code, fc_i)
	}
	fc.Name = &Cell{Type: Cell_sym, Value_sym: fcs.Name}
	fc.Type = int(fcs.Type)
	// if fcs.Args != nil {
	f := Cell_store2Cell(fcs.Args)
	fc.Args = f
	//}
	for i, _ := range fcs.Const_list {
		f := Cell_store2Cell(fcs.Const_list[i])
		fc.Const_list = append(fc.Const_list, f)
	}
	fc.Var_list = fcs.Var_list

	return fc
}

func CreateModule_store(name string) Module_store {
	ms := Module_store{}
	ms.Name = name
	return ms
}

func (ms *Module_store) AddFunc(fcs Func_store) {
	ms.Func_list = append(ms.Func_list, fcs)
}

func (ms *Module_store) AddConst(c Cell_store) {
	ms.Const_list = append(ms.Const_list, c)
}

func (ms *Module_store) AddVar(c Cell_store) {
	ms.Var_list = append(ms.Var_list, c)
}

func LoadModule(file_name string) (Module_store, error) {
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

func SaveModule(file_name string, cp Module_store) error {
	data, err := json.MarshalIndent(&cp, " ", "\t")
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	err1 := ioutil.WriteFile(file_name, data, 0644)
	return err1
}

type IntFuncDescr struct {
	FuncCode int64    `json:"func_code"`
	Name     string   `json:"name"` // имя если оно имеет смысл
	Args     []string `json:"args"` // список аргументов
}

type IntFuncFile struct {
	ID    []byte         `json:"id"`
	Name  string         `json:"name"`  // имя если оно имеет смысл
	Funcs []IntFuncDescr `json:"funcs"` // список аргументов
}

func LoadIntFuncFile(file_name string) (IntFuncFile, error) {
	var iff IntFuncFile
	data, err := ioutil.ReadFile(file_name)
	if err != nil {
		fmt.Print(err)
		return iff, err
	}

	err = json.Unmarshal(data, &iff)
	if err != nil {
		fmt.Println("error:", err)
		return iff, err
	}
	return iff, nil
}

func SaveIntFuncFile(file_name string, iff IntFuncFile) error {
	data, err := json.MarshalIndent(&iff, " ", "\t")
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	err1 := ioutil.WriteFile(file_name, data, 0644)
	return err1
}

// стек данных
type Stack_data struct {
	Stack []*Cell
}

var Nil *Cell = &Cell{Type: Cell_sym, Value_sym: "nil"}
var True *Cell = &Cell{Type: Cell_sym, Value_sym: "true"}
var Enter *Cell = &Cell{Type: Cell_sym, Value_sym: "Enter"}
var False *Cell = &Cell{Type: Cell_sym, Value_sym: "false"}

var T_int *Cell = &Cell{Type: Cell_sym, Value_sym: "int"}
var T_float *Cell = &Cell{Type: Cell_sym, Value_sym: "float"}
var T_sym *Cell = &Cell{Type: Cell_sym, Value_sym: "sym"}
var T_str *Cell = &Cell{Type: Cell_sym, Value_sym: "str"}
var T_cell *Cell = &Cell{Type: Cell_sym, Value_sym: "cell"}
var T_dict *Cell = &Cell{Type: Cell_sym, Value_sym: "dict"}
var T_array *Cell = &Cell{Type: Cell_sym, Value_sym: "array"}
var T_func *Cell = &Cell{Type: Cell_sym, Value_sym: "func"}
var T_ext *Cell = &Cell{Type: Cell_sym, Value_sym: "ext"}
var T_channel *Cell = &Cell{Type: Cell_sym, Value_sym: "channel"}
var T_next *Cell = &Cell{Type: Cell_sym, Value_sym: "new_ext"}

var Multiple *Cell = &Cell{Type: Cell_multiple}

type Extention struct {
	Type int
	Data interface{}
}

// единица представления информации - стек данных представлен массивом этих единиц
type Cell struct {
	Type            int
	Value_int       int64
	Value_float     float64
	Value_sym       string
	Value_str       string
	Value_head      *Cell
	Value_last      *Cell
	Value_dict      map[string]*Cell
	Value_array     []*Cell
	Value_func      *Func
	Value_ext       string
	Value_channel   chan *Cell
	Value_extention *Extention
	Value_dict_n    map[*Cell]*Cell
}

func CheckType(c *Cell, t int) bool {
	if c.Type == t {
		return true
	}
	return false
}

type Env struct {
	Dict  *Cell
	Stack *Cell
}

type CompilerEnv struct {
	CurrentCellList *Cell
	CurrentFuncCode []Func_code
	VarList         []string
	ConstList       []*Cell
	VarDict         map[string]int64
	ArgList         []string //?? зачем?
	FuncList        []*Func
	Stack           *[]CompilerEnv
}

type GlobalCompilerEnv struct {
	Func_dict map[string]*Func //словарь глобальных функций
}

func InitGlobalCompilerEnv() *GlobalCompilerEnv {
	var gce GlobalCompilerEnv
	gce.Func_dict = make(map[string]*Func)
	return &gce
}

func (gce *GlobalCompilerEnv) FuncToEnv(f *Func) {
	gce.Func_dict[f.Name.String(false)] = f
}
