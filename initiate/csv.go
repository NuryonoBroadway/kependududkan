package initiate

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

// const (
// 	totalWroker = 100
// )

// var dataHeaders = make([]string, 0)

func (i Init) OpenCsvFile(source string) (*File, *os.File, error) {
	log.Info("open csv file")
	drop := &File{}

	ext := filepath.Ext(source)

	var file *os.File
	var err error

	switch ext {
	case ".csv":
		file, err = os.Open(source)
		if err != nil {
			log.Info(err)
			return nil, nil, err
		}
		drop.csv = csv.NewReader(file)
	case ".xlsx":
		file, err = os.Open(source)
		if err != nil {
			log.Info(err)
			return nil, nil, err
		}
		drop.excel, err = excelize.OpenReader(file)
		if err != nil {
			return nil, nil, nil
		}
	}

	return drop, file, nil
}

func (i Init) DispatchCsvWorkers(jobs <-chan []interface{}, file *os.File, wg *sync.WaitGroup) {
	index := 0

	// log.Info(job)
	for index <= i.TotalWorker {
		go func(worker int, jobs <-chan []interface{}, file *os.File, wg *sync.WaitGroup) {
			counter := 0

			for job := range jobs {
				i.doTheCsvJob(worker, counter, file, job)
				wg.Done()
				counter++
			}
		}(index, jobs, file, wg)
		index++
	}
}

func (i Init) ReadCsvFile(reader *File, jobs chan<- []interface{}, wg *sync.WaitGroup) []string {
	header := make([]string, 0)
	for {
		row, err := reader.csv.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}

		if len(header) == 0 {
			header = row
		}

		schema := make([]interface{}, 0)
		for _, each := range row {
			schema = append(schema, each)
		}

		wg.Add(1)
		jobs <- schema
	}

	return header
}

func (i Init) doTheCsvJob(worker, counter int, file *os.File, job []interface{}) {
	for {
		var outerError error
		func(outerError *error) {
			defer func() {
				if err := recover(); err != nil {
					*outerError = fmt.Errorf("%v", err)
				}
			}()

			w := csv.NewWriter(file)
			defer w.Flush()

			builder := []string{}
			for _, j := range job {
				builder = append(builder, fmt.Sprint(j))
			}

			log.Info(builder)
			if err := w.Write(builder); err != nil {
				log.Error(err)
				return
			}

		}(&outerError)
		if outerError == nil {
			break
		}
	}

	log.Println("=> worker", worker, "inserted", counter, "data")
}
