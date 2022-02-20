package cmd

import (
	"github.com/junaidrahim/argodocs/mdgen"
	"github.com/junaidrahim/argodocs/workflow"
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
			temp, err := workflow.ParseFiles(args[0])
			if err != nil {
				panic(err)
			}
			doc, err := mdgen.GetMdDoc(temp[0])
			doc.Export("./test.md")

		},
	}

	generateCmd.Flags().StringVar(&outputPrefix, "output-prefix", "docs", "Output location prefix")
	generateCmd.Flags().StringVar(&relativeOutputPrefix, "relative-output-prefix", "./", "Relative output location prefix")

	return generateCmd
}
