package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewGenerateCommand() *cobra.Command {
	var (
		outputPrefix         string
		relativeOutputPrefix string
	)
	// generateCmd represents the generate command
	var generateCmd = &cobra.Command{
		Use:   "generate PATH --output-prefix=PREFIX --relative-output-prefix=PREFIX",
		Short: "Generate docs from workflow manifest.",
		Long:  `Generate reference docs from argo workflows.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(args)
			fmt.Println(outputPrefix)
		},
	}

	generateCmd.Flags().StringVar(&outputPrefix, "output-prefix", "docs", "Output location prefix")
	generateCmd.Flags().StringVar(&relativeOutputPrefix, "relative-output-prefix", "./", "Relative output location prefix")

	return generateCmd
}
