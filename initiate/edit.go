package initiate

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

func (i Init) SwitchAnything(file *os.File, jobs chan<- []string, wg *sync.WaitGroup, required string, path string) {
	partition := &partition{}

	parts := strings.Split(required, "|")
	for _, part := range parts {
		each := strings.Split(part, ";")
		if len(each) < 2 {
			log.Error("required not meet")
			return
		}

		partition.part = append(partition.part, row{
			id:     each[0],
			change: each[1],
			value:  each[2],
		})

	}

	header := []string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := strings.Split(scanner.Text(), ",")
		if len(header) == 0 {
			header = row
			continue
		}

		for _, part := range partition.part {
			if findHeader(header, part.change) {
				if findHeader(row, part.id) {
					log.Infof("found %v", part.id)
					row[searchIndex(header, part.change)] = part.value
					continue
				} else {
					continue
				}
			}
		}

		wg.Add(1)
		jobs <- row
	}

	if err := scanner.Err(); err != nil {
		log.Error(err)
		return
	}
	close(jobs)

}

func findHeader(header []string, target string) bool {
	for _, header := range header {
		if strings.Contains(header, target) {
			return true
		}
	}

	return false
}

func (i Init) MergeContent(jobs <-chan []string, wg *sync.WaitGroup, path string) {
	var bs []byte
	buf := bytes.NewBuffer(bs)

	// log.Info(job)

	counter := 1

	for job := range jobs {
		i.doTheEditJob(counter, buf, strings.Join(job, ","))
		wg.Done()
		counter++
	}

	log.Info(buf.Len())
	if err := os.WriteFile(path, buf.Bytes(), 0666); err != nil {
		log.Error(err)
		return
	}

	wg.Done()

}

func (i Init) doTheEditJob(counter int, buf *bytes.Buffer, job string) {
	for {
		var outerError error
		func(outerError *error) {
			defer func() {
				if err := recover(); err != nil {
					*outerError = fmt.Errorf("%v", err)
				}
			}()

			log.Info(job)
			_, err := buf.Write([]byte(job))
			if err != nil {
				log.Fatal(err)
			}
			_, err = buf.WriteString("\n")
			if err != nil {
				log.Fatal(err)
			}

		}(&outerError)
		if outerError == nil {
			break
		}
	}

	log.Println("=> inserted", counter, "data")

}
