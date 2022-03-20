gcc -D_WINDOWS -c -o lib_m.o lib_m.c
gcc -o file.dll -s -shared lib_m.o -Wl,--subsystem,windows
