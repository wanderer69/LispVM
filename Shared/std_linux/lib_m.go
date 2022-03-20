package main

import (
	"fmt"
	. "arkhangelskiy-dv.ru/LispVM/Shared/common"
	//"errors"
	"syscall"
	"strings"
)

// exported
func LoadStringDataCell(pa []CellShort, res *int) *CellShort {
        a0, err := pa[0].Str()
        if err != nil {
            *res = -1
            return nil
        }
	fmt.Printf("Hello Universe %v\r\n", a0)

	var cm CellShort
        cm.Type = CellShort_sym
        cm.Value_sym = "tested"
        *res = 0
        return &cm
}

/*
// Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
	O_RDONLY int = syscall.O_RDONLY // open the file read-only.
	O_WRONLY int = syscall.O_WRONLY // open the file write-only.
	O_RDWR   int = syscall.O_RDWR   // open the file read-write.
	// The remaining values may be or'ed in to control behavior.
	O_APPEND int = syscall.O_APPEND // append data to the file when writing.
	O_CREATE int = syscall.O_CREAT  // create a new file if none exists.
	O_EXCL   int = syscall.O_EXCL   // used with O_CREATE, file must not exist.
	O_SYNC   int = syscall.O_SYNC   // open for synchronous I/O.
	O_TRUNC  int = syscall.O_TRUNC  // truncate regular writable file when opened.
*/

/*
FileOpen file_name -> file_descriptor
FileRead file_descriptor pos size -> buffer
FileWrite file_descriptor buffer pos size 
FileClose file_descriptor
*/
func FileOpen(pa []CellShort, res *int) *CellShort {
        a0, err1 := pa[0].Str()
        if err1 != nil {
            *res = -1
            return nil
        }
        a1, err2 := pa[1].Str()
        if err2 != nil {
            *res = -2
            return nil
        }

        mode := 0
        a1l := strings.Split(a1, "|")
        for i, _ := range a1l {
             switch a1l[i] {
                  case "ro":
			mode = mode | syscall.O_RDONLY // open the file read-only.
                  case "wo":
			mode = mode | syscall.O_WRONLY // open the file write-only.
                  case "rw":
			mode = mode | syscall.O_RDWR   // open the file read-write.
			// The remaining values may be or'ed in to control behavior.
                  case "a":
			mode = mode | syscall.O_APPEND // append data to the file when writing.
                  case "c":
			mode = mode | syscall.O_CREAT  // create a new file if none exists.
                  case "ex":
			mode = mode | syscall.O_EXCL   // used with O_CREATE, file must not exist.
                  case "s":
			mode = mode | syscall.O_SYNC   // open for synchronous I/O.
                  case "t":
			mode = mode | syscall.O_TRUNC  // truncate regular writable file when opened.
             }
        }
        f, err3 := syscall.Open(a0, mode, 0755)
	if err3 != nil {
            *res = -101
            return nil		
	}
	var cm CellShort
        cm.Type = CellShort_int
        cm.Value_int = int64(f)
        *res = 0
	return &cm;
}

func FileClose(pa []CellShort, res *int) *CellShort {
        a0, err := pa[0].Int()
        if err != nil {
            *res = -1
            return nil
        }
        err = syscall.Close(a0)
	if err != nil {
            *res = -101
            return nil		
	}
	var cm CellShort
        cm.Type = CellShort_int
        cm.Value_int = 0
        *res = 0
	return &cm;
}

func FileRead(pa []CellShort, res *int) *CellShort {
        a0, err1 := pa[0].Int()
        if err1 != nil {
            *res = -1
            return nil
        }
        a1, err2 := pa[1].Int()
        if err2 != nil {
            *res = -2
            return nil
        }
        buf := make([]byte, a1) 
        _, err3 := syscall.Read(a0, buf)
	if err3 != nil {
            *res = -102
            return nil		
	}

	var cm CellShort
        cm.Type = CellShort_str
        cm.Value_str = string(buf)
        *res = 0
        return &cm
}

func FileWrite(pa []CellShort, res *int) *CellShort {
        a0, err1 := pa[0].Int()
        if err1 != nil {
            *res = -1
            return nil
        }
        a1, err2 := pa[1].Str()
        if err2 != nil {
            *res = -2
            return nil
        }
        buf := []byte(a1) 
        _, err3 := syscall.Write(a0, buf)
	if err3 != nil {
            *res = -102
            return nil		
	}

	var cm CellShort
        cm.Type = CellShort_str
        cm.Value_str = string(buf)
        *res = 0
        return &cm
}

