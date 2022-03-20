#include <stddef.h> /* for ptrdiff_t below */

#ifndef GO_CGO_EXPORT_PROLOGUE_H
#define GO_CGO_EXPORT_PROLOGUE_H

//#undef GO_CGO_GOSTRING_TYPEDEF

#ifndef GO_CGO_GOSTRING_TYPEDEF
typedef struct { const char *p; ptrdiff_t n; } _GoString_;
#endif

#endif

/* Start of preamble from import "C" comments.  */


/* End of preamble from import "C" comments.  */


/* Start of boilerplate cgo prologue.  */
//#line 1 "cgo-gcc-export-header-prolog"

#ifndef GO_CGO_PROLOGUE_H
#define GO_CGO_PROLOGUE_H

typedef signed char GoInt8;
typedef unsigned char GoUint8;
typedef short GoInt16;
typedef unsigned short GoUint16;
typedef int GoInt32;
typedef unsigned int GoUint32;
typedef long long GoInt64;
typedef unsigned long long GoUint64;
typedef GoInt64 GoInt;
typedef GoUint64 GoUint;
typedef __SIZE_TYPE__ GoUintptr;
typedef float GoFloat32;
typedef double GoFloat64;
typedef float _Complex GoComplex64;
typedef double _Complex GoComplex128;

/*
  static assertion to make sure the file is being used on architecture
  at least with matching size of GoInt.
*/
//typedef char _check_for_64_bit_pointer_matching_GoInt[sizeof(void*)==64/8 ? 1:-1];

#ifndef GO_CGO_GOSTRING_TYPEDEF
typedef _GoString_ GoString;
#endif
typedef void *GoMap;
typedef void *GoChan;
typedef struct { void *t; void *v; } GoInterface;
typedef struct { void *data; GoInt len; GoInt cap; } GoSlice;

#endif

typedef struct CellShort CellShort;
struct CellShort {
GoInt32 Type;
GoInt64 Value_int;
GoFloat64 Value_float;
GoString Value_sym;
GoString Value_str;
CellShort *Value_head;
CellShort *Value_last;
GoMap Value_dict;
GoSlice Value_array;
};

typedef CellShort * P_CellShort;

/* End of boilerplate cgo prologue.  */
/*
#ifdef __cplusplus
extern "C" {
#endif

extern __declspec(dllexport) GoInt sqrt_(GoFloat64 a, GoFloat64* c);
extern __declspec(dllexport) GoInt LoadData(GoSlice file_name_d, GoSlice* data_p);
extern __declspec(dllexport) GoInt ToStruct(GoSlice data, GoSlice* data_p);

#ifdef __cplusplus
}
#endif
*/