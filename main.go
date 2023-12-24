package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Operator interface {
	Operate(string)
}

type Parameters struct {
	Operation string
	FieldsArr []int
	Delimeter string
	Source    io.ReadCloser
}

func NewParamenters() Parameters {
	return Parameters{}
}

func main() {
	pars := NewParamenters()
	fmt.Println("Start")
	last := false
	for i := 1; i < len(os.Args); i++ {
		if strings.Contains(os.Args[i], "-f") {
			str := strings.Replace(os.Args[i], "-f", "", -1)
			fieldsArr := []int{}
			if strings.Contains(str, "\"") {
				for _, v := range strings.Split(strings.Replace(str, "\"", "", -1), " ") {
					fieldId, err := strconv.Atoi(v)
					if err != nil {
						panic(err)
					}
					fieldsArr = append(fieldsArr, fieldId)
				}
			} else {
				for _, v := range strings.Split(str, ",") {
					fieldId, err := strconv.Atoi(v)
					if err != nil {
						panic(err)
					}
					fieldsArr = append(fieldsArr, fieldId)
				}
				pars.FieldsArr = fieldsArr
			}

			pars.Operation = "f"
			last = i == len(os.Args)-1
		}
		if strings.Contains(os.Args[i], "-d") {
			str := strings.Replace(os.Args[i], "-d", "", -1)
			pars.Delimeter = str
			last = i == len(os.Args)-1
		}
		if i == len(os.Args)-1 && !last {
			if os.Args[i] == "-" {
				pars.Source = os.Stdin
			} else {
				f, err := os.Open(os.Args[i])
				if err != nil {
					panic(err)
				}
				pars.Source = f
			}
		}
	}

	if pars.Delimeter == "" {
		pars.Delimeter = "\t"
	}

	if pars.Operation == "" {
		panic("no operation")
	}

	if len(pars.FieldsArr) == 0 {
		panic("invalid fields number")
	}

	if pars.Source == nil {
		pars.Source = os.Stdin
	}

	var oper Operator
	if pars.Operation == "f" {
		oper = NewFieldOPerator(pars.Delimeter, pars.FieldsArr, os.Stdout)
	}

	b := bufio.NewScanner(pars.Source)
	for b.Scan() {
		row := b.Text()
		oper.Operate(row)
	}

	pars.Source.Close()
}

type FieldOperation struct {
	delimeter string
	fields    []int
	w         io.WriteCloser
}

func NewFieldOPerator(d string, f []int, w io.WriteCloser) *FieldOperation {
	return &FieldOperation{
		delimeter: d,
		fields:    f,
		w:         w,
	}
}

func (fo *FieldOperation) SetDelimiter(d string) {
	fo.delimeter = d
}

func (fo *FieldOperation) Operate(row string) {
	arr := strings.Split(row, fo.delimeter)
	for _, v := range fo.fields {
		fo.w.Write([]byte(arr[v-1] + fo.delimeter))
	}

	fo.w.Write([]byte("\n"))
}
