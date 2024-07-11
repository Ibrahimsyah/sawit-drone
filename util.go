package main

import "fmt"

type Util struct{}

func NewUtil() *Util {
	return &Util{}
}

func (u *Util) Scanln(target ...any) {
	fmt.Scanln(target...)
}
