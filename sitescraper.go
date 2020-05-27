package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Weird file permissions is because of a shared folder inside a VM
	// TODO: user could specify the directory name, everything is saved here
	os.Mkdir("SCRAPE", os.FileMode(0770))

	url := "http://10.10.10.187"
	fmt.Println("Connecting:", url)
	// We'll save the front page first
	saveFrontPage(url)

	// New request to parse the contents, as the response io.ReadCloser
	// is not reusable.
	fmt.Println("Making 2nd request to find out file references.")
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	// We'll have all the text elements in the 2nd returned element
	// TODO: remove empty and meaningless elements
	var values, _ = parse(response)
	var filenames []string

	// If there is a dot, it is likely a file or a link that can be followed
	// Other cases relevant for site scraping?
	for i := 0; i < len(values); i++ {
		if strings.Contains(values[i], ".") {
			filenames = append(filenames, values[i])
		}
	}
	// TODO: Inform the user for the number of files
	// Maybe ask to list the files found? Sizes?
	downloadSiteFiles(url, filenames)
	return
}

func parse(response *http.Response) ([]string, []string) {

	tok := html.NewTokenizer(response.Body)
	var values []string
	var textTokens []string

	for {
		tokenType := tok.Next()
		if tokenType == html.ErrorToken {
			// TODO: make the ending explicit on io.EOF
			fmt.Println("Reached end of request.")
			return values, textTokens
		}

		// <script> and <noscript> are  escaped: '<' becomes '&lt;'. These
		// cases along with other items such as newlines and tabs are here.
		// <noscript> at least might contain html elements inside.
		if tokenType == html.TextToken {
			textTokens = append(textTokens, string(tok.Text()))
			continue
		}

		// Get the key=value pairs out of every html tag
		_, att := tok.TagName() // Tag names ignored
		if att == true {
			for {
				_, value, more := tok.TagAttr() // keys are not used for now
				values = append(values, string(value))
				if more == false {
					break
				}
			}
		}
	}
}

func downloadSiteFiles(url string, filenames []string) {

	// Create a directory structure?
	for i := range filenames {
		fullurl := url + "/" + filenames[i]
		fmt.Println("Downloading:", fullurl)
		// Replacing slashes with underscores...
		filename := strings.ReplaceAll(filenames[i], "/", "_")
		getFile(fullurl, filename)
	}
}

func getFile(url string, filename string) {

	// TODO: error handling
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	out, err := os.Create("SCRAPE/" + filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer out.Close()

	bytes, err := io.Copy(out, resp.Body)
	fmt.Println(strconv.FormatInt(bytes, 10) + " bytes --> SCRAPE/" + filename)
	return
}

func saveFrontPage(url string) {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	out, err := os.Create("SCRAPE/ROOT.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer out.Close()

	bytes, err := io.Copy(out, resp.Body)
	fmt.Println("Saving front page html from:", url)
	fmt.Println(strconv.FormatInt(bytes, 10) + " bytes --> SCRAPE/ROOT.html")

	return
}
