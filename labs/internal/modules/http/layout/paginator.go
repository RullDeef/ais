package layout

import "strconv"

type Page struct {
	Index     string
	Current   bool
	Elipsised bool
}

// 1 2 [3] 4 (total=4, max>total)
// 1 2 [3] 4 5 6 7 8 ... (max < total && curr <= (max + 1)/2)
// ... 7 8 9 [10] 11 12 13 ... (max < total && curr > (max + 1)/2)
// ... 7 8 9 [10] 11 12 13 14 (max < total && curr > (max + 1)/2 && curr + (max + 1)/2 >= total)
func FormatPages(totalPages int, maxPages int, currentPage int) []Page {
	if totalPages <= maxPages {
		pages := make([]Page, maxPages)
		for i := 1; i <= totalPages; i++ {
			pages[i-1] = Page{
				Index:   strconv.FormatInt(int64(i), 10),
				Current: i == int(currentPage),
			}
		}
		return pages
	} else {
		elipsisedStart := currentPage > (maxPages+1)/2
		elipsisedEnd := currentPage+(maxPages+1)/2 < totalPages+2
		startPage := min(currentPage-(maxPages+1)/2, totalPages-maxPages+1)
		startPage = max(1, startPage)
		pages := make([]Page, maxPages)
		for i := 0; i < maxPages; i++ {
			pages[i] = Page{
				Index:   strconv.FormatInt(int64(startPage+i), 10),
				Current: startPage+i == int(currentPage),
			}
		}
		if elipsisedStart {
			pages[0].Index = "..."
			pages[0].Elipsised = true
		}
		if elipsisedEnd {
			pages[maxPages-1].Index = "..."
			pages[maxPages-1].Elipsised = true
		}
		return pages
	}
}
