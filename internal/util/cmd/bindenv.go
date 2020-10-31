package cmd

import (
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

// BindEnv binds environment variables to the flags.
// Env variable name is upper cased flag name and replaced "-"
// and "." into "_".
//
// Value set via flag has priority over value set via env variable.
func BindEnv(fs *pflag.FlagSet) error {
	set := make(map[string]bool)

	fs.Visit(func(f *pflag.Flag) {
		set[f.Name] = true
	})

	var flagError error

	fs.VisitAll(func(f *pflag.Flag) {
		if flagError != nil {
			return
		}

		replacer := strings.NewReplacer("-", "_", ".", "_")
		envVar := replacer.Replace(strings.ToUpper(f.Name))

		val := os.Getenv(envVar)
		if val == "" {
			return
		}

		if !set[f.Name] {
			t := f.Value.Type()
			if t == "stringArray" || t == "stringSlice" {
				vals := strings.Split(val, " ")
				for _, v := range vals {
					if err := fs.Set(f.Name, v); err != nil {
						flagError = errors.Wrapf(err, "wrapping %s with %v", f.Name, v)

						return
					}
				}

				return
			}

			if err := fs.Set(f.Name, val); err != nil {
				flagError = errors.Wrapf(err, "wrapping %s with %v", f.Name, val)

				return
			}
		}
	})

	return errors.WithStack(flagError)
}
