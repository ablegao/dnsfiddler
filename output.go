package main

import (
	"bytes"
	"fmt"
	"time"
)

type Output struct {
	*bytes.Buffer
	uuid string
}

func (out *Output) Println(argv ...interface{}) {
	s := fmt.Sprintln(argv...)
	out.WriteString(" ---- " + out.uuid + time.Now().Format("15:03:04.0000 ") + s)

}

func (out *Output) NonPrefixPrintln(argv ...interface{}) {

	s := fmt.Sprintln(argv...)
	out.WriteString("|  " + s)
}
