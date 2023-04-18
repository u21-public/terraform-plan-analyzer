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


			var reporter_type string
			if cCtx.Bool("github"){
				reporter_type = "github"
			} else {
				reporter_type = "basic"
			}

			reporter, err_reporter := PlanAnalyzer.NewReporter(reporter_type, report)
			if err_reporter != nil {
				fmt.Println(err_reporter)
				os.Exit(1)
			}
			err_report := reporter.PostReport()
			if err_report != nil {
				fmt.Println(err_report)
				os.Exit(1)
			}
			return nil
		},
	}
	if err_cli := app.Run(os.Args); err_cli != nil {
		log.Fatal(err_cli)
	}
}
