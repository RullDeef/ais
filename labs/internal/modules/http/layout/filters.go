package layout

import "strings"

type GenreOption struct {
	ID       string
	Label    string
	Selected string
}

type FilterParams struct {
	GenreOptions []GenreOption
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

	genresOpts := make([]GenreOption, 0, len(genres))
	for _, genre := range genres {
		id := "genre-" + strings.ReplaceAll(strings.ToLower(genre), " ", "-")
		genresOpts = append(genresOpts, GenreOption{
			ID:       id,
			Label:    genre,
			Selected: "",
		})
	}

	return FilterParams{
		GenreOptions: genresOpts,
	}
}
