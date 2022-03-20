#!/bin/bash
cd ./Compiler_Main
source makl.sh
cd ..
cd ./VM_Main
source makl.sh
cd ..
cd ../Shared/std_linux/lib_m.so
source mak_c_l.sh
cd ../..
