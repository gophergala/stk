package main

import (
	"fmt"
)

func main() {
	/*result, _ := Search("drush failed")
	for _, item := range result.Items {
		println(item.Title)
	}*/

	//printExplanation("TEst", "https://script.fail")

	printError("Error occured", "Something went wrong", "http://script.fail")

}

func printError(errstr string, maybeReason string, detailUrl string) {
	fmt.Println(errstr)
	fmt.Println()
	fmt.Println(bold("Possible reason:"))
	fmt.Println(maybeReason)
	fmt.Println()
	fmt.Println(bold("Details: "))
	fmt.Println(underline(detailUrl))
	fmt.Println()
}

func bold(text string) string {
	return "\033[1m" + text + "\033[0m"
}

func underline(text string) string {
	return "\033[4m" + text + "\033[0m"
}
