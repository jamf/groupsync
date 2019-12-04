package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/google/logger"
	"github.com/jamf/groupsync/cmd"
	"github.com/jamf/groupsync/services"
	"github.com/logrusorgru/aurora"
)

func main() {
	defer logger.Init("groupsyncLogger", debug, true, ioutil.Discard).Close()
	logger.SetFlags(log.Lshortfile)

	if debug {
		fmt.Println(
			aurora.Red("You're using a dev build of groupsync. "),
		)
	}

	err := services.Init()
	if err != nil {
		logger.Fatal(err)
	}

	cmd.Execute()
}
