package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// https://learnxinyminutes.com/forth/
func TestLearnXinYmins1(t *testing.T) {

	for i, tt := range []struct {
		input string
		want  string
	}{
		{input: "\\ This is a comment"},
		{input: "( This is also a comment but it's only used when defining words )"},
		{
			input: "5 2 3 56 76 23 65 .s",
			want:  "<7> 5 2 3 56 76 23 65",
		},
		// arithmetic
		{input: "5 4 + .", want: "9"},
		{input: "6 7 * .", want: "42"},
		{input: "1360 23 - .", want: "1337"},
		{input: "12 12 / .", want: "1"},
		{input: "3 2 mod .", want: "1"},
		{input: "99 negate .", want: "-99"},
		{input: "-99 abs .", want: "99"},
		{input: "52 23 max .", want: "52"},
		{input: "52 23 min .", want: "23"},
		// stack manipulation
		{input: "3 dup - .s", want: "<1> 0"},
		{input: "2 5 swap / .s", want: "<1> 2"},
		{input: "6 4 5 rot .s", want: "<3> 4 5 6"},
		{input: "4 0 drop 2 / .s", want: "<1> 2"},
		{input: "1 2 3 nip .s", want: "<2> 1 3"},
		// advanced stack manipulation
		// note: index is 0-based from top of stack down
		{input: "1 2 3 4 tuck .s", want: "<5> 1 2 4 3 4"},
		{input: "1 2 3 4 over .s", want: "<5> 1 2 3 4 3"},
		{input: "1 2 3 4 2 roll .s", want: "<4> 1 3 4 2"},
		{input: "1 2 3 4 2 pick .s", want: "<5> 1 2 3 4 2"},
	} {
		r, w, _ := os.Pipe()
		os.Stdout = w

		f := NewForth()
		f.run(parseLine(tt.input))

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		got := buf.String()

		if got != tt.want {
			t.Errorf("test %d: got %q, want %q", i, got, tt.want)
		}
	}
}

func TestLearnXinYmins2(t *testing.T) {
	f := NewForth()

	for i, tt := range []struct {
		input string
		want  string
	}{
		// creating words
		{input: ": square ( n -- n ) dup * ;"},
		{input: "5 square .", want: "25"},
		{input: "see square", want: ": square dup * ;"},
	} {
		r, w, _ := os.Pipe()
		os.Stdout = w

		f.run(parseLine(tt.input))

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		got := buf.String()

		if got != tt.want {
			t.Errorf("test %d: got %q, want %q", i, got, tt.want)
		}
	}
}
