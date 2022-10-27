package recon

import (
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func SubFinder(fullUrl string, simpleUrl string) {
	path := strings.ReplaceAll(simpleUrl, ":", "")
	if CheckIfExists(path + "/Surface/Subfinder.txt") {
		return
	}

	// Subfinder command
	subFinder := exec.Command("subfinder", "-d", simpleUrl, "| sort -u")
	results, _ := subFinder.Output()

	file, _ := os.Create(path + "/Surface/Subfinder.txt")
	file.WriteString(string(results))
	file.Close()

	// Httpx command
	absPath, _ := filepath.Abs(path + "/Surface/Subfinder.txt")
	absPath = strings.ReplaceAll(absPath, "\\", "/")

	httpx := exec.Command("httpx", "-nc", "-status-code", "-title", "-tech-detect", "-list", absPath)
	results, _ = httpx.Output()

	file, _ = os.Create(path + "/Surface/Httpx.txt")
	file.WriteString(string(results))
	file.Close()
}

// Fetch some basic information about the target
func TargetInfo(fullUrl string, simpleUrl string) {
	path := strings.ReplaceAll(simpleUrl, ":", "")
	// Check if target info file exists, and if return
	if CheckIfExists(path + "/TargetInfo.txt") {
		return
	}

	// Create the info file
	infoFile, _ := os.Create(path + "/TargetInfo.txt")

	// First line
	infoFile.WriteString("---- Target " + simpleUrl + " ----\n")

	// Get the ip list from the target
	infoFile.WriteString("\n# Net Info\n")
	ips, _ := net.LookupIP(simpleUrl)
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			infoFile.WriteString("IPv4: " + ipv4.String() + "\n")
		}
	}

	// Get the cname
	cname, _ := net.LookupCNAME(simpleUrl)
	infoFile.WriteString("CNAME: " + cname + "\n")

	// Has script tag
	results := RegexFile(`<noscript>`, path+"/Target.html")
	// if we get a script tag
	if len(results) != 0 {
		infoFile.WriteString("\nTarget has a <noscript> tag\n")
	}

	// Close the info file
	defer infoFile.Close()

	// Get the robots.txt file from the target
	robotsFile, _ := os.Create(path + "/robots.txt")

	// Request the file
	r, er := http.Get(fullUrl + "robots.txt")
	if er == nil {
		// Write the robots.txt file
		robots, err := io.ReadAll(r.Body)
		if err != nil {
			println(err)
		} else {
			robotsFile.WriteString("\n" + string(robots))
		}

		defer robotsFile.Close()
	} else {
		println(er.Error())
	}
	defer r.Body.Close()

}
