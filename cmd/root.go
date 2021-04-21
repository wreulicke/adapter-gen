package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wreulicke/adapter-gen/generator"
	"golang.org/x/tools/imports"
)

func splitPackageAndTarget(v string) (string, string, bool) {
	i := strings.LastIndex(v, ".")

	packageName := v[:i]
	targetName := v[i+1:]
	var shouldBePointer bool
	if trimed := strings.TrimPrefix(packageName, "*"); packageName != trimed {
		packageName = trimed
		shouldBePointer = true
	}
	return packageName, targetName, shouldBePointer
}

func NewRootCommand() *cobra.Command {
	var output string
	cmd := &cobra.Command{
		Use:   "adapter-gen",
		Short: "adapter-gen is code generator for adapter",
		RunE: func(cmd *cobra.Command, args []string) error {
			packageName, targetName, shouldBePointer := splitPackageAndTarget(args[0])
			b := &bytes.Buffer{}
			gen, err := generator.New(b, packageName, targetName, shouldBePointer)
			if err != nil {
				return err
			}
			if err := gen.Generate(); err != nil {
				return err
			}
			bs, err := imports.Process(output, b.Bytes(), &imports.Options{})
			if err != nil {
				return err
			}
			_ = os.MkdirAll(filepath.Dir(output), os.ModePerm)
			f, err := os.Create(output)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(f, bytes.NewBuffer(bs))
			return err
		},
		Args: cobra.ExactArgs(1),
	}
	cmd.Flags().StringVarP(&output, "output", "o", "adapter/adapter.go", "output name")
	return cmd
}

func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
