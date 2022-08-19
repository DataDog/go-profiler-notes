//go:build darwin
// +build darwin

package main

import (
	"debug/gosym"
	"debug/macho"
	"os"
)

// from https://github.com/lizrice/debugger-from-scratch/blob/master/symbols.go
func goSymTable() (*gosym.Table, error) {
	exe, err := macho.Open(os.Args[0])
	if err != nil {
		return nil, nil
	}
	defer exe.Close()

	addr := exe.Section("__text").Addr

	lineTableData, err := exe.Section("__gopclntab").Data()
	if err != nil {
		return nil, nil
	}
	lineTable := gosym.NewLineTable(lineTableData, addr)
	if err != nil {
		return nil, nil
	}

	symTableData, err := exe.Section("__gosymtab").Data()
	if err != nil {
		return nil, nil
	}

	return gosym.NewTable(symTableData, lineTable)
}
