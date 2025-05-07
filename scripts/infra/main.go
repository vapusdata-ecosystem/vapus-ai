package main

import (
	"flag"
	"log"

	vapuspublish "github.com/vapusdata-ecosystem/vapusai/scripts/goscripts/publish"
)

var (
	forHelmpublish    = false
	forInstallerSetup = false
)

func init() {
	flag.BoolVar(&forHelmpublish, "helm-publish", false, "activate helm publish scripts")
	flag.BoolVar(&forInstallerSetup, "platform-installer", false, "activate platform installer setup scripts")
	flag.Parse()
}

func main() {
	if forHelmpublish {
		chartVersion := vapuspublish.HelmChartOps()
		log.Println(chartVersion)
	}
	if forInstallerSetup {
		log.Println("Installer setup")
	}
}
