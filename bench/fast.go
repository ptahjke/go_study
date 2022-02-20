package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type User struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Browsers []string `json:"browsers"`
}

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	seenBrowsers := make([]string, 0, 100)
	// Approach with map takes a little bit more memory but mor easy in use
	// seenBrowsers := make(map[string]bool, 100)

	scanner := bufio.NewScanner(file)

	fmt.Fprintln(out, "found users:")

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Bytes()
		if !bytes.Contains(line, []byte("Android")) && !bytes.Contains(line, []byte("MSIE")) {
			continue
		}

		var user User
		err := user.UnmarshalJSON(line)
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		notSeenBefore := false
		for _, browser := range user.Browsers {
			switch {
			case strings.Contains(browser, "Android"):
				isAndroid = true
				notSeenBefore = true
			case strings.Contains(browser, "MSIE"):
				isMSIE = true
				notSeenBefore = true
			default:
				continue
			}

			// map approach
			// seenBrowsers[browser] = true

			if notSeenBefore {
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}

				if notSeenBefore {
					seenBrowsers = append(seenBrowsers, browser)
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		email := strings.Replace(user.Email, "@", " [at] ", 1)
		fmt.Fprintf(out, "[%d] %s <%s>\n", i, user.Name, email)
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}
