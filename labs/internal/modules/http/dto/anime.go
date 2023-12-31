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
	Liked    bool
	Disliked bool
}

func NewAnimeDTO(anime model.Anime, mark model.PreferenceMark) AnimeDTO {
	dto := AnimeDTO{
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
	switch mark.MarkWeight {
	case model.PreferenceMarkFavourite:
		dto.Liked = true
	case model.PreferenceMarkUnfavourite:
		dto.Disliked = true
	}
	return dto
}
