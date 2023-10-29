package csv

import (
	"anicomend/model"
	"encoding/csv"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
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

		aired := strings.Split(record[11], " to ")
		if len(aired) == 1 {
			aired = append(aired, "Jan 1, 2024")
		} else if aired[1] == "?" {
			aired[1] = "Jan 1, 2024"
		}

		animes[i] = model.Anime{
			Id:        id,
			Title:     record[1],
			ImageURL:  record[5],
			Type:      record[6],
			Source:    record[7],
			Studio:    record[27],
			Genres:    strings.Split(record[28], ", "),
			Duration:  float32(duration),
			Year:      uint32(year),
			AiredFrom: parseDateTime(aired[0]),
			AiredTo:   parseDateTime(aired[1]),
		}

		// filter empty genres
		animes[i].Genres = slices.DeleteFunc(animes[i].Genres, func(g string) bool {
			return len(strings.TrimSpace(g)) == 0
		})
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

func parseDateTime(st string) time.Time {
	if st != "?" {
		t, _ := time.Parse("Jan 2, 2006", st)
		return t
	}
	return time.Now()
}
