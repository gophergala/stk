package main

import (
	"fmt"
	"log"
)

func main() {
	// LOGIC FOR CAPTURING STDERR

	reason, url := findReason("drush failed", "", "")
	printError("Error occured", reason, url)
}

func findReason(strerr, command, parameters string) (reason string, url string) {
	res, err := Search(strerr)

	if err != nil {
		log.Fatal(err)
	}

	if len(res.Items) == 0 {
		return "", ""
	}

	answerId := res.Items[0].AcceptedAnswerId
	answer, err := GetAnswers(answerId)

	if err != nil {
		log.Fatal(err)
	}

	if len(answer.Items) == 0 {
		return "", ""
	}

	reason = answer.Items[0].Body
	url = res.Items[0].Link
	return
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
