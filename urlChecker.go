package main

import (
	//"flag"

	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/mb0/glob"
	flag "github.com/spf13/pflag"
)

func main() {
	globFlag := flag.Bool("g", false, "Glob pattern")

	//add flags of -j, --jason to enable the program to output JSON
	jflag := flag.BoolP("json", "j", false, "json output")

	// checking version flag
	vflag := flag.BoolP("version", "v", false, "version check")

	// ignore url flag
	ignoreFlag := flag.BoolP("ignore", "i", false, "ignore url patterns")

	telescopeFlag := flag.BoolP("telescope", "t", false, "check from telescope")

	flag.Parse()
	//deal with non-file path, giving usage message
	if len(os.Args) == 1 {
		fmt.Println("help/usage message: To run this program, please pass an argument to it,i.e.: go run urlChecker.go urls.txt")

	} else {
		//feature of checking version
		if *vflag {
			fmt.Println("  *****  urlChecker Version 0.2  *****  ")
			return
		}

		if *globFlag {
			//Assign the glob pattern provided to a local variable
			pattern := flag.Args()[0]
			//Read all files in the current directory
			files, _ := ioutil.ReadDir(".")
			//Create a globber object
			globber, _ := glob.New(glob.Default())
			//Loop through all files
			for _, file := range files {
				//Check if the file name match the glob pattern provided
				matched, _ := globber.Match(pattern, file.Name())
				//If matched then run the url check on that file
				if matched {
					//open file and read it
					content, err := ioutil.ReadFile(file.Name())
					if err != nil {
						log.Fatal(err)
					}
					textContent := string(content)

					fmt.Println(">>  ***** UrlChecker is working now...... *****  <<")
					fmt.Println("--------------------------------------------------------------------------------------------------")
					//call functions to check the availability of each url
					checkURL(extractURL(textContent))
				}
			}
			return
		}

		if *ignoreFlag {
			combineStrIgnorePatterns := ""
			ignoreFilePath := flag.Args()[0]

			ignoreList := parseIgnoreURL(ignoreFilePath)

			filepath := flag.Arg(1)

			//If the user did not provide a second arg, exit with status code 1
			if filepath == "" {
				fmt.Println("A filepath as a second arg is required")
				os.Exit(1)
			}

			//Extract the URLs from filepath provided
			URLList := parseUniqueURLsFromFile(filepath)

			//If there are links to ignore, then filter them out from URLs list extracted above, else check regularly
			if len(ignoreList) != 0 {
				URLListWithoutIgnores := []string{}
				for index, pattern := range ignoreList {
					if index != len(ignoreList)-1 {
						combineStrIgnorePatterns += pattern + "|"
					} else {
						combineStrIgnorePatterns += pattern
					}
				}

				//Create regex object to match any ignore links in list
				re := regexp.MustCompile("^(" + combineStrIgnorePatterns + ")")

				//Filter out the ignored links
				for _, link := range URLList {
					if !re.Match([]byte(link)) {
						URLListWithoutIgnores = append(URLListWithoutIgnores, link)
					}
				}

				//Check with URL List that has no ignored links
				checkURL(URLListWithoutIgnores)
			} else {
				//Check regularly if there is nothing to ignore
				checkURL(URLList)
			}
			return
		}

		if *telescopeFlag {
			var urls []string
			urls = parseFromTelescope()
			if *jflag {

				checkURLJson(urls)
			} else {

				fmt.Println()
				fmt.Println(">>  ***** UrlChecker is working now...... *****  <<")
				fmt.Println("--------------------------------------------------------------------------------------------------")
				checkURL(urls)
			}
		}

		//use for loop to deal with multiple file paths
		i := 1
		for i+1 <= len(os.Args) {

			var urls []string

			if os.Args[i][0] != '-' {

				//call functions to check the availability of each url
				urls = parseUniqueURLsFromFile(os.Args[i])

				//check if there are flags for JSON output or not
				if *jflag {

					checkURLJson(urls)
				} else {

					fmt.Println()
					fmt.Println(">>  ***** UrlChecker is working now...... *****  <<")
					fmt.Println("--------------------------------------------------------------------------------------------------")
					checkURL(urls)
				}
			}
			i++

		}

	}
}
