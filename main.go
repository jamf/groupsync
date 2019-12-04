package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/google/logger"
	"github.com/jamf/groupsync/cmd"
	"github.com/logrusorgru/aurora"
)

func main() {
	defer logger.Init("groupsyncLogger", debug, true, ioutil.Discard).Close()
	logger.SetFlags(log.Lshortfile)

	if debug {
		fmt.Print(
			aurora.Red("You're using a dev build of groupsync. "),
		)
		fmt.Println("Verbose logging enabled.")
	}

	cmd.Execute()
}
