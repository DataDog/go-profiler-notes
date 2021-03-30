//+build linux

package main

import (
	"debug/elf"
	"errors"
	"fmt"
	"os"
)

func gopclntab() ([]byte, error) {
	file, err := elf.Open(os.Args[0])
	if err != nil {
		return nil, fmt.Errorf("elf.Open: %w", err)
	}
	for _, s := range file.Sections {
		if s.Name == ".gopclntab" {
			return s.Data()
		}
	}
	return nil, errors.New("could not find .gopclntab")
}
