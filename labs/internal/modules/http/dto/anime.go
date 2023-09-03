package dto

import (
	"anicomend/model"
	"fmt"
	"strings"
)

type AnimeDTO struct {
	Id       uint64
	Title    string
	ImageURL string
	Type     string
	Source   string
	Studio   string
	Genres   string
	Duration string
	Year     string
}

func NewAnimeDTO(anime model.Anime) AnimeDTO {
	return AnimeDTO{
		Id:       anime.Id,
		Title:    anime.Title,
		ImageURL: anime.ImageURL,
		Type:     anime.Type,
		Source:   anime.Source,
		Studio:   anime.Studio,
		Genres:   strings.Join(anime.Genres, ", "),
		Duration: fmt.Sprintf("%.0f min", anime.Duration),
		Year:     fmt.Sprintf("%d", anime.Year),
	}
}
