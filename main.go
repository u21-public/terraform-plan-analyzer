package main

import (
	"fmt"
	"log"
	"os"

	cli "github.com/urfave/cli/v2"

	"github.com/u21-public/terraform-bulk-analyzer/internal/plananalyzer"
	"github.com/u21-public/terraform-bulk-analyzer/internal/reporter"
)

func main() {
	app := &cli.App{
		Name:  "Terraform Plan Analyzer",
		Usage: "Reads Plans -> Analyzes them -> prints report",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "tfplans",
				Usage:    "Relative path to folder holding tfplans",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "pretty",
				Usage: "Pretty prints to console",
			},
			&cli.BoolFlag{
				Name:  "github",
				Usage: "Posts report to github PR",
			},
			&cli.BoolFlag{
				Name:  "summary",
				Usage: "Prints machine-readable summary line to stdout (V8: only under this flag; never injected into PR comment body)",
			},
		},
		Action: func(cCtx *cli.Context) error {
			plans := plananalyzer.ReadPlans(cCtx.String("tfplans"))
			analyzedPlans := plananalyzer.NewPlanAnalyzer(plans)
			analyzedPlans.ProcessPlans()

			if cCtx.Bool("summary") {
				s := analyzedPlans.Summarize()
				fmt.Printf("::summary:: created=%d modified=%d destroyed=%d replaced=%d workspace_count=%d\n",
					s.Created, s.Modified, s.Destroyed, s.Replaced, s.WorkspaceCount)
			}

			report := analyzedPlans.GenerateReport()

			var reporterType string
			if cCtx.Bool("github") {
				reporterType = reporter.GithubReporterType
			} else {
				reporterType = reporter.BasicReporterType
			}

			reporter, err := reporter.NewReporter(reporterType, report)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if err = reporter.PostReport(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return nil
		},
	}
	if errCli := app.Run(os.Args); errCli != nil {
		log.Fatal(errCli)
	}
}
