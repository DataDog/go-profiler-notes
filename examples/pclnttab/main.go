package main

import (
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
	symTable, err := goSymTable()
	if err != nil {
		return err
	}
	for _, pc := range callers() {
		file, line, fn := symTable.PCToLine(uint64(pc))
		fmt.Printf("%x: %s() %s:%d\n", pc, fn.Name, file, line)
	}
	return nil
}

func callers() []uintptr {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(2, pcs)
	pcs = pcs[0:n]
	return pcs
}
