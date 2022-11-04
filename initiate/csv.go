package initiate

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

// const (
// 	totalWroker = 100
// )

// var dataHeaders = make([]string, 0)

func (i Init) OpenFile(source string, target string) (*File, error) {
	log.Info("open csv file")
	drop := &File{}

	var err error

	drop.CsvFile, err = os.Open(source)
	if err != nil {
		log.Info(err)
		return nil, err
	}
	drop.csv = csv.NewReader(drop.CsvFile)

	drop.ExcelFile, err = os.Open(target)
	if err != nil {
		log.Info(err)
		return nil, err
	}
	drop.excel, err = excelize.OpenReader(drop.ExcelFile)
	if err != nil {
		return nil, nil
	}

	return drop, nil
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
			continue
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
