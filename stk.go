package main

import (
	"fmt"
	"log"
	"os"

	sto "github.com/gophergala/stk/stackoverflow"
)

//could use the command as a possible tag

func main() {
	// LOGIC FOR CAPTURING STDERR

	reason, url := findReason("drush failed", "", "")
	printError("Error occured", reason, url)
}

func findReason(strerr, command, parameters string) (reason string, url string) {
	res, err := sto.Search(strerr)

	if err != nil {
		log.Fatal(err)
	}

	if len(res.Items) == 0 {
		return "", ""
	}

	answerID := res.Items[0].AcceptedAnswerId
	answer, err := sto.GetAnswers(answerID)

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

func printError(errstr string, maybeReason string, detailURL string) {
	fmt.Println(errstr)
	fmt.Println()
	fmt.Println(bold("Possible reason:"))
	fmt.Println(maybeReason)
	fmt.Println()
	fmt.Println(bold("Details: "))
	fmt.Println(underline(detailURL))
	fmt.Println()
}

func bold(text string) string {
	if os.Getenv("TERM") == "xterm" {
		return "\033[1m" + text + "\033[0m"
	}
	return text
}

func underline(text string) string {
	if os.Getenv("TERM") == "xterm" {
		return "\033[4m" + text + "\033[0m"
	}
	return text
}
