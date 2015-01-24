package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

//Assumptions in this version:
//Program receives a string of the command that is supposed to be run.
//ie stk go run execTest.go

//If we end up needing something a bit closer to the metal,
//look at os.StartProcess before getting hackish

//To truly get stderr, we would need to intercept any write call to the STDERR
//But that's hard, so we are going to use exec.Cmd on the first go around.

//Any init code that we need will eventually be put in here
func init() {
	fmt.Println("Starting Up.")
}

//the main loop is probably going to look like:
//1.Process provided string into an executable command
//2.Exec them
//3.Have a go routine running to capture any err output then pass them off to
//  the API call,
// 4. Get results, prepend file name to whatever the output was from the api
func main() {
	//This will choke if more than one cmd is passed
	cmd, err := cleanInput(os.Args[1:]...)
	if err != nil {
		log.Fatalf("The provided command is not installed: %T %v",
			err,
			err)
	}
	execCmd(cmd)
}

//CleanInput takes all the relevant arguments from os.Args
//and tries to break it down into an exec.Cmd struct
//This will need a lot of tuning as it will be fragile
func cleanInput(arg ...string) (cmd *exec.Cmd, err error) {
	if len(arg) <= 0 {
		log.Fatalln("Must provide input.")
	}
	fmt.Printf("Args: %v\n", arg)
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
		log.Fatal("Pip conn err: ", e)
	}
	reader := bufio.NewReader(stderr)
	if e := cmd.Start(); e != nil {
		log.Fatal("Process Start Failed", e)
	}
	processErrs(reader)
}

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
