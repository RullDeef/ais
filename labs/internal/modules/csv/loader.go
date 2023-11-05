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

var filteredGenres []string = []string{
	"Dementia",
	"Harem",
	"Kids",
	"Ecchi",
	"Shounen Ai",
	"Yuri",
	"Hentai",
	"Yaoi",
	"Shoujo Ai",
}

type CSVAnimeLoader struct {
	filename       string
	filterAdult    bool
	updateImageCDN bool
}

func NewAnimeLoader(filename string) *CSVAnimeLoader {
	return &CSVAnimeLoader{
		filename:       filename,
		filterAdult:    true,
		updateImageCDN: true,
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

		episodes, err := strconv.ParseUint(record[8], 10, 32)
		if err != nil {
			return nil, err
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
			Episodes:  uint32(episodes),
			Year:      uint32(year),
			AiredFrom: parseDateTime(aired[0]),
			AiredTo:   parseDateTime(aired[1]),
		}

		// filter empty genres
		animes[i].Genres = slices.DeleteFunc(animes[i].Genres, func(g string) bool {
			return len(strings.TrimSpace(g)) == 0
		})
	}

	if c.filterAdult {
		animes = filterGenres(animes)
	}

	// minAired := time.Now().UTC()
	// maxAired := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)

	if c.updateImageCDN {
		for i, anime := range animes {
			animes[i].ImageURL = strings.ReplaceAll(anime.ImageURL, "myanimelist.cdn-dena.com", "cdn.myanimelist.net")
			// if anime.AiredFrom.Year() != 1 && anime.AiredFrom.Before(minAired) {
			// 	minAired = anime.AiredFrom
			// }
			// if anime.AiredFrom.After(maxAired) {
			// 	maxAired = anime.AiredFrom
			// }
		}
	}

	// fmt.Printf("min aired: %+v\n", minAired)	// min aired: 1942-03-15 00:00:00 +0000 UTC
	// fmt.Printf("max aired: %+v\n", maxAired)	// max aired: 2018-05-25 00:00:00 +0000 UTC

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

func filterGenres(animes []model.Anime) []model.Anime {
	filtered := make([]model.Anime, 0)
	for _, anime := range animes {
		hasFilteredGenre := false
		for _, genre := range anime.Genres {
			if slices.Contains(filteredGenres, genre) {
				hasFilteredGenre = true
				break
			}
		}
		if !hasFilteredGenre && len(anime.Genres) > 0 {
			filtered = append(filtered, anime)
		}
	}
	return filtered
}
