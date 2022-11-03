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

	source := flag.String("source", "", "source to define a source file")
	target := flag.String("target", "", "target to define a target file")
	dukuh := flag.String("desa", "", "target to define a dukuh")

	flag.Parse()

	in := initiate.NewInit(*dukuh)
	readerCsv, csvFile, err := in.OpenCsvFile(*source)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	readerExel, excelFile, err := in.OpenCsvFile(*target)
	if err != nil {
		log.Fatal(err)
	}
	defer excelFile.Close()

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(path)

	info, err := csvFile.Stat()
	if err != nil {
		log.Fatal(err)
	}

	exit, err := os.Create(fmt.Sprintf("%v/output/%v", path, info.Name()))
	if err != nil {
		log.Info("create failed")
		log.Fatal(err)
	}
	defer exit.Close()

	originSource := make(chan []interface{})
	targetSource := make(chan []interface{})

	wg := new(sync.WaitGroup)

	go in.DispatchCsvWorkers(originSource, exit, wg)
	in.SourceHeaders = in.ReadCsvFile(readerCsv, originSource, wg)

	go in.DispatchExelWorkers(targetSource, exit, wg)
	in.ReadExcelFile(readerExel, targetSource, wg)

	wg.Wait()

	duration := time.Since(start)
	fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")
}
