package model

import (
	"fmt"
	"slices"
)

type GenreSimilaritier struct {
	genreList  []string
	genrelinks []int
	distances  [][]int
}

var (
	genreList []string = []string{
		"Action",        // 0
		"Adventure",     // 1
		"Cars",          // 2
		"Comedy",        // 3
		"Demons",        // 4
		"Drama",         // 5
		"Fantasy",       // 6
		"Game",          // 7
		"Historical",    // 8
		"Horror",        // 9
		"Josei",         // 10
		"Magic",         // 11
		"Martial Arts",  // 12
		"Mecha",         // 13
		"Military",      // 14
		"Music",         // 15
		"Mystery",       // 16
		"Parody",        // 17
		"Police",        // 18
		"Psychological", // 19
		"Romance",       // 20
		"Samurai",       // 21
		"School",        // 22
		"Sci-Fi",        // 23
		"Seinen",        // 24
		"Shoujo",        // 25
		"Shounen",       // 26
		"Slice of Life", // 27
		"Space",         // 28
		"Sports",        // 29
		"Supernatural",  // 30
		"Super Power",   // 31
		"Thriller",      // 32
		"Vampire",       // 33
	}
	genreLinks []int = []int{
		26, // "Action", -> shounen
		26, // "Adventure", -> shounen
		29, // "Cars", -> sports
		27, // "Comedy", -> sliceoflife
		30, // "Demons", -> supernatural
		27, // "Drama", -> sliceoflife
		16, // "Fantasy", -> mystery
		26, // "Game", -> shounen
		0,  // "Historical", -> action
		32, // "Horror", -> thriller
		25, // "Josei", -> Shoujo
		6,  // "Magic", -> fantasy
		0,  // "Martial Arts", -> action
		23, // "Mecha", -> scifi
		18, // "Military", -> police
		27, // "Music", -> sliceoflife
		-1, // "Mystery",
		27, // "Parody", -> sliceoflife
		0,  // "Police", -> action
		24, // "Psychological", -> seinen
		25, // "Romance", -> Shoujo
		8,  // "Samurai", -> historical
		27, // "School", -> sliceoflife
		-1, // "Sci-Fi",
		-1, // "Seinen",
		-1, // "Shoujo",
		-1, // "Shounen",
		-1, // "Slice of Life",
		23, // "Space", -> scifi
		26, // "Sports", -> shounen
		16, // "Supernatural", -> mystery
		30, // "Super Power", -> supernatural
		16, // "Thriller", -> mystery
		30, // "Vampire", -> supernatural
	}
)

func NewGenreSimilaritier() *GenreSimilaritier {
	gs := GenreSimilaritier{
		genreList:  genreList,
		genrelinks: genreLinks,
	}

	var buildPath func(int) []int
	buildPath = func(g int) []int {
		if genreLinks[g] == -1 {
			return []int{g}
		} else {
			return append(buildPath(genreLinks[g]), g)
		}
	}

	symDiffLen := func(s1, s2 []int) int {
		diff := 0
		for _, s := range s1 {
			if !slices.Contains(s2, s) {
				diff++
			}
		}
		for _, s := range s2 {
			if !slices.Contains(s1, s) {
				diff++
			}
		}
		return diff
	}

	// recalc distances
	gs.distances = make([][]int, len(genreList))
	for i := 0; i < len(genreList)-1; i++ {
		gs.distances[i] = make([]int, len(genreList))
		pi := buildPath(i)
		for j := 0; j < i; j++ {
			gs.distances[j][i] = gs.distances[i][j]
		}
		for j := i + 1; j < len(genreList); j++ {
			pj := buildPath(j)
			gs.distances[i][j] = symDiffLen(pi, pj)
		}
	}

	return &gs
}

func (gs *GenreSimilaritier) Distance(g1, g2 string) int {
	i1 := slices.Index(gs.genreList, g1)
	i2 := slices.Index(gs.genreList, g2)

	if i1 == -1 {
		panic(fmt.Errorf("genre %s not found", g1))
	}
	if i2 == -1 {
		panic(fmt.Errorf("genre %s not found", g2))
	}

	return gs.distances[i1][i2]
}

func (gs *GenreSimilaritier) Similarity(genres1, genres2 []string) float64 {
	jac := float64(set_intersect_len(genres1, genres2)) / float64(set_union_len(genres1, genres2))
	gs1, gs2 := set_diff(genres1, genres2), set_diff(genres2, genres1)

	if len(gs1) == 0 {
		gs1 = genres1
	}
	if len(gs2) == 0 {
		gs2 = genres2
	}

	sim, max_dist := 0.0, 1.0
	for _, g1 := range gs1 {
		for _, g2 := range gs2 {
			dist := float64(gs.Distance(g1, g2))
			max_dist = max(max_dist, dist)
			sim -= dist
		}
	}
	sim += max_dist * float64(len(gs1)*len(gs2))
	sim /= float64(len(gs1) * len(gs2))
	return min(1, jac+sim/6)
}

func set_diff(s1, s2 []string) []string {
	res := make([]string, 0, len(s1))
	for _, s := range s1 {
		if !slices.Contains(s2, s) {
			res = append(res, s)
		}
	}
	return res
}

func set_intersect_len(s1, s2 []string) int {
	res := 0
	for _, s := range s1 {
		if slices.Contains(s2, s) {
			res++
		}
	}
	return res
}

func set_union_len(s1, s2 []string) int {
	res := len(s1)
	for _, s := range s2 {
		if !slices.Contains(s1, s) {
			res++
		}
	}
	return res
}
