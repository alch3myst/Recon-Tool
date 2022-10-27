package main

import (
	"flag"
	recon "recon/src"
)

// TODO: headers
// TODO: optimize fot multiple inputs

func main() {
	// Flags

	// Target url
	var target string
	flag.StringVar(&target, "t", "", "Target")

	// Mode flag
	var mode string
	flag.StringVar(&mode, "m", "", "Select the mode to run")

	// Set if page is http
	var isHttp bool
	flag.BoolVar(&isHttp, "http", false, "Select the mode to run")

	// Parse flags
	flag.Parse()

	recon.UiBanner()

	// If has no  target print the help message
	if target == "" {
		recon.UiHelp()
		return
	}

	// Parse the target url
	fullUrl, simpleUrl := recon.ParseTargetURL(target, isHttp)

	switch mode {
	case "recon":
		Recon(fullUrl, simpleUrl)
	default:
		recon.UiHelp()
	}
}

func Recon(fullUrl string, simpleUrl string) {

	// Base content -------
	// Setup the target, creating folders and stuff
	recon.SetupTarget(fullUrl, simpleUrl)

	// Get basic information about the target
	recon.TargetInfo(fullUrl, simpleUrl)

	// Get the page
	recon.GetPage(fullUrl, simpleUrl)
	// !Base content -------

	// Scrapper and discover
	// Scrap information from the target
	recon.Scrapper(fullUrl, simpleUrl)

	// Build a basic attack surface
	recon.SubFinder(fullUrl, simpleUrl)

	recon.DownloadJS(fullUrl, simpleUrl)
	// !Scrapper and discover
}
