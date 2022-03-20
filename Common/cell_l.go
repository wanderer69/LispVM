package common

import (
	"fmt"
	"strconv"

	//	"unicode/utf8"
	"encoding/json"
	//	"hash/maphash"
	"strings"
)

func (c *Cell) IntToStr() string {
	return ""
}

func (c *Cell) StrToStr() string {
	return ""

}

func (c *Cell) HeadLastToStr(level int) string {
	var str_l []string
	if level == 0xFFFFFFFF {
		if c.Type == Cell_cell {
			str_l = append(str_l, "(")
		}
		str_l = append(str_l, ")")

		if c.Type == Cell_cell {
			str_l = append(str_l, "(")
		}
	}
	return ""
}

func (c *Cell) DictToStr(level int) string {
	return ""
}

func (c *Cell) ArrayToStr(level int) string {
	return ""
}

func (c *Cell) FuncToStr(level int) string {
	return ""
}

func (c *Cell) CellToStr() string {
	return ""
}

func CellCons(h *Cell, l *Cell) *Cell {
	res := &Cell{Type: Cell_cell, Value_head: h, Value_last: l}
	return res
}

func CellCreateFromInt(i int) *Cell {
	res := &Cell{Type: Cell_int, Value_int: int64(i)}
	return res
}

func CreateCellInt(i int64) *Cell {
	res := &Cell{Type: Cell_int, Value_int: int64(i)}
	return res
}

func CreateCellFloat(i float64) *Cell {
	res := &Cell{Type: Cell_float, Value_float: float64(i)}
	return res
}

func CreateCellSym(str string) *Cell {
	res := &Cell{Type: Cell_int, Value_sym: str}
	return res
}

func CreateCellStr(str string) *Cell {
	res := &Cell{Type: Cell_int, Value_str: str}
	return res
}

/*
func CreateCell(i ) *Cell {
	res := &Cell{Type: Cell_, Value_: i}
	return res
}
*/
func CreateCellCell(h *Cell, l *Cell) *Cell {
	res := &Cell{Type: Cell_cell, Value_head: h, Value_last: l}
	return res
}

func CreateCellDict(i map[string]*Cell) *Cell {
	res := &Cell{Type: Cell_dict, Value_dict: i}
	return res
}

func CreateCellArray(i []*Cell) *Cell {
	res := &Cell{Type: Cell_array, Value_array: i}
	return res
}

func CreateCellFunc(i *Func) *Cell {
	res := &Cell{Type: Cell_func, Value_func: i}
	return res
}

func CreateCellExt(i string) *Cell {
	res := &Cell{Type: Cell_ext, Value_ext: i}
	return res
}

func CreateCellСhannel(i chan *Cell) *Cell {
	res := &Cell{Type: Cell_channel, Value_channel: i}
	return res
}

func CreateCell(i *Extention) *Cell {
	res := &Cell{Type: Cell_next, Value_extention: i}
	return res
}

func CreateCellDict_n(i map[*Cell]*Cell) *Cell {
	res := &Cell{Type: Cell_dict_n, Value_dict_n: i}
	return res
}

