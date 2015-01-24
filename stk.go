package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"

	sto "github.com/gophergala/stk/stackoverflow"
)

//could use the command as a possible tag
//Assumptions in this version:
//Program receives a string of the command that is supposed to be run.
//ie stk go run execTest.go

//If we end up needing something a bit closer to the metal,
//look at os.StartProcess before getting hackish

//To truly get stderr, we would need to intercept any write call to the STDERR
//But that's hard, so we are going to use exec.Cmd on the first go around.

//Any init code that we need will eventually be put in here
func init() {
	log.Println("Starting Up.")
}

//the main loop is probably going to look like:
//1.Process provided string into an executable command
//2.Exec them
//3.Have a go routine running to capture any err output then pass them off to
//  the API call,
// 4. Get results, prepend file name to whatever the output was from the api
func main() {
	//This will choke if more than one cmd is passed
	/*	cmd, err := cleanInput(os.Args[1:]...)
		if err != nil {
			log.Fatalf("The provided command is not installed: %T %v",
				err,
				err)
		}
		execCmd(cmd)
	*/
	// LOGIC FOR CAPTURING STDERR

	reason, url := findReason("drush failed", "", "")
	printError("Error occured", reason, url)
}

//CleanInput takes all the relevant arguments from os.Args
//and tries to break it down into an exec.Cmd struct
//This will need a lot of tuning as it will be fragile
func cleanInput(arg ...string) (cmd *exec.Cmd, err error) {
	if len(arg) <= 0 {
		log.Fatalln("Must provide input.")
	}
	log.Printf("Args: %v\n", arg)
	if len(os.Args) > 2 {
		cmd = exec.Command(os.Args[1], os.Args[2:]...)
	} else {
		cmd = exec.Command(os.Args[1], "")
	}
	log.Printf("cmd.Args: %#v", cmd.Args)
	return
}

//This is going to be the main event loop in all actuality
//It will launch the provided task and attaches itself to stdErr,
//blocking on the exit of the cmd
//Redirects the stderr(which expects an io.Writer) into a channel,
//which the API is blocking on in order to launch a request.
func execCmd(cmd *exec.Cmd) {
	stderr, e := cmd.StderrPipe()
	if e != nil {
		log.Fatal("Pipe conn err: ", e)
	}
	reader := bufio.NewReader(stderr)
	if e := cmd.Start(); e != nil {
		log.Fatal("Process Start Failed", e)
	}
	go processErrs(reader)

	if err := cmd.Wait(); err != nil {
		//Type is exit error
		log.Fatal(err)
	}
}

//processErrs is the function that launches the requests to the API
func processErrs(reader *bufio.Reader) {
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Read err", err)
			return
		}
		log.Println("Captured: ", s)
	}
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
