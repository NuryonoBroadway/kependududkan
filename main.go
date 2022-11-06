package main

import (
	"flag"
	"fmt"
	"kependudukan-patalan/initiate"
	"math"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	start := time.Now()

	add_cmd := flag.NewFlagSet("add", flag.ExitOnError)
	add_source := add_cmd.String("source", "", "source to define a source file")
	target := add_cmd.String("target", "", "target to define a target file")
	dukuh := add_cmd.String("dukuh", "", "target to define a dukuh")
	output := add_cmd.String("output", ".", "to define a output")

	edit_cmd := flag.NewFlagSet("edit", flag.ExitOnError)
	edit_source := edit_cmd.String("source", "", "source to define a source file")
	part := edit_cmd.String("part", "", "target to define which column to edit. i.e: primary_key;header_nam:value|etch...")

	if len(os.Args) < 2 {
		fmt.Println("expected 'one' or 'two' subcommands")
		os.Exit(1)
	}

	in := initiate.NewInit(*dukuh, os.Args[1])
	wg := new(sync.WaitGroup)

	switch os.Args[1] {
	case "add":
		add_cmd.Parse(os.Args[2:])
		reader, err := in.OpenFile(*add_source, *target)
		if err != nil {
			log.Fatal(err)
		}
		defer reader.CsvFile.Close()
		defer reader.ExcelFile.Close()

		info, err := reader.CsvFile.Stat()
		if err != nil {
			log.Info("cant find info")
			log.Fatal(err)
		}

		exit, err := os.Create(fmt.Sprintf("%v/%v", *output, info.Name()))
		if err != nil {
			log.Info("create failed")
			log.Fatal(err)
		}
		defer exit.Close()

		targetSource := make(chan []interface{})

		go in.DispatchWorkers(targetSource, exit, wg)
		in.ReadFile(reader, targetSource, wg)

		wg.Wait()

	case "edit":
		edit_cmd.Parse(os.Args[2:])
		reader, err := in.OpenFile(*edit_source, *target)
		if err != nil {
			log.Fatal(err)
		}
		defer reader.CsvFile.Close()

		targetSource := make(chan []string)

		wg.Add(1)
		go in.MergeContent(targetSource, wg, *edit_source)

		in.SwitchAnything(reader.CsvFile, targetSource, wg, *part, *edit_source)

		wg.Wait()
	}

	duration := time.Since(start)
	fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")
}
