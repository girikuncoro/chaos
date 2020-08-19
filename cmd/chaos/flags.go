package main

import (
	"github.com/girikuncoro/chaos/pkg/cli/output"
	"github.com/spf13/cobra"
)

const outputFlag = "output"

// bindOutputFlag will add output flag to given command
func bindOutputFlag(cmd *cobra.Command, varRef *output.Format) {
	cmd.Flags().VarP(newOutputValue(output.Table, varRef), outputFlag, "o",
		"prints output in the specified format. Allowed values: table")
}

type outputValue output.Format

func newOutputValue(defaultValue output.Format, p *output.Format) *outputValue {
	*p = defaultValue
	return (*outputValue)(p)
}

func (o *outputValue) Type() string {
	return "format"
}

func (o *outputValue) String() string {
	return string(*o)
}

func (o *outputValue) Set(s string) error {
	outfmt, err := output.ParseFormat(s)
	if err != nil {
		return err
	}
	*o = outputValue(outfmt)
	return nil
}
