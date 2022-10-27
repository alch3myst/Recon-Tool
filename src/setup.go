package recon

import (
	"os"
	"strings"
)

// Create folder structure
func SetupTarget(fullUrl string, simpleURL string) {
	path := strings.ReplaceAll(simpleURL, ":", "")

	// Main folder
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		println("\nCreating folder structure at " + path)
		os.Mkdir(path, 0777)
	}

	// Scraps folder
	_, err = os.Stat(path + "/Scraps")
	if os.IsNotExist(err) {
		os.Mkdir(path+"/Scraps", 0777)
	}

	// Scraps folder
	_, err = os.Stat(path + "/Surface")
	if os.IsNotExist(err) {
		os.Mkdir(path+"/Surface", 0777)
	}

	// JS folder
	_, err = os.Stat(path + "/JS")
	if os.IsNotExist(err) {
		os.Mkdir(path+"/JS", 0777)
	}
}
