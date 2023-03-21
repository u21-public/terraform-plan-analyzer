package main

import (
	"fmt"
	"github.com/mattthaber/terraform-bulk-analyzer/internal/PlanAnalyzer"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "Terraform Bulk Analyzer",
		Usage: "Reads Plans -> Analayzes them -> prints report",
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

			reporter,err := PlanAnalyzer.NewReporter(cCtx.Bool("github"), report)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			reporter.PostReport()

			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
