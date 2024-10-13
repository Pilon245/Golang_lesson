package main

import (
	"bufio"
	"fmt"
	"github.com/mailru/easyjson"
	"hw3/models"
	"io"
	"os"
	"strings"
)

// вам надо написать более быструю оптимальную этой функции

//type User struct {
//	Browsers []string `json:"browsers"`
//	Email    string   `json:"email"`
//	Name     string   `json:"name"`
//}

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Fprintln(out, "found users:")

	seenBrowsers := make(map[string]interface{})
	index := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var user models.User
		err := easyjson.Unmarshal([]byte(line), &user)
		if err != nil {
			panic(err)
		}

		// Обрабатывай каждую строку

		isAndroid := false
		isMSIE := false

		for _, browser := range user.Browsers {
			if strings.Contains(browser, "Android") {
				isAndroid = true
				if _, exists := seenBrowsers[browser]; !exists {
					seenBrowsers[browser] = struct{}{}
				}
			}
			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				if _, exists := seenBrowsers[browser]; !exists {
					seenBrowsers[browser] = struct{}{}
				}
			}

		}

		// Если пользователь использует и Android, и MSIE, добавляем его в найденные
		if isAndroid && isMSIE {
			email := strings.Replace(user.Email, "@", " [at] ", -1)
			fmt.Fprintf(out, "[%d] %s <%s>\n", index, user.Name, email)
		}
		index++
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}
