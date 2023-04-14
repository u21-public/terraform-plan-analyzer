package main

import (
	"fmt"
	cli "github.com/urfave/cli/v2"
	"log"
	"os"

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
			fmt.Println(report)
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
