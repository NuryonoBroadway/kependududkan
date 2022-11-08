package main

import (
	"flag"
	"fmt"
	"kependudukan-patalan/initiate"
	"math"
	"os"
	"strings"
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

	edit_cmd := flag.NewFlagSet("edit", flag.ExitOnError)
	edit_source := edit_cmd.String("source", "", "source to define a source file")
	part := edit_cmd.String("part", "", "target to define which column to edit. i.e: primary_key;header_nam:value|etch...")

	if len(os.Args) < 2 {
		log.Fatal("expected 'one' or 'two' subcommands")
	}

	wg := new(sync.WaitGroup)

	switch os.Args[1] {
	case "add":
		add_cmd.Parse(os.Args[2:])

		if !findMe(*dukuh) {
			log.Fatal("dukuh not found")
		}

		in := initiate.NewInit(*dukuh, os.Args[1])
		reader, err := in.OpenFile(*add_source, *target)
		if err != nil {
			log.Fatal(err)
		}
		defer reader.CsvFile.Close()
		defer reader.ExcelFile.Close()
		log.Info(in.Dukuh)

		targetSource := make(chan []string)

		wg.Add(1)
		go in.DispatchWorkers(targetSource, wg, *add_source)
		in.ReadFile(reader, reader.CsvFile, targetSource, wg)

		wg.Wait()

	case "edit":
		edit_cmd.Parse(os.Args[2:])

		in := initiate.NewInit("", os.Args[1])
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

func findMe(location string) bool {
	log.Infof("find location: %v", location)
	dukuh := map[string]bool{
		"bakulan kulon":  true, // 1
		"bakulan wetan":  true, // 2
		"ngaglik":        true, // 3
		"gelangan":       true, // 4
		"tanjung lor":    true, // 5
		"jetis":          true, // 6
		"tanjung karang": true, // 7
		"gaduh":          true, // 8
		"patalan":        true, // 9
		"karang asem":    true, // 10
		"panjang jiwo":   true, // 11
		"gerselo":        true, // 12
		"sulang lor":     true, // 13
		"sulang kidul":   true, // 14
		"dukuh sukun":    true, // 15
		"butuh":          true, // 16
		"boto":           true, // 17
		"kategan":        true, // 18
		"ketandan":       true, // 19
		"bobok":          true, // 20
	}

	return dukuh[strings.ToLower(location)]
}
