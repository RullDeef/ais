package csv

import (
	"anicomend/model"

	"go.uber.org/fx"
)

var Module = fx.Module("csv",
	fx.Provide(
		fx.Annotate(
			func(filename string) model.AnimeLoader {
				return NewAnimeLoader(filename)
			},
			fx.ParamTags(`name:"dataset-path"`),
		),
	),
)
