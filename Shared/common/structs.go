package structs

import (
//	"fmt"
	"errors"
)

type CellShort struct {
	Type        int32
	Value_int   int64
	Value_float float64
	Value_sym   string
	Value_str   string
	Value_head  *CellShort
	Value_last  *CellShort
	Value_dict  map[string]*CellShort
	Value_array []*CellShort
}

const CellShort_sym = 1
const CellShort_int = 2
const CellShort_float = 3
const CellShort_str = 4
const CellShort_cell = 5
const CellShort_dict = 6
const CellShort_array = 7
const CellShort_func = 8
const CellShort_object = 9
//const CellShort_channel = 10
const CellShort_dict_n = 11
//const CellShort_multiple = 12

func (cs CellShort) Int() (int, error) {
    if cs.Type != CellShort_int {
          return 0, errors.New("Wrong type int")
    }
    return int(cs.Value_int), nil
}

func (cs CellShort) Float() (float64, error) {
    if cs.Type != CellShort_float {
          return 0, errors.New("Wrong type float")
    }
    return cs.Value_float, nil
}

func (cs CellShort) Sym() (string, error) {
    if cs.Type != CellShort_sym {
          return "", errors.New("Wrong type sym")
    }
    return cs.Value_sym, nil
}

func (cs CellShort) Str() (string, error) {
    if cs.Type != CellShort_str {
          return "", errors.New("Wrong type str")
    }
    return cs.Value_str, nil
}

func (cs CellShort) Cell_last() (*CellShort, error) {
    if cs.Type != CellShort_cell {
          return nil, errors.New("Wrong type cell")
    }
    return cs.Value_last, nil
}

func (cs CellShort) Cell_head() (*CellShort, error) {
    if cs.Type != CellShort_cell {
          return nil, errors.New("Wrong type cell")
    }
    return cs.Value_head, nil
}

func (cs CellShort) Dict() (map[string]*CellShort, error) {
    if cs.Type != CellShort_cell {
          return nil, errors.New("Wrong type dict")
    }
    return cs.Value_dict, nil
}

func (cs CellShort) Array() ([]*CellShort, error) {
    if cs.Type != CellShort_array {
          return nil, errors.New("Wrong type array")
    }
    return cs.Value_array, nil
}
/*
	Type        int32
	Value_int   int64
	Value_float float64
	Value_sym   string
	Value_str   string
	Value_head  *CellShort
	Value_last  *CellShort
	Value_dict  map[string]*CellShort
	Value_array []*CellShort
*/
type TExtFunc func (pa []CellShort, res *int) *CellShort
