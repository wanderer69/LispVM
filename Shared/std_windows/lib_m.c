/* add_basic.c
   Demonstrates creating a DLL with an exported function, the inflexible way.
*/

#include <math.h>
#include <string.h>
#include <stdlib.h>
#include <stdio.h>
#include <inttypes.h>

#if defined(_WINDOWS)
    #define EXPORT __declspec(dllexport)
    #define CDECL __cdecl
#elif defined(_LINUX)
    #define EXPORT __attribute__((visibility("default")))
    #define IMPORT
    #define CDECL
#else
    //  do nothing and hope for the best?
    #define EXPORT
    #define IMPORT
    #pragma warning Unknown dynamic link import/export semantics.
#endif

#include "lib_m.h"

int PrintTestC(unsigned char * s)
{
  printf("'%s'\r\n", s);
  return strlen(s);
}

unsigned char * GoString2CharA(CellShort a)
{
  unsigned char *ss = (unsigned char *)malloc(a.Value_str.n + 1);
  memset(ss, 0, a.Value_str.n + 1);
  memcpy(ss, a.Value_str.p, a.Value_str.n);
//  printf("'%s'\r\n", ss);
  return ss;
}

EXPORT CellShort* CDECL LoadStringDataCell(P_CellShort *pa, unsigned int l_in, unsigned int *l_out, unsigned int *res)
{
  CellShort *a = *pa;
  unsigned char *ss = GoString2CharA(a[0]);
/*
  unsigned char *ss = (unsigned char *)malloc(a[0].Value_str.n + 1);
  memset(ss, 0, a[0].Value_str.n + 1);
  memcpy(ss, a[0].Value_str.p, a[0].Value_str.n);
  printf("'%s'\r\n", ss);
*/
  printf("LoadStringDataCell 0 t %d n %ld '%s'\r\n", a[0].Type, a[0].Value_str.n, ss);
  printf("1 %d %" PRId64 "\r\n", a[1].Type, a[1].Value_int);

  *l_out = (unsigned int)(l_in);
  CellShort *b = (CellShort *)malloc(sizeof(CellShort));
  b->Type = a[0].Type;
  b->Value_str.p = "Nortrop";
  b->Value_str.n = 7;
  *res = 0;
  free(ss);
  printf("-> %d %s\r\n", b->Type, b->Value_str.p);
  return b;
}

/*
FileOpen file_name -> file_descriptor
FileRead file_descriptor pos size -> buffer
FileWrite file_descriptor buffer pos size 
FileClose file_descriptor
*/
EXPORT CellShort* CDECL FileOpen(P_CellShort *pa, unsigned int l_in, unsigned int *l_out, unsigned int *res)
{
  CellShort *a = *pa;
//  printf("l_in %d\r\n", l_in);
  if(l_in != 2) {
      *res = -1;
      return NULL;
  }
  unsigned char *file_name = GoString2CharA(a[0]);

//  printf("FileOpen 0 t %d n %d '%s'\r\n", a[0].Type, a[0].Value_str.n, file_name);
/*
  unsigned char *file_name = (unsigned char *)malloc(a[0].Value_str.n + 1);
  memset(file_name, 0, a[0].Value_str.n + 1);
  memcpy(file_name, a[0].Value_str.p, a[0].Value_str.n);
  printf("'%s'\r\n", file_name);
*/
  unsigned char *mode = GoString2CharA(a[1]);
//  printf("1 t %d n %d '%s'\r\n", a[1].Type, a[1].Value_str.n, mode);
/*
  unsigned char *mode = (unsigned char *)malloc(a[1].Value_str.n + 1);
  memset(mode, 0, a[1].Value_str.n + 1);
  memcpy(mode, a[1].Value_str.p, a[1].Value_str.n);
  printf("'%s'\r\n", mode);
*/
  FILE *f = fopen(file_name, mode);
  if (f == NULL) {
      *res = -2;
      return NULL;
  }
  *l_out = 1;
  *res = 0;
  // a[0].Value_float = sqrt(a[0].Value_float);
  CellShort *b = (CellShort *)malloc(sizeof(CellShort));
  b->Type = 2;
  b->Value_int = (GoInt64)f;
  free(file_name);
  free(mode);
//  printf("-> %d %d\r\n", b->Type, b->Value_int);
  return b;
}

EXPORT CellShort* CDECL FileClose(P_CellShort *pa, unsigned int l_in, unsigned int *l_out, unsigned int *res)
{
  CellShort *a = *pa;
//  printf("l_in %d\r\n", l_in);
  if(l_in != 1) {
      *res = -1;
      return NULL;
  }
//  printf("FileClose 0 t %d value %d\r\n", a[0].Type, a[0].Value_int);
  FILE *f = (FILE *)a[0].Value_int;
  fclose(f);
  *res = 0;
  *l_out = 1;
  CellShort *b = (CellShort *)malloc(sizeof(CellShort));
  b->Type = 2;
  b->Value_int = 0;
  return b;
}

EXPORT CellShort* CDECL FileRead(P_CellShort *pa, unsigned int l_in, unsigned int *l_out, unsigned int *res)
{
  CellShort *a = *pa;
//  printf("l_in %d \r\n", l_in);
  if(l_in != 2) {
      *res = -1;
      return NULL;
  }
//  printf("FileRead 0 t %d value %d\r\n", a[0].Type, a[0].Value_int);
//  printf("1 t %d value %d\r\n", a[1].Type, a[1].Value_int);
  FILE *f = (FILE *)a[0].Value_int;
  unsigned int size = (unsigned int)a[1].Value_int;

  unsigned char *buf = (unsigned char *)malloc(sizeof(unsigned char) * (size + 1));

  size_t r = fread(buf, sizeof(unsigned char), size, f);
  if (r >= 0) {
  } else {
      *res = -2;
      return NULL;
  }
  // printf("-> %s\r\n", buf);
  *l_out = 1;
  *res = 0;
  CellShort *b = (CellShort *)malloc(sizeof(CellShort));
  b->Type = 4;
  b->Value_str.p = buf;
  b->Value_str.n = r;
  return b;
}

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
  fpos_t *pos; /* fpos_t определен в stdio.h */
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
  fpos_t *pos; /* fpos_t определен в stdio.h */
  pos = (fpos_t *)(&l);
  fgetpos (f, pos);

  *res = 0;
  *l_out = 1;
  CellShort *b = (CellShort *)malloc(sizeof(CellShort));
  b->Type = 2;
  b->Value_int = l;
  return b;
}

/*
EXPORT CellShort* CDECL FileOpen(CellShort a[], unsigned int l_in, unsigned int *l_out)
{
  printf("LoadStringDataCell %d %d\r\n", a[0].Type, a[0].Value_str.n);
  unsigned char *ss = (unsigned char *)malloc(a[0].Value_str.n + 1);
  memset(ss, 0, a[0].Value_str.n + 1);
  memcpy(ss, a[0].Value_str.p, a[0].Value_str.n);
  printf("'%s'\r\n", ss);
  *l_out = (unsigned int)(l_in);
  // a[0].Value_float = sqrt(a[0].Value_float);
  CellShort *b = (CellShort *)malloc(sizeof(CellShort));
  b->Type = a[0].Type;

//  GoString *gs = (GoString *)malloc(sizeof(GoString));
//  gs->p = "Nortrop";
//  gs->n = 7;
//  b->Value_str = gs;

  b->Value_str.p = "Nortrop";
  b->Value_str.n = 7;
  printf("-> %d %s\r\n", b->Type, b->Value_str.p);
  return b;
}
*/