func FilePos(pa []CellShort, res *int) *CellShort {
        a0, err1 := pa[0].Int()
        if err1 != nil {
            *res = -1
            return nil
        }
        a1, err2 := pa[1].Int()
        if err2 != nil {
            *res = -2
            return nil
        }
        a2, err3 := pa[2].Int()
        if err3 != nil {
            *res = -3
            return nil
        }
        var pos int64
        if((a2 >= 0) || (a2 < 3)) {
                pos_, err := syscall.Seek(a0, int64(a1), a2)
		if err != nil {
                    *res = -103
                    return nil		
		}
		pos = pos_
        }

	var cm CellShort
        cm.Type = CellShort_int
        cm.Value_int = int64(pos)
        *res = 0
        return &cm
}

func FileSize(pa []CellShort, res *int) *CellShort {
        a0, err1 := pa[0].Int()
        if err1 != nil {
            *res = -1
            return nil
        }
        var pos int64
        pos_, err := syscall.Seek(a0, 0, 1) // SEEK_CUR
	if err != nil {
            *res = -103
            return nil		
	}
	pos = pos_

        pos_, err = syscall.Seek(a0, 0, 2) // SEEK_END
	if err != nil {
            *res = -103
            return nil		
	}

        _, err = syscall.Seek(a0, pos, 0) // SEEK_SET
	if err != nil {
            *res = -103
            return nil		
	}

	var cm CellShort
        cm.Type = CellShort_int
        cm.Value_int = int64(pos_)
        *res = 0
        return &cm
}

/*
EXPORT CellShort* CDECL FileWrite(P_CellShort *pa, unsigned int l_in, unsigned int *l_out, unsigned int *res)
{
  CellShort *a = *pa;
  if(l_in != 2) {
      *res = -1;
      return NULL;
  }
  unsigned char *buf = GoString2CharA(a[1]);

  // printf("FileWrite 0 t %d value %d\r\n", a[0].Type, a[0].Value_int);
  // printf("1 t %d value %s\r\n", a[1].Type, buf);
  FILE *f = (FILE *)a[0].Value_int;
  unsigned int size = (unsigned int)a[1].Value_str.n;

  // unsigned char *buf = (unsigned char *)malloc(sizeof(unsigned char) * (size + 1));

  // printf("-> %p\r\n", buf);

  size_t r = fwrite(buf, sizeof(unsigned char), size, f);
  if (r >= 0) {
  } else {
      *res = -2;
      return NULL;
  }
  // printf("-> %s\r\n", buf);
  *res = 0;
  *l_out = 1;
  CellShort *b = (CellShort *)malloc(sizeof(CellShort));
  b->Type = 2;
  b->Value_int = r;
  return b;
}

EXPORT CellShort* CDECL FilePos(P_CellShort *pa, unsigned int l_in, unsigned int *l_out, unsigned int *res)
{
  CellShort *a = *pa;
  if(l_in != 3) {
      *res = -1;
      return NULL;
  }
  //printf("FilePos 0 t %d value %d\r\n", a[0].Type, a[0].Value_int);
  //printf("1 t %d value %d\r\n", a[1].Type, a[1].Value_int);
  //printf("2 t %d value %d\r\n", a[2].Type, a[2].Value_int);
  FILE *f = (FILE *)a[0].Value_int;
  unsigned int offset = (unsigned int)a[1].Value_int;
  int origin = (unsigned int)a[2].Value_int;

  if((origin >= 0) || (origin < 3)) {
  	fseek(f, offset, origin);
  }

  long l;
  int i;
  fpos_t *pos; // fpos_t определен в stdio.h
  pos = (fpos_t *)(&l);
  fgetpos (f, pos);

  *res = 0;
  *l_out = 1;
  CellShort *b = (CellShort *)malloc(sizeof(CellShort));
  b->Type = 2;
  b->Value_int = l;
  return b;
}

EXPORT CellShort* CDECL FileSize(P_CellShort *pa, unsigned int l_in, unsigned int *l_out, unsigned int *res)
{
  CellShort *a = *pa;
  if(l_in != 1) {
      *res = -1;
      return NULL;
  }
  // printf("FileSize 0 t %d value %d\r\n", a[0].Type, a[0].Value_int);
  FILE *f = (FILE *)a[0].Value_int;
  fseek(f, 0, 2);
  long l;
  int i;
  fpos_t *pos; // fpos_t определен в stdio.h
  pos = (fpos_t *)(&l);
  fgetpos (f, pos);

  *res = 0;
  *l_out = 1;
  CellShort *b = (CellShort *)malloc(sizeof(CellShort));
  b->Type = 2;
  b->Value_int = l;
  return b;
}
*/
