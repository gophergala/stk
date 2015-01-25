package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	sto "github.com/gophergala/stk/stackoverflow"
	"gopkg.in/alecthomas/kingpin.v1"
)

//could use the command as a possible tag
//Assumptions in this version:
//Program receives a string of the command that is supposed to be run.
//ie stk go run execTest.go

//If we end up needing something a bit closer to the metal,
//look at os.StartProcess before getting hackish

//To truly get stderr, we would need to intercept any write call to the STDERR
//But that's hard, so we are going to use exec.Cmd on the first go around.

var (
	errFileFlag = kingpin.Flag("errFile",
		"Output errors to a file in the pwd with the timestamp for a name.").Default("false").Short('e').Bool()
	commandArgs = kingpin.Arg("command", "Command being run").
			Required().
			Strings()
	cmd     *exec.Cmd
	err     error
	errFile *os.File
	nohtml  *regexp.Regexp
)

//Any init code that we need will eventually be put in here
func init() {
	kingpin.Parse()
	cleanInput()
	if cmd.Path == "" {
		log.Fatalln("The provided command is not installed")
	}
	if *errFileFlag {
		errFile, err = os.Create(time.Now().UTC().
			Format("Jan 2, 2006 at 3:04pm (MST)"))
		if err != nil {
			log.Fatalln("File Creation err: ", err)
		}
	}
	log.Printf("Starting Up. %#v", commandArgs)
	nohtml, _ = regexp.Compile("<[^>]*>")

}

func stripHtml(htmlContent string) string {
	return nohtml.ReplaceAllString(htmlContent, "")
}

//the main loop is probably going to look like:
//1.Process provided string into an executable command
//2.Exec them
//3.Have a go routine running to capture any err output then pass them off to
//  the API call,
// 4. Get results, prepend file name to whatever the output was from the api
func main() {
	if *errFileFlag {
		defer errFile.Close()
	}
	//This will choke if more than one cmd is passed
	execCmd()
	//	stderr := "The drush command could not be found"

	//	reason, url := findReason(stderr, "", "")
	//	sanitized := stripHtml(reason)

	//	printError(stderr, sanitized, url)
}

//CleanInput takes all the relevant arguments from os.Args
//and tries to break it down into an exec.Cmd struct
//This will need a lot of tuning as it will be fragile
func cleanInput() {
	if len(*commandArgs) <= 0 {
		log.Fatalln("Must provide input.")
	}
	if len(os.Args) > 2 {
		cmd = exec.Command((*commandArgs)[0], (*commandArgs)[1:]...)
	} else {
		cmd = exec.Command((*commandArgs)[0])
	}
	log.Printf("cmd.Args: %#v", cmd.Args)
	return
}

//This is going to be the main event loop in all actuality
//It will launch the provided task and attaches itself to stdErr,
//blocking on the exit of the cmd
//Redirects the stderr(which expects an io.Writer) into a channel,
//which the API is blocking on in order to launch a request.
func execCmd() {
	stderr, e := cmd.StderrPipe()
	if e != nil {
		log.Fatal("Pipe conn err: ", e)
	}
	reader := bufio.NewReader(stderr)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("Pipe conn err: ", err)
	}
	r := bufio.NewReader(stdout)
	if e := cmd.Start(); e != nil {
		log.Fatal("Process Start Failed", e)
	}
	errChan := make(chan string)
	go passStdOut(r)
	go processErrs(reader, errChan)

	//Problem? If the command exits it passes back a Proc state to err which will prompt an exit before the go routine can even process.
	//Solution is channels.
	if err := cmd.Wait(); err != nil {
		//Type is exit error
		//log.Fatal("Problem?", err)
		select {
		case <-errChan:
			log.Fatal(err)
		}
	}

}

//processErrs is the function that launches the requests to the API
func processErrs(reader *bufio.Reader, errChan chan<- string) {
	var writer *bufio.Writer
	if *errFileFlag {
		writer = bufio.NewWriter(errFile)
	}
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println("Read err", err)
				errChan <- err.Error()
			}
			continue
		} else {
			log.Println("Captured: ", s)
			reason, url := findReason(s, (*commandArgs)[0], "")
			printError("Error Captured:", stripHtml(reason), url)
			if *errFileFlag {
				n, e := writer.WriteString(s + "\n")
				if e != nil {
					log.Printf("Bytes written: %d.Err:%v",
						n, err)

				}
				//I want to defer this flush till exit but
				//that would mean adding it to the main func
				//which would require a new global var
				writer.Flush()
			}
			errChan <- s
		}
	}
}

func passStdOut(r *bufio.Reader) {
	_, err := r.WriteTo(os.Stdout)
	if err != nil {
		log.Println("Write err", err)
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

func xterm(code string) func(s string) string {
	env := os.Getenv("TERM")
	isXterm := strings.Contains(env, "xterm")

	return func(text string) (output string) {
		if isXterm {
			output = code + text + "\033[0m"
		} else {
			output = text
		}
		return
	}
}

func bold(text string) string {
	return xterm("\033[1m")(text)
}

func underline(text string) string {
	return xterm("\033[4m")(text)
}
