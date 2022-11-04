package initiate

import (
	"encoding/csv"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

type M map[string]interface{}

type Init struct {
	TotalWorker int
	Dukuh       string
}

type File struct {
	excel     *excelize.File
	csv       *csv.Reader
	ExcelFile *os.File
	CsvFile   *os.File
}

func NewInit(dukuh string) Init {
	return Init{
		TotalWorker: 100,
		Dukuh:       strings.ToLower(dukuh),
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
