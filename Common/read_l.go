package common

import (
	"fmt"
	//	"strconv"
	"unicode/utf8"
	//. "../Common"
)

func Load_list(pos_begin int, pos_end int, lent int, text string, level int, flag string, debug int) (*Cell, int, int, int) {
	// debug = 20
	prev := pos_begin
	i_prev := pos_begin
	var root_cell *Cell = nil
	var current_cell *Cell = nil
	var current_key *Cell = nil
	if debug > 2 {
		fmt.Printf("<-- text '%v'\r\n", text[pos_begin:pos_end])
	}
	flag_s := false
	g_s_flag := false
	//pos_begin_n := pos_begin
	for i, w := pos_begin, 0; i < pos_end; {
		runeValue, width := utf8.DecodeRuneInString(text[i:])
		if debug > 4 {
			fmt.Printf("%v%#U starts at byte position %d %v level %v\n", LevelShift(level), runeValue, i, string(runeValue), level)
		}
		w = width
		s1 := string(runeValue)
		if !g_s_flag {
			// если это скобка - строим новую ячейку, если на этом уровне нет или идем вглубь.
			if s1 == "(" {
				if debug > 13 {
					fmt.Printf("i %v\r\n", i)
				}
				if root_cell == nil {
					mi := Cell{}
					mi.Value_head = Nil
					mi.Value_last = Nil
					mi.Type = Cell_cell
					root_cell = &mi
					current_cell = root_cell
				} else {
					level_g := 1
					prev = i
					i = i + w
					for {
						runeValue, width := utf8.DecodeRuneInString(text[i:])
						if debug > 4 {
							fmt.Printf("!%#U starts at byte position %d %v\n", runeValue, i, string(runeValue))
						}
						w = width
						s1 := string(runeValue)
						if s1 == "(" {
							level_g = level_g + 1
						} else {
							if s1 == ")" {
								level_g = level_g - 1
								if level_g > 0 {

								} else {
									break
								}
							}
						}
						i = i + w
						if i > len(text) {
							fmt.Printf("Error i %v len(text) %v\r\n", i, len(text))
							panic("Error!")
						}
					}
					if debug > 12 {
						fmt.Printf("text '%v' prev %v i %v w %v\r\n", text[prev:i+w], prev, i, w)
						fmt.Printf("prev %v, i+w %v, i-prev %v, text %v, level+1 %v\r\n", prev, i+w, i+w-prev, text, level+1)
					}
					list_cell, _, _, _ := Load_list(prev, i+w, i+w-prev, text, level+1, "(", debug)
					if current_cell.Value_head != Nil {
						mi := Cell{}
						mi.Value_head = list_cell
						mi.Value_last = Nil
						mi.Type = Cell_cell
						current_cell.Value_last = &mi
						current_cell = &mi
					} else {
						current_cell.Value_head = list_cell
						current_cell.Value_last = Nil
					}
				}
				i = i + w
			} else if (s1 == " ") || (s1 == "\t") || (s1 == "\r") || (s1 == "\n") {
				// это разделитель! если до этого были отличные символы - строим строку.
				if debug > 5 {
					fmt.Printf("space flag_s %v\r\n", flag_s)
				}
				if flag_s {
					flag_s = false
					if debug > 12 {
						fmt.Printf("i_prev %v, i %v, w %v\r\n", i_prev, i, w)
					}
					if debug > 5 {
						fmt.Printf("text separator %v\r\n", text[i_prev:i])
					}
					n := CellCreateFromStr(GetSlice(text, i_prev, i))
					if debug > 12 {
						fmt.Printf("mi.Value %v\r\n", n.String(false))
					}
					if current_cell != nil {
						if current_cell.Value_head != Nil {
							mi := Cell{}
							mi.Value_head = n
							mi.Value_last = Nil
							mi.Type = Cell_cell
							current_cell.Value_last = &mi
							current_cell = &mi
						} else {
							current_cell.Value_head = n
							current_cell.Value_last = Nil
						}
					} else {
						mi := Cell{}
						mi.Value_head = n
						mi.Value_last = Nil
						mi.Type = Cell_cell
						root_cell = &mi
						current_cell = root_cell
					}
				} else {
					// нет ничего это повторный разделитель
				}
				i = i + w
			} else if s1 == ")" {
				// скобка закрывающая завершаем работу и выходим
				if debug > 5 {
					fmt.Printf("close bracket flag_s %v\r\n", flag_s)
				}
				if flag_s {
					flag_s = false
					if debug > 12 {
						fmt.Printf("i_prev %v, i %v\r\n", i_prev, i)
					}
					if debug > 5 {
						fmt.Printf("text %v\r\n", text[i_prev:i])
					}
					n := CellCreateFromStr(GetSlice(text, i_prev, i))
					if debug > 12 {
						fmt.Printf("mi.Value ) %v\r\n", n.String(false))
					}
					if current_cell != nil {
						if current_cell.Value_head != Nil {
							mi := Cell{}
							mi.Value_head = n
							mi.Value_last = Nil
							mi.Type = Cell_cell
							current_cell.Value_last = &mi
							current_cell = &mi
						} else {
							current_cell.Value_head = n
							current_cell.Value_last = Nil
						}
						//
					} else {
						mi := Cell{}
						mi.Value_head = n
						mi.Value_last = Nil
						mi.Type = Cell_cell
						root_cell = &mi
						current_cell = root_cell
					}
				}
				if debug > 2 {
					if root_cell != nil {
						ss := root_cell.Print(false)
						fmt.Printf("--> %v\r\n", ss)
					} else {
						fmt.Printf("--> nil!!!\r\n")
					}
				}
				return root_cell, i + 1, pos_end, 0
			} else if s1 == "[" {
				if debug > 13 {
					fmt.Printf("[ i %v\r\n", i)
				}
				if root_cell == nil {
					mi := CellCreateArrayEmpty()
					root_cell = mi
					current_cell = root_cell
					flag = "["
					// fmt.Printf("create array\r\n")
				} else {
					level_g := 1
					prev = i
					i = i + w
					for {
						runeValue, width := utf8.DecodeRuneInString(text[i:])
						if debug > 4 {
							fmt.Printf("[>> !%#U starts at byte position %d %v\r\n", runeValue, i, string(runeValue))
						}
						w = width
						s1 := string(runeValue)
						//fmt.Printf("s1 %v\r\n", s1)
						if s1 == "[" {
							level_g = level_g + 1
						} else {
							if s1 == "]" {
								level_g = level_g - 1
								if level_g > 0 {

								} else {
									break
								}
							}
						}
						i = i + w
						if i > len(text) {
							fmt.Printf("Error i %v len(text) %v\r\n", i, len(text))
							panic("Error!")
						}
					}
					if debug > 12 {
						fmt.Printf("text [ '%v' prev %v i %v w %v\r\n", text[prev:i+w], prev, i, w)
						fmt.Printf("prev [ %v, i+w %v, i-prev %v, text %v, level+1 %v\r\n", prev, i+w, i+w-prev, text, level+1)
					}
					list_cell, _, _, _ := Load_list(prev, i+w, i+w-prev, text, level+1, "[", debug)
					switch current_cell.Type {
					case Cell_cell:
						if current_cell.Value_head != Nil {
							mi := Cell{}
							mi.Value_head = list_cell
							mi.Value_last = Nil
							mi.Type = Cell_cell
							current_cell.Value_last = &mi
							current_cell = &mi
						} else {
							current_cell.Value_head = list_cell
							current_cell.Value_last = Nil
						}
					case Cell_array:
						CellAppendArray(current_cell, list_cell)
						// fmt.Printf("add to array\r\n")
					}
					// fmt.Printf("array loaded\r\n")
				}
				i = i + w
			} else if s1 == "]" {
				// скобка закрывающая завершаем работу и выходим
				if debug > 5 {
					fmt.Printf("close bracket flag_s %v\r\n", flag_s)
				}
				if flag_s {
					//flag_s = false
					if debug > 12 {
						fmt.Printf("i_prev %v, i %v\r\n", i_prev, i)
					}
					if debug > 5 {
						fmt.Printf("text ] %v\r\n", text[i_prev:i])
					}
					n := CellCreateFromStr(GetSlice(text, i_prev, i))
					if debug > 12 {
						fmt.Printf("--> mi.Value %v\r\n", n.String(false))
					}
					if current_cell != nil {
						CellAppendArray(current_cell, n)
					} else {
						mi := CellCreateArray(n)
						root_cell = mi
						current_cell = root_cell
					}
				}
				if debug > 2 {
					if root_cell != nil {
						ss := root_cell.Print(false)
						fmt.Printf("--> %v\r\n", ss)
					} else {
						fmt.Printf("--> nil!!!\r\n")
					}
				}
				return root_cell, i + 1, pos_end, 0
			} else if s1 == "{" {
				if debug > 13 {
					fmt.Printf("i %v\r\n", i)
				}
				if root_cell == nil {
					mi := CellCreateDict()
					root_cell = mi
					current_cell = root_cell
					flag = "{"
				} else {
					level_g := 1
					prev = i
					i = i + w
					for {
						runeValue, width := utf8.DecodeRuneInString(text[i:])
						if debug > 4 {
							fmt.Printf("[>>> !%#U starts at byte position %d %v\r\n", runeValue, i, string(runeValue))
						}
						w = width
						s1 := string(runeValue)
						// fmt.Printf("s1 %v\r\n", s1)
						if s1 == "{" {
							level_g = level_g + 1
						} else {
							if s1 == "}" {
								level_g = level_g - 1
								if level_g > 0 {

								} else {
									break
								}
							}
						}
						i = i + w
						if i > len(text) {
							fmt.Printf("Error i %v len(text) %v\r\n", i, len(text))
							panic("Error!")
						}
					}
					if debug > 12 {
						fmt.Printf("text { '%v' prev %v i %v w %v\r\n", text[prev:i+w], prev, i, w)
						fmt.Printf("prev { %v, i+w %v, i-prev %v, text %v, level+1 %v\r\n", prev, i+w, i+w-prev, text, level+1)
					}
					list_cell, _, _, _ := Load_list(prev, i+w, i+w-prev, text, level+1, "{", debug)
					switch current_cell.Type {
					case Cell_cell:
						if current_cell.Value_head != Nil {
							mi := Cell{}
							mi.Value_head = list_cell
							mi.Value_last = Nil
							mi.Type = Cell_cell
							current_cell.Value_last = &mi
							current_cell = &mi
						} else {
							current_cell.Value_head = list_cell
							current_cell.Value_last = Nil
						}
					case Cell_array:
						CellAppendArray(current_cell, list_cell)
					case Cell_dict:
						CellAddDict(current_cell, current_key, list_cell)
					}
				}
				i = i + w
			} else if s1 == "}" {
				// скобка закрывающая завершаем работу и выходим
				if debug > 5 {
					fmt.Printf("close bracket flag_s %v\r\n", flag_s)
				}
				if flag_s {
					//flag_s = false
					if debug > 12 {
						fmt.Printf("i_prev %v, i %v\r\n", i_prev, i)
					}
					if debug > 5 {
						fmt.Printf("text ] %v\r\n", text[i_prev:i])
					}
					n := CellCreateFromStr(GetSlice(text, i_prev, i))
					if debug > 12 {
						fmt.Printf("mi.Value %v\r\n", n.String(false))
					}
					if current_key != nil {
						if current_cell != nil {
							CellAddDict(current_cell, current_key, n)
						} else {
							mi := CellCreateDict()
							CellAddDict(current_cell, current_key, n)
							root_cell = mi
							current_cell = root_cell
						}
						current_key = nil
					}
				}
				if debug > 2 {
					if root_cell != nil {
						ss := root_cell.Print(false)
						fmt.Printf("--> %v\r\n", ss)
					} else {
						fmt.Printf("--> nil!!!\r\n")
					}
				}
				return root_cell, i + 1, pos_end, 0
			} else if s1 == "\"" {
				i_prev = i // + w
				if flag_s {
					flag_s = false
				} else {
					flag_s = true
				}
				g_s_flag = true
				i = i + w
			} else if s1 == ":" {
				// fmt.Printf("flag %v\r\n", flag)
				if flag == "{" {
					// это разделитель! если до этого были отличные символы - строим строку.
					if debug > 5 {
						fmt.Printf("Comma flag_s %v\r\n", flag_s)
					}
					if flag_s {
						flag_s = false
						if debug > 12 {
							fmt.Printf("i_prev %v, i %v, w %v\r\n", i_prev, i, w)
						}
						if debug > 5 {
							fmt.Printf("text comma %v\r\n", text[i_prev:i])
						}
						n := CellCreateFromStr(GetSlice(text, i_prev, i))
						if debug > 12 {
							fmt.Printf("mi.Value %v\r\n", n.String(false))
						}
						if debug > 12 {
							fmt.Printf("current_cell %v\r\n", current_cell)
						}
						//if current_key == nil {
						current_key = n
						//}
						if debug > 12 {
							fmt.Printf("current_key %v\r\n", current_key.String(false))
						}
					} else {
						// нет ничего это повторный разделитель
						//fmt.Printf("{ flag_s false\r\n")
						/*
							n := CellCreateFromStr(GetSlice(text, i_prev, i))
							current_key = n
						*/
						//fmt.Printf("current_key\r\n")
						if debug > 12 {
							fmt.Printf("current_key %v\r\n", current_key.String(false))
						}
					}
					i = i + w
				} else {
					i = i + w
				}
			} else if s1 == "," { // || s1 == " "
				if flag == "[" {
					// это разделитель! если до этого были отличные символы - строим строку.
					if debug > 5 {
						fmt.Printf("Comma flag_s %v\r\n", flag_s)
					}
					if flag_s {
						flag_s = false
						if debug > 12 {
							fmt.Printf("i_prev %v, i %v, w %v\r\n", i_prev, i, w)
						}
						if debug > 5 {
							fmt.Printf("text comma %v\r\n", text[i_prev:i])
						}
						n := CellCreateFromStr(GetSlice(text, i_prev, i))
						if debug > 12 {
							fmt.Printf("-> mi.Value %v\r\n", n.String(false))
						}
						if debug > 12 {
							fmt.Printf("current_cell %v\r\n", current_cell)
						}
						if current_cell != nil {
							CellAppendArray(current_cell, n)
						} else {
							mi := CellCreateArray(n)
							root_cell = mi
							current_cell = root_cell
						}
						if debug > 12 {
							fmt.Printf("current_cell %v\r\n", current_cell.String(false))
						}
					} else {
						// нет ничего это повторный разделитель
					}
				} else if flag == "{" {
					// это разделитель! если до этого были отличные символы - строим строку.
					if debug > 5 {
						fmt.Printf("Comma { flag_s %v\r\n", flag_s)
					}
					if flag_s {
						flag_s = false
						if debug > 12 {
							fmt.Printf("i_prev %v, i %v, w %v\r\n", i_prev, i, w)
						}
						if debug > 5 {
							fmt.Printf("text comma %v\r\n", text[i_prev:i])
						}
						n := CellCreateFromStr(GetSlice(text, i_prev, i))
						if debug > 12 {
							fmt.Printf("mi.Value %v\r\n", n.String(false))
						}
						if debug > 12 {
							fmt.Printf("current_cell %v\r\n", current_cell)
						}
						if current_key != nil {
							if current_cell != nil {
								CellAddDict(current_cell, current_key, n)
							} else {
								mi := CellCreateDict()
								CellAddDict(current_cell, current_key, n)
								root_cell = mi
								current_cell = root_cell
							}
							current_key = nil
						}
						if debug > 12 {
							fmt.Printf("current_cell %v\r\n", current_cell.String(false))
						}
					} else {
						// нет ничего это повторный разделитель
						fmt.Printf(">>>\r\n")
					}
				}
				i = i + w
			} else {
				// это символ!!
				if !flag_s {
					i_prev = i // + w
					flag_s = true
				} else {
				}
				i = i + w
			}
		} else {
			if s1 == "\"" {
				// это разделитель! если до этого были отличные символы - строим строку.
				if debug > 5 {
					fmt.Printf("space 2 flag_s %v\r\n", flag_s)
				}
				i = i + w
				if flag_s {
					flag_s = false
					if debug > 12 {
						fmt.Printf("i_prev %v, i %v, w %v\r\n", i_prev, i, w)
					}
					if debug > 5 {
						fmt.Printf("text %v\r\n", text[i_prev:i])
					}
					n := CellCreateFromStr(GetSlice(text, i_prev, i))
					if debug > 12 {
						fmt.Printf("---> mi.Value %v\r\n", n.String(false))
					}
					// fmt.Printf("flag %v\r\n", flag)
					if flag == "{" {
						// fmt.Printf("current_cell %v\r\n", current_cell)
						if current_cell != nil {
							// fmt.Printf("current_key %v\r\n", current_key)
							if current_key != nil {
								CellAddDict(current_cell, current_key, n)
							} else {
								// fmt.Printf("No current_key %v\r\n", n)
								current_key = n
							}
						} else {
							mi := CellCreateDict()
							CellAddDict(current_cell, current_key, n)
							root_cell = mi
							current_cell = root_cell
						}
						if debug > 12 {
							fmt.Printf("current_cell %v\r\n", current_cell.String(false))
						}
					} else {
						if flag == "[" {
							if current_cell != nil {
								CellAppendArray(current_cell, n)
							} else {
								mi := CellCreateArray(n)
								root_cell = mi
								current_cell = root_cell
							}
							if debug > 12 {
								fmt.Printf("current_cell %v\r\n", current_cell.String(false))
							}
						} else {
							if current_cell != nil {
								if current_cell.Value_head != Nil {
									mi := CellCreateCell(n, Nil)
									current_cell.Value_last = mi
									current_cell = mi
								} else {
									current_cell.Value_head = n
									current_cell.Value_last = Nil
								}
							} else {
								mi := CellCreateCell(n, Nil)
								root_cell = mi
								current_cell = root_cell
							}
						}
					}
				} else {
					// нет ничего это повторный разделитель
				}
				g_s_flag = false
				//flag_s = true
				i_prev = i // + w
				/*
					} else if s1 == "," { // || s1 == " "
						if debug > 5 {
							fmt.Printf("!!!! Comma flag_s %v\r\n", flag_s)
						}
						if flag == "[" {
							// это разделитель! если до этого были отличные символы - строим строку.
							if debug > 5 {
								fmt.Printf("Comma flag_s %v\r\n", flag_s)
							}
							if flag_s {
								flag_s = false
								if debug > 12 {
									fmt.Printf("i_prev %v, i %v, w %v\r\n", i_prev, i, w)
								}
								if debug > 5 {
									fmt.Printf("text comma %v\r\n", text[i_prev:i])
								}
								n := CellCreateFromStr(GetSlice(text, i_prev, i))
								if debug > 12 {
									fmt.Printf("mi.Value %v\r\n", n.Value_head)
								}
								if current_cell != nil {
									CellAppendArray(current_cell, n)
								} else {
									mi := CellCreateArray(n)
									root_cell = mi
									current_cell = root_cell
								}
							} else {
								// нет ничего это повторный разделитель
							}
						}
						i = i + w
				*/
			} else {
				i = i + w
				flag_s = true
			}
		}

		if i >= pos_end {
			// все закончилось....
			// fmt.Printf("end flag_s %v i_prev %v, i %v\r\n", flag_s, i_prev, i)
			if flag_s {
				flag_s = false
				if debug > 12 {
					fmt.Printf("i_prev %v, i %v\r\n", i_prev, i)
				}
				if debug > 5 {
					fmt.Printf("text all %v\r\n", text[i_prev:i])
				}
				root_cell = CellCreateFromStr(GetSlice(text, i_prev, i))
				if debug > 12 {
					fmt.Printf("mi.Value %v\r\n", root_cell.Value_head)
				}
			}
		}
	}
	return root_cell, 0, 0, 0
}
