package main

import (
	"bufio"
	"computerClub/app"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no input file provided")
		return
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	c, es, err := app.ParseInput(bufio.NewScanner(f))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	logger := log.Default()
	logger.SetFlags(0)
	a := app.App{
		Club: c,
		Log:  logger,
	}
	a.Run(es)
}
