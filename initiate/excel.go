package initiate

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (i Init) ReadExcelFile(reader *File, jobs chan<- []interface{}, wg *sync.WaitGroup) {
	// family := make(map[string][]interface{})

	log.Info(i.SourceHeaders)
	for _, sheet := range reader.excel.GetSheetMap() {
		// sheet := "RT 65"
		log.Infof("sheets %v", sheet)
		rt := strings.Split(sheet, " ")[1]

		rows, err := reader.excel.GetRows(sheet)
		if err != nil {
			fmt.Println(err)
			return
		}

		kk := ""
		for i, d := range rows {
			if len(d) != len(rows[0]) {
				for l := len(d) + 1; l < len(rows[0]); l++ {
					rows[i] = append(rows[i], "#")
				}
			}

			if d[0] != "" {
				kk = d[0]
			} else {
				rows[i][0] = kk
			}

			for j, k := range d {
				if j == 0 {
					continue
				}

				if k == "" {
					rows[i][j] = "-"
				}
			}
		}

		for l := 2; l < len(rows); l++ {
			schema := make([]interface{}, 0)
			for where, lobby := range i.SourceHeaders {
				if where == 0 {
					loc, _ := time.LoadLocation("Asia/Jakarta")
					schema = append(schema, time.Now().In(loc).Format("01/02/2006 15:04:05"))
					continue
				} else if where == 1 {
					schema = append(schema, cases.Title(language.Indonesian, cases.Compact).String(i.Dukuh))
					continue
				} else if strings.Contains(lobby, "RT") {
					if where == whichRow(i.Dukuh) {
						schema = append(schema, rt)
						continue
					} else {
						schema = append(schema, "")
						continue
					}
				} else {
					for index, head := range rows[0] {
						if lobby == head || i.sameHeader(lobby, head) {
							schema = append(schema, rows[l][index])
							break
						}
					}
				}
			}

			wg.Add(1)
			jobs <- schema
		}

	}
	close(jobs)
}

func (i Init) DispatchExelWorkers(jobs <-chan []interface{}, file *os.File, wg *sync.WaitGroup) {
	index := 0

	// log.Info(job)
	for index <= i.TotalWorker {
		go func(worker int, jobs <-chan []interface{}, file *os.File, wg *sync.WaitGroup) {
			counter := 0

			for job := range jobs {
				i.doTheExcelJob(worker, counter, file, job)
				wg.Done()
				counter++
			}
		}(index, jobs, file, wg)
		index++
	}
}

func (i Init) sameHeader(lobby, header string) bool {
	if strings.Contains(strings.ToLower(header), "kartu keluarga") && strings.Contains(strings.ToLower(lobby), "kk") {
		return true
	}

	if header == "No." {
		return false
	}

	if strings.Contains(lobby, header) {
		return true
	}

	return false
}

func (i Init) doTheExcelJob(worker, counter int, file *os.File, job []interface{}) {
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
