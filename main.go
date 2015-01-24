package main

import (
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
	//	cmd := cleanInput(os.Args[1:])
	//This will choke if to cmd are passed at once,ie
	cmd := cleanInput(os.Args[1:]...)
	execCmd(&cmd)
}

//CleanInput takes all the relevant arguments from os.Args
//and tries to break it down into an exec.Cmd struct
//This will need a lot of tuning as it will be fragile
func cleanInput(arg ...string) exec.Cmd {
	fmt.Printf("Args: %v\n", arg)
	return exec.Cmd{}
}

//This is going to be the main event loop in all actuality
//It will launch the provided task and attaches itself to stdErr,
//blocking on the exit of the cmd
//Redirects the stderr(which expects an io.Writer) into a channel,
//which the API is blocking on in order to launch a request.
func execCmd(cmd *exec.Cmd) {
	//Make sure the actual command exists, ie go
	path, err := exec.LookPath(cmd.Path)
	if err != nil {
		log.Fatalf("Bloody hell, at least install the command: %s\n",
			cmd.Path)
	}
	fmt.Printf("Yay?%s\n", path)
}
