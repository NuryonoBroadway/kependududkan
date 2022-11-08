package initiate

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

func (i Init) ReadFile(reader *File, file *os.File, jobs chan<- []string, wg *sync.WaitGroup) {
	row := make([][]string, 0)
	header := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		each := strings.Split(scanner.Text(), ",")
		if len(header) == 0 {
			header = each
			continue
		}

		row = append(row, each)
	}

	if err := scanner.Err(); err != nil {
		log.Error(err)
		return
	}

	func() {
		wg.Add(1)
		jobs <- header
	}()

	log.Info(header)
	for _, sheet := range reader.excel.GetSheetMap() {
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
					rows[i] = append(rows[i], "-")
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
			tolerance := 0
			schema := make([]string, 0)
			for where, lobby := range header {
				if where == 0 {
					schema = append(schema, time.Now().Local().Format("01/02/2006 15:04:05"))
					continue
				} else if where == 1 {
					schema = append(schema, strings.ToUpper(i.Dukuh))
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
						if lobby == head {
							if rows[l][index] == "-" {
								tolerance += 1
							}
							schema = append(schema, rows[l][index])
						} else if i.sameHeader(lobby, head) {
							schema = append(schema, strings.ReplaceAll(rows[l][index], " ", ""))
						}
					}
				}
			}

			if tolerance > 3 {
				tolerance = 0
				continue
			}

			if slices.Contains(rows[0], "NIK") {
				for i, each := range row {
					if i == 0 {
						continue
					}
					if slices.Contains(each, rows[l][searchIndex(rows[0], "NIK")]) {
						row[i] = schema
					}
				}
			}

			tolerance = 0
		}
	}

	for _, row := range row {
		if len(row) < 1 {
			continue
		} else {
			wg.Add(1)
			jobs <- row
		}
	}

	close(jobs)
}

func searchIndex(row []string, target string) int {
	for i, value := range row {
		if strings.Contains(value, target) {
			return i
		}
	}

	return -1
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

func (i Init) DispatchWorkers(jobs <-chan []string, wg *sync.WaitGroup, path string) {
	var bs []byte
	buf := bytes.NewBuffer(bs)

	counter := 1

	for job := range jobs {
		i.doTheExcelJob(counter, buf, strings.Join(job, ","))
		wg.Done()
		counter++
	}

	log.Info(buf.Len())
	if err := os.WriteFile(path, buf.Bytes(), 0777); err != nil {
		log.Error(err)
		return
	}

	wg.Done()
}

func (i Init) doTheExcelJob(counter int, buf *bytes.Buffer, job string) {
	for {
		var outerError error
		func(outerError *error) {
			defer func() {
				if err := recover(); err != nil {
					*outerError = fmt.Errorf("%v", err)
				}
			}()

			// log.Info(job)
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

	// log.Println("=> inserted", counter, "data")
}
