package initiate

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

type M map[string]interface{}

type Init struct {
	TotalWorker int
	Dukuh       string
	Mode        string
}

type File struct {
	excel     *excelize.File
	csv       *csv.Reader
	ExcelFile *os.File
	CsvFile   *os.File
}

type partition struct {
	part []row
}

type row struct {
	id     string
	change string
	value  string
}

func NewInit(dukuh string, mode string) Init {
	return Init{
		TotalWorker: 100,
		Dukuh:       strings.ToLower(dukuh),
		Mode:        mode,
	}
}

func (i Init) OpenFile(source string, target string) (*File, error) {
	drop := &File{}
	var err error

	switch i.Mode {
	case "add":
		log.Info("open csv file")
		drop.CsvFile, err = os.Open(source)
		if err != nil {
			log.Info(err)
			return nil, err
		}
		drop.csv = csv.NewReader(drop.CsvFile)

		log.Info("open excel file")
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

	case "edit":
		log.Info("open csv file")
		drop := &File{}

		var err error

		drop.CsvFile, err = os.Open(source)
		if err != nil {
			log.Info(err)
			return nil, err
		}
		drop.csv = csv.NewReader(drop.CsvFile)

		return drop, nil

	default:
		return nil, fmt.Errorf("mode not found")
	}
}

func whichRow(dukuh string) int {
	switch dukuh {
	case "bakulan kulon":
		return 2
	case "bakulan wetan":
		return 3
	case "ngaglik":
		return 5
	case "gelangan":
		return 4
	case "tanjung lor":
		return 6
	case "jetis":
		return 9
	case "tanjung karang":
		return 10
	case "gaduh":
		return 11
	case "patalan":
		return 7
	case "karang asem":
		return 8
	case "panjang jiwo":
		return 12
	case "gerselo":
		return 13
	case "sulang lor":
		return 14
	case "sulang kidul":
		return 16
	case "dukuh sukun":
		return 17
	case "butuh":
		return 15
	case "boto":
		return 19
	case "kategan":
		return 18
	case "ketandan":
		return 20
	case "bobok":
		return 21
	default:
		return 0
	}
}
