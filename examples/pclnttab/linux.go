//go:build linux
// +build linux

package main

import (
	"debug/elf"
	"errors"
	"fmt"
	"os"
)

// from https://github.com/lizrice/debugger-from-scratch/blob/master/symbols.go
func goSymTable() (*gosym.Table, error) {
	exe, err := elf.Open(os.Args[0])
	if err != nil {
		return nil, nil
	}
	defer exe.Close()

	addr := exe.Section(".text").Addr

	lineTableData, err := exe.Section(".gopclntab").Data()
	if err != nil {
		return nil, nil
	}
	lineTable := gosym.NewLineTable(lineTableData, addr)
	if err != nil {
		return nil, nil
	}

	symTableData, err := exe.Section(".gosymtab").Data()
	if err != nil {
		return nil, nil
	}

	return gosym.NewTable(symTableData, lineTable)
}
