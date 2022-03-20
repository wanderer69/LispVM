#!/bin/bash
DIR_DEPLOY="C:/Deploymented" #or /home/user/Deploymented
DIR_LISPVM=$DIR_DEPLOY/LispVM
DIR_BIN=$DIR_LISPVM/bin
FILE_COMPILER="./Compiler_Main/Compiler_Main.exe"
FILE_VM="./VM_Main/VM_Main.exe"
DIR_SHARED="shared"
DIR_STD="Shared"
DIR_STDL="std_linux"
DIR_STDW="std_windows"
FILE_STDL="./Shared/std_linux/lib_m.so"
FILE_STDW="./Shared/std_windows/file.dll"
DIR_TST="tests"
FILE_TST="./Tests/*.lisp"   
if [ -d "$DIR_DEPLOY" ]; then
  # Take action if $DIR exists. #
  echo "Installing LispVM files in ${DIR_DEPLOY}..."
else
  mkdir $DIR_DEPLOY
fi
if [ -d "$DIR_LISPVM" ]; then
  # Take action if $DIR exists. #
  if [ "$(ls -A $DIR_LISPVM)" ]; then
     echo "Take action $DIR_LISPVM is not Empty"
     echo "Rename LispVM files in ${DIR_LISPVM}..."
     OLD=$(date +%Y-%m-%d_%s)
     mv ${DIR_LISPVM} ${DIR_LISPVM}_${OLD}
     echo "Create ${DIR_LISPVM}..."
     mkdir $DIR_LISPVM
  else
     echo "$DIR_LISPVM is Empty"
  fi
else
  echo "Create ${DIR_LISPVM}..."
  mkdir $DIR_LISPVM
fi
mkdir $DIR_BIN
if [ -f "$FILE_COMPILER" ]; then
    echo "$FILE_COMPILER exists."
else
    echo "Please make $FILE_COMPILER"
    exit -1
fi
if [ -f "$FILE_VM" ]; then
    echo "$FILE_VM exists."
else
    echo "Please make $FILE_VM"
    exit -2
fi
cp $FILE_COMPILER $DIR_BIN
cp $FILE_VM $DIR_BIN
#echo $OSTYPE
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        mkdir $DIR_LISPVM/$DIR_SHARED
        mkdir $DIR_LISPVM/$DIR_SHARED/$DIR_STDL
        cp $FILE_STDL $DIR_LISPVM/$DIR_SHARED/$DIR_STDL
elif [[ "$OSTYPE" == "darwin"* ]]; then
        # Mac OSX
        exit -4
elif [[ "$OSTYPE" == "cygwin" ]]; then
        # POSIX compatibility layer and Linux environment emulation for Windows
        exit -4
elif [[ "$OSTYPE" == "msys" ]]; then
        # Lightweight shell and GNU utilities compiled for Windows (part of MinGW)
        mkdir $DIR_LISPVM/$DIR_SHARED
        mkdir $DIR_LISPVM/$DIR_SHARED/$DIR_STDW
        cp $FILE_STDW $DIR_LISPVM/$DIR_SHARED/$DIR_STDW
elif [[ "$OSTYPE" == "win32" ]]; then
        # I'm not sure this can happen.
        exit -4
elif [[ "$OSTYPE" == "freebsd"* ]]; then
        # FreeBSD
        exit -4
else
        # Unknown.
        exit -4
fi
mkdir $DIR_LISPVM/$DIR_TST
cp $FILE_TST $DIR_LISPVM/$DIR_TST
exit 0