func CellCreateFromStr(str_ string) *Cell {
	str := strings.Trim(str_, " \r\n\t")
	//fmt.Printf("CellCreateFromStr str '%v'\r\n", str)
	// смотрим, какие символы - обрамляют
	var res *Cell
	if len(str) > 0 {
		if str[0] == '"' {

			if str[len(str)-1] == '"' {
				res = &Cell{Type: Cell_str, Value_str: str[1 : len(str)-1]}
			}
			//fmt.Printf("CellCreateFromStr res %#v\r\n", res)
		} else {
			if str[0] == '[' && str[len(str)-1] == ']' {
				// это массив с разделителями запятыми
				//fmt.Printf("CellCreateFromStr array '%v'\r\n", str)
			} else {
				if str[0] == '{' && str[len(str)-1] == '}' {
					// это словарь надо поделить на разделы : ключ всегда строка а значение может быть любой объект.

				} else {
					// попробуем число
					f_int, err := strconv.ParseInt(str, 0, 64)
					if err != nil {
						//fmt.Printf("err %v\r\n", err)
						f_float, err1 := strconv.ParseFloat(str, 64)
						//fmt.Printf("err1 '%v' f_float %v \r\n", err1, f_float)
						if err1 == nil {
							res = &Cell{Type: Cell_float, Value_float: f_float}
						} else {
							// это символ
							// fmt.Printf("!!!!!! str %v\r\n", str)
							switch str {
							case "Nil":
								res = Nil
								//fmt.Printf(">> %v\r\n", res)
							case "True":
								res = True
							case "False":
								res = False
							case "nil":
								res = Nil
								//fmt.Printf(">> %v\r\n", res)
							case "true":
								res = True
							case "false":
								res = False
							case "int":
								res = T_int
							case "float":
								res = T_float
							case "sym":
								res = T_sym
							case "str":
								res = T_str
							case "cell":
								res = T_cell
							case "dict":
								res = T_dict
							case "array":
								res = T_array
							case "func":
								res = T_func
							case "ext":
								res = T_ext
							case "channel":
								res = T_channel
							case "new_ext":
								res = T_next
							default:
								res = &Cell{Type: Cell_sym, Value_sym: str}
							}
						}
					} else {
						res = &Cell{Type: Cell_int, Value_int: f_int}
					}
				}
			}
		}
		return res
	}
	return nil
}

func CellCreateFromCell(head *Cell, last *Cell) *Cell {
	return nil
}

func CellCreateFromString(str string) *Cell {
	var res Cell = Cell{Type: Cell_str, Value_str: str}
	return &res
}

func ToSpace(level int) string {
	v := ""
	for i := 0; i < level; i++ {
		v = v + "  "
	}
	return v
}

func (c *Cell) Print(debug bool) string {
	result := c.String(debug)
	fmt.Printf("%v\r\n", result)
	return result
}

func (c *Cell) String(debug bool) string {
	res, _ := c.String_item(0, debug)
	return res
}

func CellCreateExt(id []byte, name string, strl []string) *Cell {
	et := Ext_Type_store{id, name, strl}
	b, err := json.Marshal(et)
	if err != nil {
		fmt.Printf("%v\r\n", err)
		return nil
	}
	var res Cell = Cell{Type: Cell_ext, Value_ext: string(b)}
	return &res
}

func (c *Cell) GetCellExtValue() ([]byte, string, []string, bool) {
	var et Ext_Type_store
	if c.Type == Cell_ext {
		err := json.Unmarshal([]byte(c.Value_ext), &et)
		if err != nil {
			fmt.Printf("%v\r\n", err)
			return nil, "", nil, false
		}
		return et.ID, et.Name, et.Value, true
	}
	return nil, "", nil, false
}

func CellCreateCell(h *Cell, l *Cell) *Cell {
	var res Cell = Cell{Type: Cell_cell, Value_head: h, Value_last: l}
	return &res
}

func CellCreateArrayEmpty() *Cell {
	var res Cell = Cell{Type: Cell_array, Value_array: []*Cell{}}
	return &res
}

func CellCreateArray(l *Cell) *Cell {
	var res Cell = Cell{Type: Cell_array, Value_array: []*Cell{l}}
	return &res
}

func CellAppendArray(a *Cell, l *Cell) *Cell {
	a.Value_array = append(a.Value_array, l)
	return a
}

func CellCreateDict() *Cell {
	dict := make(map[string]*Cell)
	var res Cell = Cell{Type: Cell_dict, Value_dict: dict}
	return &res
}

func CellCreateDictN() *Cell {
	dict := make(map[*Cell]*Cell)
	var res Cell = Cell{Type: Cell_dict_n, Value_dict_n: dict}
	return &res
}

