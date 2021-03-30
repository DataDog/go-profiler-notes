//+build darwin

package main

import (
	"debug/macho"
	"errors"
	"fmt"
	"os"
)

func gopclntab() ([]byte, error) {
	file, err := macho.Open(os.Args[0])
	if err != nil {
		return nil, fmt.Errorf("elf.Open: %w", err)
	}
	for _, s := range file.Sections {
		if s.Name == "__gopclntab" {
			return s.Data()
		}
	}
	return nil, errors.New("could not find .gopclntab")
}
