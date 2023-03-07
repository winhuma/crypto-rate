package myfunc

import (
	"fmt"
	"runtime"
	"strings"
)

func MyErrFormat(err error) error {
	pc, file, line, _ := runtime.Caller(1)
	funcName := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	return fmt.Errorf("[%s][%s] %v", funcName[len(funcName)-1], file+":"+fmt.Sprint(line), err)
}
