package mtools

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/go-modulus/modulus/errors"
)

type CliErrorHandler struct {
}

func NewCliErrorHandler() *CliErrorHandler {
	return &CliErrorHandler{}
}

func (c *CliErrorHandler) HandleError(err error) {
	hint := errors.Hint(err)
	fmt.Println(color.RedString(hint))
	cause := errors.Cause(err)
	if cause != nil {
		fmt.Println("  Previous error:", color.RedString(cause.Error()))
	}
	meta := errors.Meta(err)
	if len(meta) > 0 {
		fmt.Println("  Error Info:")
		for k, v := range meta {
			fmt.Println("    ", color.YellowString(k), v)
		}
	}

	os.Exit(1)
}
