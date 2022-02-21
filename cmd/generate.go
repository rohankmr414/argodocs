package cmd

import (
	"fmt"
	"github.com/junaidrahim/argodocs/mdgen"
	"github.com/junaidrahim/argodocs/workflow"
	"github.com/spf13/cobra"
	"os"
	"strings"
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
			dir, err := os.Getwd()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Generating docs from", dir)
			for _, arg := range args {
				templates, err := workflow.ParseFiles(arg)
				if err != nil {
					panic(err)
				}
				for _, wf := range templates {
					doc, err := mdgen.GetMdDoc(wf)
					if err != nil {
						panic(err)
					}

					wfPathList := strings.Split(wf.FilePath, "/")
					//wfFileName := wfPathList[len(wfPathList)-1]
					if len(wfPathList) > 0 {
						wfPathList = wfPathList[:len(wfPathList)-1]
					}
					wfPath := strings.Join(wfPathList, "/")
					if len(wfPath) > 0 {
						wfPath += "/"
					}
					err = os.MkdirAll(relativeOutputPrefix+"/"+outputPrefix+"/"+wfPath, os.ModePerm)
					if err != nil {
						panic(err)
					}
					fname := strings.Split(wf.FilePath, ".")
					if len(fname) > 0 {
						fname = fname[:len(fname)-1]
					}
					writeFilePath := strings.Join(fname, "")
					path := relativeOutputPrefix + "/" + outputPrefix + "/" + writeFilePath + ".md"
					err = doc.Export(path)
					if err != nil {
						panic(err)
					}
				}
			}
		},
	}

	generateCmd.Flags().StringVar(&outputPrefix, "output-prefix", "docs", "Output location prefix")
	generateCmd.Flags().StringVar(&relativeOutputPrefix, "relative-output-prefix", "./", "Relative output location prefix")

	return generateCmd
}
