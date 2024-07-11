package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	osExit = func(code int) {}
	assert.NotPanics(t, main)
}
