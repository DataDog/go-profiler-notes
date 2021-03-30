package main

import (
	"debug/gosym"
	"fmt"
	"os"
	"runtime"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	data, err := gopclntab()
	if err != nil {
		return fmt.Errorf("gopclntab: %w", err)
	}
	return debugSymtab(data)
}

type StackTrace struct {
	Frames []StackFrame
}

type StackFrame struct {
	PC       int
	Function string
	File     string
	Line     int
}

func callers() []uintptr {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(2, pcs)
	pcs = pcs[0:n]
	return pcs
}

func debugSymtab(gopclntab []byte) error {
	table, err := gosym.NewTable(nil, gosym.NewLineTable(gopclntab, 0))
	if err != nil {
		return fmt.Errorf("gosym.NewTable: %w", err)
	}

	for _, pc := range callers() {
		file, line, fn := table.PCToLine(uint64(pc))
		fmt.Printf("%s() %s:%d\n", fn.Name, file, line)
	}
	return nil
}
