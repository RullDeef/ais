package csv

import (
	"anicomend/model"
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

type CSVAnimeLoader struct {
	filename string
}

func NewAnimeLoader(filename string) *CSVAnimeLoader {
	return &CSVAnimeLoader{
		filename: filename,
	}
}

func (c *CSVAnimeLoader) Load() ([]model.Anime, error) {
	records, err := readCsvFile(c.filename)
	if err != nil {
		return nil, err
	}

	records = records[1:]

	animes := make([]model.Anime, len(records))
	for i, record := range records {
		id, err := strconv.ParseUint(record[0], 10, 64)
		if err != nil {
			return nil, err
		}
		duration, err := strconv.ParseFloat(record[31], 64)
		if err != nil {
			return nil, err
		}
		year, err := strconv.ParseFloat(record[32], 64)
		if err != nil {
			return nil, err
		}
		animes[i] = model.Anime{
			Id:       id,
			Title:    record[1],
			ImageURL: record[5],
			Type:     record[6],
			Source:   record[7],
			Studio:   record[27],
			Genres:   strings.Split(record[28], ", "),
			Duration: float32(duration),
			Year:     uint32(year),
		}
	}

	return animes, nil
}

func readCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}
