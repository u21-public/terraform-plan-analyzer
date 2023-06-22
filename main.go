package main

import (
	"fmt"
	"log"
	"os"

	cli "github.com/urfave/cli/v2"

	"github.com/u21-public/terraform-bulk-analyzer/internal/PlanAnalyzer"
	"github.com/u21-public/terraform-bulk-analyzer/internal/Reporter"
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
		},
		Action: func(cCtx *cli.Context) error {
			plans := PlanAnalyzer.ReadPlans(cCtx.String("tfplans"))
			analyzedPlans := PlanAnalyzer.NewPlanAnalyzer(plans)
			analyzedPlans.ProcessPlans()
			report := analyzedPlans.GenerateReport()

			var reporterType string
			if cCtx.Bool("github") {
				reporterType = Reporter.GithubReporterType
			} else {
				reporterType = Reporter.BasicReporterType
			}

			reporter, err := Reporter.NewReporter(reporterType, report)
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
