#!/bin/bash

run_test_aux() {
  echo Passing $1...
  REF=$1
  GOFILE=${REF%%.ref}.go
  EXE=${REF%%.ref}
  SATELLITE=${REF%%.ref}.SATELLITE.ref
  SATELLITEGO=${REF%%.ref}.SATELLITE.go

  GOLINE=go
  
  if [ -e $SATELLITE ]; then
    R05CCOMP= R05PATH= ../build/refal5t $REF $SATELLITE 2>__error.txt
    if [ $? -ge 200 ]; then
      echo COMPILER ON $REF FAILS, SEE __error.txt
      exit
    fi
  else
    R05CCOMP= R05PATH= ../build/refal5t $REF 2>__error.txt
    if [ $? -ge 200 ]; then
      echo COMPILER ON $REF FAILS, SEE __error.txt
      exit
    fi
  fi

  # R05CCOMP= R05PATH= ../build/refal5t $REF $SATELLITE 2>__error.txt
  # if [ $? -ge 200 ]; then
    # echo COMPILER ON $REF FAILS, SEE __error.txt
    # exit
  # fi
  rm __error.txt
  if [ ! -e $GOFILE ]; then
    echo COMPILATION FAILED
    exit
  fi

  #   if [ ! -e $SATELLITEC ]; then
  #     echo COMPILATION FAILED
  #     exit
  #   fi
  # else
  #   SATELLITEC=

  $GOLINE build -o $EXE $GOFILE
  if [ $? -gt 0 ]; then
    echo COMPILATION FAILED
    exit
  fi

  ./$EXE 2> __dump.txt
  # команда [ в условии меняет код возврата, поэтому нужна переменная
  EXIT_CODE=$?
  if [ $EXIT_CODE -ge 200 ]; then
    echo TEST FAILED, SEE __dump.txt
    exit
  elif [ $EXIT_CODE -gt 0 ]; then
    echo "TEST FAILED (INTERNAL ERROR)"
    exit
  fi

  rm $GOFILE $EXE $SATELLITEC
  [ -e __dump.txt ] && rm __dump.txt

  echo $1 "Ok!"
  echo
}

run_test_aux.BAD-SYNTAX() {
  echo Passing $1...
  REF=$1
  GOFILE=${REF%%.ref}.go
  EXE=${REF%%.ref}

  R5TCCOMP= R5TPATH= ../build/refal5t $REF 2>__error.txt
  if [ $? -ge 200 ]; then
    echo COMPILER ON $REF FAILS, SEE __error.txt
    exit
  fi
  rm __error.txt
  if [ -e $GOFILE ]; then
    echo COMPILATION SUCCESSED, BUT EXPECTED SYNTAX ERROR
    rm $GOFILE
    exit
  fi

  echo "Ok! Compiler didn't crash on invalid syntax"
  echo
}

run_test_aux.SATELLITE() {
  :
}

run_test() {
  REF=$1
  SUFFIX=`echo ${REF%%.ref} | sed 's/[^.]*\(\.[^.]*\)*/\1/'`
  run_test_aux$SUFFIX $1
}

# source ../c-plus-plus.conf.sh

if [ -z "$1" ]; then
  for s in *.ref; do
    run_test $s
  done
else
  for s in $*; do
    run_test $s
  done
fi
