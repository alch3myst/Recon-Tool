package recon

import (
	"bufio"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/go-rod/rod"
)

func Scrapper(fullUrl string, simpleUrl string) {
	println("\nStarting Scrapper")
	// Href, src ping stuff that has some kind of link/content to see
	println("Searching for links and sources...")
	results := RegexFile(`(href|src|ping)="[a-zA-Z/0-9.:?=\-&;%_,]*.`, simpleUrl+"/Target.html")
	WriteRegexToFile(results, simpleUrl+"/Scraps/href_src.txt", false)

	println("Searching for js links on page source...")
	results = RegexFile(`[\S]*\.js"`, simpleUrl+"/Target.html")
	WriteRegexToFile(results, simpleUrl+"/Scraps/anyJS.txt", false)

	println("Searching for any link on page source...")
	results = RegexFile(`(http|https)://[\S]*('|")`, simpleUrl+"/Target.html")
	WriteRegexToFile(results, simpleUrl+"/Scraps/links.txt", true)

	println("Target source scrapped.\n")
}

// , wg *sync.WaitGroup
func DownloadJS(fullUrl string, simpleUrl string) {
	path := strings.ReplaceAll(simpleUrl, ":", "")

	println("\nDownloading scripts found (from src in anyJS.txt)")

	// Open the js file list
	jsList, _ := os.Open(path + "/Scraps/anyJS.txt")

	// Create buffer to read line by line
	scan := bufio.NewScanner(jsList)

	// Final url list
	var list []string

	// Break the js file list line by line
	for scan.Scan() {
		// Get current line
		tmpLine := scan.Text()

		// Create url from src=""
		if strings.Contains(tmpLine, "src=") {
			// Remove the src tag
			tmpLine = strings.Replace(tmpLine, "src=\"", "", -1)

			// Create the url
			list = append(list, fullUrl+tmpLine[:len(tmpLine)-1])
		}
	}

	// Download js file
	for _, jsLink := range list {
		// regex the js file name
		fn := regexp.MustCompile(`(sj\.)[a-z0-9-_.]*`)

		// Create the fine name string
		fileName := reverse(string(fn.Find([]byte(reverse(jsLink)))))

		// Request the data
		jsData := RequestUrl(jsLink)

		// Create and write the file
		f, _ := os.Create(path + "/JS/" + fileName)
		f.Write(jsData)
		f.Close()
	}

	list = nil

	println("Download finished.")
}

func reverse(s string) string {
	rns := []rune(s)
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}
	return string(rns)
}

func GetPage(fullUrl string, simpleURL string) {
	path := strings.ReplaceAll(simpleURL, ":", "")

	_, err := os.Stat(simpleURL + "/Target.html")
	if !os.IsNotExist(err) {
		return
	}

	println("Getting the target page...")
	// Create the headless browser instance
	headlessBrowser := rod.New().MustConnect()

	// Headless browser to get the page content
	page := headlessBrowser.MustPage(fullUrl)

	// Print the page
	// this awaits for the page js to load
	println("Taking a page screenshot...")
	page.MustWaitLoad().MustScreenshot(path + "/Screenshot.png")

	// Get the source code after the js compile
	pageString, _ := page.HTML()

	// Write the page html source to a file
	os.WriteFile(path+"/Target.html", []byte(pageString), 0644)
}

// Parse the target url
func ParseTargetURL(target string, isHttp bool) (string, string) {
	// Convert the simple url to a complete url for the headless browser
	var httpsTargetRawURL string = target
	if !strings.Contains(target, "https://") && !strings.Contains(target, "http://") {
		// Only https for now.

		if isHttp {
			httpsTargetRawURL = "http://www." + target + "/"
		} else {
			httpsTargetRawURL = "https://www." + target + "/"
		}

		// Remove / if have
		simpleURL := strings.ReplaceAll(target, "/", "")

		return httpsTargetRawURL, simpleURL
	}

	// Remove all the / form the url
	simpleURL, _ := url.Parse(target)

	return httpsTargetRawURL, simpleURL.Hostname()
}

// Regex a file and return the found results
func RegexFile(pattern string, filePath string) [][]byte {
	file, _ := os.ReadFile(filePath)
	reg := regexp.MustCompile(pattern)
	match := reg.FindAll(file, -1)

	return match
}

// Write regex matches to file
func WriteRegexToFile(results [][]byte, path string, removeLast bool) {
	file, _ := os.Create(path)

	// Remove last character from string before writing to the file
	if removeLast {
		for _, v := range results {
			_, _ = file.WriteString(string(v)[:len(v)-1] + "\n")
		}
		file.Close()

		return
	}

	// Write the output file
	for _, v := range results {
		_, _ = file.WriteString(string(v) + "\n")
	}
	file.Close()
}
