package cmd

import "github.com/spf13/pflag"

func getNullableBool(flag string, flags *pflag.FlagSet) (*bool, error) {
	if flags.Changed(flag) {
		val, err := flags.GetBool(flag)
		if err != nil {
			return nil, err
		}
		return &val, nil
	}
	return nil, nil
}
