package spreadsheet

import (
	"fmt"
	"strings"

	l "github.com/djaustin/tractor-beam/logger"
	"github.com/xuri/excelize/v2"
)

func ExtractMaps(path, sheet, keyHeader string, valueHeaders ...string) (map[string]map[string]interface{}, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	l.Logger.Debugf("opened spreadsheet file '%s'", path)

	cols, err := f.GetCols(sheet)
	if err != nil {
		return nil, err
	}

	colIndices, err := findColumnIndices(cols, append(valueHeaders, keyHeader)...)
	if err != nil {
		return nil, err
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		l.Logger.Fatal(err)
	}

	results := make(map[string]map[string]interface{}, len(rows)-1)
	for i, row := range rows {
		if i < 1 {
			continue
		}
		valueMap := make(map[string]interface{}, len(valueHeaders))
		for header, index := range colIndices {
			valueMap[header] = row[index]
		}
		key := row[colIndices[strings.ToLower(keyHeader)]]
		results[key] = valueMap
	}
	return results, nil
}

func difference(searchSet map[string]struct{}, foundSet map[string]int) []string {
	var diff []string

	for k := range searchSet {
		if _, found := foundSet[k]; !found {
			diff = append(diff, k)
		}
	}
	return diff
}

func findColumnIndices(cols [][]string, headers ...string) (map[string]int, error) {

	// Create a set of headers to search for O(1) lookup
	searchSet := make(map[string]struct{}, len(headers))
	for _, v := range headers {
		searchSet[strings.ToLower(v)] = struct{}{}
	}

	foundColumns := make(map[string]int, len(headers))

	for i, col := range cols {
		normalisedHeader := strings.ToLower(col[0])
		if _, ok := searchSet[normalisedHeader]; ok {
			foundColumns[normalisedHeader] = i
		}
		if len(foundColumns) == len(searchSet) {
			// All columns found
			return foundColumns, nil
		}
	}

	return nil, fmt.Errorf("failed to find columns %+v in spreadsheet", difference(searchSet, foundColumns))
}
