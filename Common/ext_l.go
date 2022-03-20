package common

//	"fmt"
//"strconv"
//"strings"

func CreateExtCell(val interface{}) *Cell {
	c := Cell{Type: Cell_next}
	e := Extention{Type: Cell_next, Data: val}
	c.Value_extention = &e
	return &c
}
