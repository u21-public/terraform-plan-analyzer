package main

import (
	"fmt"
	"log"
	"os"

	cli "github.com/urfave/cli/v2"

	"github.com/u21-public/terraform-bulk-analyzer/internal/PlanAnalyzer"
)

func main() {
	app := &cli.App{
		Name:  "Terraform Bulk Analyzer",
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
				reporterType = "github"
			} else {
				reporterType = "basic"
			}

			reporter, err_reporter := PlanAnalyzer.NewReporter(reporterType, report)
			if errReporter != nil {
				fmt.Println(errReporter)
				os.Exit(1)
			}
			errReport := reporter.PostReport()
			if errReport != nil {
				fmt.Println(errReport)
				os.Exit(1)
			}
			return nil
		},
	}
	if errCli := app.Run(os.Args); errCli != nil {
		log.Fatal(errCli)
	}
}
