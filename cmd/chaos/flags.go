package main

import (
	"github.com/girikuncoro/chaos/pkg/cli/output"
	"github.com/girikuncoro/chaos/pkg/cli/values"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const outputFlag = "output"

func addValueOptionsFlags(f *pflag.FlagSet, v *values.Options) {
	f.StringArrayVar(&v.Values, "set", []string{}, "set values on the command line")
}

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