func CellAddDict(d *Cell, k *Cell, v *Cell) *Cell {
	/*
		var h2 maphash.Hash
		h2.SetSeed(MakeSeed())
		h2.WriteString(k.String(false))
		kk := fmt.Sprintf("%x", h2.Sum64())
	*/
	kk := ""
	switch k.Type {
	case Cell_str:
		kk = k.Value_str
	case Cell_sym:
		kk = k.Value_sym
	default:
		kk = k.String(false)
	}
	d.Value_dict[kk] = v
	return d
}

func CellAddDictN(d *Cell, k *Cell, v *Cell) *Cell {
	d.Value_dict_n[k] = v
	return d
}

func (c *Cell) String_item(level int, debug bool) (string, bool) {
	var res bool
	res = false
	result := ""
	if debug {
		fmt.Printf("c.Type %v\r\n", c.Type)
	}
	switch c.Type {
	case Cell_sym:
		result = c.Value_sym
		if debug {
			fmt.Printf("%v%v\r\n", ToSpace(level), c.Value_sym)
		}
	case Cell_int:
		result = fmt.Sprintf("%v", c.Value_int)
		if debug {
			fmt.Printf("%v%v\r\n", ToSpace(level), c.Value_int)
		}
	case Cell_str:
		result = fmt.Sprintf("\"%s\"", c.Value_str)
		if debug {
			fmt.Printf("%v'%v'\r\n", ToSpace(level), c.Value_str)
		}
	case Cell_float:
		result = fmt.Sprintf("%f", c.Value_float)
		if debug {
			fmt.Printf("%v%f\r\n", ToSpace(level), c.Value_float)
		}
	case Cell_cell:
		// печатаем новый уровень который начинается со скобки
		h := c.String_level(level+1, true, debug)
		result = fmt.Sprintf("\r\n%v%s", ToSpace(level), h)
	case Cell_dict:
		result = result + "{"
		i := 0
		for k, v := range c.Value_dict {
			s, res1 := v.String_item(level, debug)
			res = res1
			if i > 0 {
				result = result + ", " + fmt.Sprintf("%s:%s", k, s)
			} else {
				result = result + fmt.Sprintf("%s:%s", k, s)
			}
			i = i + 1
		}
		result = result + "}"
	case Cell_array:
		result = result + "["
		for i, v := range c.Value_array {
			s, res1 := v.String_item(level, debug)
			res = res1
			if i > 0 {
				result = result + ", " + fmt.Sprintf("%s", s)
			} else {
				result = result + fmt.Sprintf("%s", s)
			}
		}
		result = result + "]"
	case Cell_func:
		result = fmt.Sprintf("<%s>", c.Value_func)

	case Cell_ext:
		result = fmt.Sprintf("<%s>", c.Value_ext)

	case Cell_next:
		result = fmt.Sprintf("<%s>", c.Value_extention)
	case Cell_multiple:
		result = fmt.Sprintf("$%v$", c.Value_int)
	}

	return result, res
}

func (c *Cell) String_level(level int, flag bool, debug bool) string {
	// новый уровень добавляем скобку
	h := ""
	tl := 0
	cc := c
	result := ""
	for {
		// цикл по уровню пока не Nil хвост
		if cc.Value_head != nil {
			h, _ = cc.Value_head.String_item(level, debug)
			if len(result) > 0 {
				result = result + " " + h
			} else {
				result = "(" + h
			}
			if debug {
				fmt.Printf("%vh>%v\r\n", ToSpace(level), h)
			}
		}
		if cc.Value_last != nil {
			tl = cc.Value_last.Type
			if tl == Cell_cell {
				cc = cc.Value_last
			} else {
				if cc.Value_last == Nil {
					result = result + ")"
					break
				} else {
					h, _ = cc.Value_last.String_item(level, debug)
					result = result + "." + h
					break
				}
			}
		}
	}
	return result
}
