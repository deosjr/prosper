package main

import (
	"fmt"
	"strings"
)

type Forth struct {
	stack         stack
	env           env
	mode          mode
	compileBuffer []expression
}

func NewForth() *Forth {
	environment := env{}
	for w := range builtIns {
		environment[w] = function{isBuiltin: true}
	}
	return &Forth{
		stack: stack{},
		env:   environment,
	}
}

type stack []expression

func (f *Forth) push(e expression) {
	f.stack = append(f.stack, e)
}

func (f *Forth) pushN(n number) {
	f.push(expression{valueN: n, kind: numberKind})
}

func (f *Forth) pop() expression {
	n := len(f.stack)
	if n == 0 {
		panic("pop on empty stack")
	}
	e := f.stack[n-1]
	f.stack = f.stack[:n-1]
	return e
}

func (f *Forth) popN() number {
	e := f.pop()
	if e.kind != numberKind {
		panic("popN expected a number")
	}
	return e.valueN
}

func (f *Forth) popW() word {
	e := f.pop()
	if e.kind != wordKind {
		panic("popW expected a word")
	}
	return e.valueW
}

func (f *Forth) peek() expression {
	n := len(f.stack)
	if n == 0 {
		panic("peek on empty stack")
	}
	return f.stack[n-1]
}

func (f *Forth) peekAt(index int) expression {
	n := len(f.stack)
	if n <= index {
		panic("peekAt beyond stack")
	}
	return f.stack[n-index-1]
}

func (f *Forth) cut(index int) expression {
	n := len(f.stack)
	if n <= index {
		panic("cut beyond stack")
	}
	i := n - index - 1
	e := f.stack[i]
	f.stack = append(f.stack[:i], f.stack[i+1:]...)
	return e
}

type mode uint8

const (
	immediate mode = iota
	compile
)

func (f *Forth) run(program []expression) {
	for _, e := range program {
		switch f.mode {
		case immediate:
			f.eval(e)
		case compile:
			// todo: separate dict for immediate words?
			if e.valueW == ";" {
				f.eval(e)
				continue
			}
			if e.valueW == "[" {
				f.mode = immediate
				continue
			}
			if e.valueW == "literal" {
				f.compileBuffer = append(f.compileBuffer, f.pop())
				continue
			}
			f.compileBuffer = append(f.compileBuffer, e)
		}
	}
}

func (f *Forth) eval(e expression) {
	if e.kind == numberKind {
		f.push(e)
		return
	}
	if strings.Contains(string(e.valueW), " ") {
		// this feels like a hack to support 'see w' only atm
		fields := strings.Fields(string(e.valueW))
		cmd, arg := fields[0], fields[1]
		switch cmd {
		case "see":
			fmt.Printf(": %s", arg)
			for _, e := range f.env[word(arg)].body {
				fmt.Print(" ", e)
			}
			fmt.Print(" ;")
		case "char":
			f.push(expression{valueN: number(arg[0]), kind: numberKind})
		default:
			panic("lookahead hack abused!")
		}
		return
	}
	fn, ok := f.env[e.valueW]
	if !ok {
		panic(fmt.Sprintf("unknown word %s", e.valueW))
	}
	if fn.isBuiltin {
		builtIns[e.valueW](f)
		return
	}
	for _, e := range fn.body {
		f.eval(e)
	}
}
