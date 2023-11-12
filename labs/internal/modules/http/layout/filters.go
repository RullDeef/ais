package layout

import (
	"anicomend/service/filter"
	"strings"
)

type GenreOption struct {
	ID    string
	Label string
}

type DurationCategory struct {
	ID    string
	Label string
}

type TypeOption struct {
	ID    string
	Label string
}

type FilterParams struct {
	GenreOptions       []GenreOption
	DurationCategories []DurationCategory
	TypeOptions        []TypeOption
}

func NewFilterParams() FilterParams {
	genres := []string{
		"Psychological",
		"Action",
		"Shounen",
		"Supernatural",
		"Sports",
		"Martial Arts",
		"Historical",
		"Demons",
		"Josei",
		"Space",
		"Mystery",
		"Vampire",
		"Cars",
		"Super Power",
		"Seinen",
		"Sci-Fi",
		"Magic",
		"Parody",
		"Thriller",
		"Music",
		"Game",
		"Fantasy",
		"Adventure",
		"Romance",
		"Police",
		"Drama",
		"Samurai",
		"School",
		"Comedy",
		"Shoujo",
		"Military",
		"Horror",
		"Slice of Life",
		"Mecha",
	}

	durations := []DurationCategory{
		{ID: "very-short", Label: "очень короткая"},
		{ID: "short", Label: "короткая"},
		{ID: "not-long", Label: "не долгая"},
		{ID: "not-very-long", Label: "не очень долгая"},
		{ID: "not-very-short", Label: "не очень короткая"},
		{ID: "not-short", Label: "не короткая"},
		{ID: "long", Label: "долгая"},
		{ID: "very-long", Label: "очень долгая"},
	}

	genresOpts := make([]GenreOption, 0, len(genres))
	for _, genre := range genres {
		id := "genre-" + strings.ReplaceAll(strings.ToLower(genre), " ", "-")
		genresOpts = append(genresOpts, GenreOption{
			ID:    id,
			Label: genre,
		})
	}

	typesOpts := make([]TypeOption, 0, len(filter.AnimeTypes))
	for _, animeType := range filter.AnimeTypes {
		id := "type-" + strings.ReplaceAll(strings.ToLower(animeType), " ", "-")
		typesOpts = append(typesOpts, TypeOption{
			ID:    id,
			Label: animeType,
		})
	}

	return FilterParams{
		GenreOptions:       genresOpts,
		DurationCategories: durations,
		TypeOptions:        typesOpts,
	}
}
