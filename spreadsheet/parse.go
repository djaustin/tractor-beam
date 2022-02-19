package spreadsheet

import (
	"fmt"
	"strings"

	l "github.com/djaustin/tractor-beam/logger"
	"github.com/xuri/excelize/v2"
)

func ExtractPairs(path, sheet, keyHeader, valHeader string) (map[string]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	l.Logger.Infof("opened spreadsheet file '%s'", path)

	cols, err := f.GetCols(sheet)
	if err != nil {
		return nil, err
	}

	keyColIdx, valColIdx, err := findColumnIndices(cols, keyHeader, valHeader)
	if err != nil {
		return nil, err
	}
	rows, err := f.GetRows(sheet)
	if err != nil {
		l.Logger.Fatal(err)
	}

	results := make(map[string]string, len(rows)-1)

	for idx, row := range rows {
		if idx < 1 {
			// Skip column header
			continue
		}
		key, value := row[keyColIdx], row[valColIdx]
		results[key] = value
	}
	l.Logger.Infof("scanned %d rows of data", len(rows)-1)
	return results, nil
}

func findColumnIndices(cols [][]string, keyHeader, valHeader string) (keyColIdx, valColIdx int, err error) {
	var keyColFound, valColFound bool

	for idx, column := range cols {
		if strings.EqualFold(keyHeader, column[0]) {
			keyColIdx = idx
			keyColFound = true
			l.Logger.Infof("found column '%s' at index %d", keyHeader, idx)
		}
		if strings.EqualFold(valHeader, column[0]) {
			valColIdx = idx
			valColFound = true
			l.Logger.Infof("found column '%s' at index %d", valHeader, idx)
		}
		if valColFound && keyColFound {
			break
		}
	}
	if !valColFound {
		return keyColIdx, valColIdx, fmt.Errorf("failed to find value column '%s'", valHeader)
	}
	if !keyColFound {
		return keyColIdx, valColIdx, fmt.Errorf("failed to find key column '%s'", keyHeader)
	}
	return keyColIdx, valColIdx, nil
}
