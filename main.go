package main

import (
	"fmt"
	"strconv"
	"strings"
)

// a typical Forth implements floating point values using
// a completely separate stack
type expression struct {
	valueN number
	valueW word
	kind   kind
}

type kind uint8

const (
	numberKind kind = iota
	wordKind
)

type number int

type word string

func (e expression) String() string {
	if e.kind == wordKind {
		return string(e.valueW)
	}
	return fmt.Sprintf("%d", e.valueN)
}

type env map[word]function

type function struct {
	isBuiltin bool
	body      []expression
}

var builtIns = map[word]func(*Forth){
	// arithmetic
	"+":      func(f *Forth) { f.pushN(f.popN() + f.popN()) },
	"-":      func(f *Forth) { a, b := f.popN(), f.popN(); f.pushN(b - a) },
	"*":      func(f *Forth) { f.pushN(f.popN() * f.popN()) },
	"/":      func(f *Forth) { a, b := f.popN(), f.popN(); f.pushN(b / a) },
	"mod":    func(f *Forth) { a, b := f.popN(), f.popN(); f.pushN(b % a) },
	"negate": func(f *Forth) { f.pushN(-f.popN()) },
	"abs": func(f *Forth) {
		if x := f.popN(); x < 0 {
			f.pushN(-x)
		} else {
			f.pushN(x)
		}
	},
	"max": func(f *Forth) { f.pushN(max(f.popN(), f.popN())) },
	"min": func(f *Forth) { f.pushN(min(f.popN(), f.popN())) },

	// stack manipulation
	"dup":        func(f *Forth) { x := f.pop(); f.push(x); f.push(x) },
	"swap":       func(f *Forth) { a, b := f.pop(), f.pop(); f.push(a); f.push(b) },
	"rot":        func(f *Forth) { a, b, c := f.pop(), f.pop(), f.pop(); f.push(b); f.push(a); f.push(c) },
	"drop":       func(f *Forth) { f.pop() },
	"nip":        func(f *Forth) { x := f.pop(); f.pop(); f.push(x) },
	"clearstack": func(f *Forth) { f.stack = nil },

	// advanced stack manipulation
	"tuck": func(f *Forth) { a, b := f.pop(), f.pop(); f.push(a); f.push(b); f.push(a) },
	"over": func(f *Forth) { a, b := f.pop(), f.pop(); f.push(b); f.push(a); f.push(b) },
	"roll": func(f *Forth) { i := f.popN(); x := f.cut(int(i)); f.push(x) },
	"pick": func(f *Forth) { i := f.popN(); x := f.peekAt(int(i)); f.push(x) },

	// i/o
	"CR": func(f *Forth) { fmt.Println() },
	".":  func(f *Forth) { fmt.Print(f.pop()) },
	".s": func(f *Forth) {
		fmt.Printf("<%d>", len(f.stack))
		for i := 0; i < len(f.stack); i++ {
			fmt.Print(" ", f.stack[i])
		}
	},
	"emit": func(f *Forth) { fmt.Printf("%c", f.popN()) },
	"see":  func(f *Forth) { f.parameterFn = func(arg expression) {
		w := arg.valueW
		fmt.Printf(": %s", w)
		for _, e := range f.env[word(w)].body {
			fmt.Print(" ", e)
		}
		fmt.Print(" ;")
	}},
	"char": func(f *Forth) { f.parameterFn = func(arg expression) {
		f.push(expression{valueN: number(arg.valueW[0]), kind: numberKind})
	}},

	// compile mode
	":": func(f *Forth) {
		if f.mode == compile {
			panic("already in compile mode")
		}
		f.mode = compile
	},
	";": func(f *Forth) {
		if len(f.compileBuffer) < 2 {
			panic("invalid word declaration")
		}
		w := f.compileBuffer[0]
		if w.kind != wordKind {
			panic("cannot redeclare a number")
		}
		decl := make([]expression, len(f.compileBuffer)-1)
		copy(decl, f.compileBuffer[1:])
		f.env[w.valueW] = function{body: decl}
		f.compileBuffer = nil
		f.mode = immediate
	},
	"]": func(f *Forth) { f.mode = compile },
}

func parseLine(input string) []expression {
	var out []expression
	var commentMode bool
	for _, s := range strings.Fields(input) {
		if s == "\\" {
			break
		}
		if s == "(" {
			commentMode = true
			continue
		}
		if commentMode {
			if s == ")" {
				commentMode = false
			}
			continue
		}
		n, err := strconv.Atoi(s)
		if err == nil {
			out = append(out, expression{valueN: number(n), kind: numberKind})
			continue
		}
		out = append(out, expression{valueW: word(s), kind: wordKind})
	}
	return out
}

func main() {
	fmt.Println("GO FORTH, AND PROSPER")
	f := NewForth()
	input := "25 10 * 50 + CR ."
	e := parseLine(input)
	f.run(e)
}
