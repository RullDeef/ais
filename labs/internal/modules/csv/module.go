package csv

import (
	"anicomend/model"

	"go.uber.org/fx"
)

var Module = fx.Module("csv",
	fx.Provide(
		fx.Annotate(
			NewAnimeLoader,
			fx.ParamTags(`name:"dataset-path"`),
		),
		func(loader *CSVAnimeLoader) model.AnimeLoader { return loader },
	),
)
