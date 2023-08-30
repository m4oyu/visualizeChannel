package main

import (
	"regexp"
	"strings"
)

func main() {
	// if len(os.Args) < 2 {
	// 	fmt.Println("Usage: mytool <go-file>")
	// 	return
	// }

	// goFilePath := os.Args[1]
	// goFileContent, err := os.ReadFile(goFilePath)
	// if err != nil {
	// 	fmt.Println("Error reading file:", err)
	// 	return
	// }

	// fmt.Println(string(goFileContent))

	// // Inject channel operation into the file content
	// modifiedContent := injectChannelOperation(string(goFileContent))

	// modifiedFilePath := "modified_" + goFilePath
	// err = os.WriteFile(modifiedFilePath, []byte(modifiedContent), 0644)
	// if err != nil {
	// 	fmt.Println("Error writing modified file:", err)
	// 	return
	// }

	// fmt.Println("Modified file written to:", modifiedFilePath)

	regexp.MustCompilePOSIX()

}

func injectChannelOperation(input string) string {
	// Inject channel operation here
	// For example, replace "ch <- value" with "ch <- value + 1"

	// Inject make
	// ch := make(chan int, int)
	1 := make(chan int, 3)

	// Inject send

	// Inject recv

	return strings.Replace(input, "ch <- i", "ch <- i + 1", -1)
}
