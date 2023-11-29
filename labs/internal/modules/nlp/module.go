package nlp

import "go.uber.org/fx"

var Module = fx.Module("nlp",
	fx.Provide(
		fx.Annotate(
			New,
			fx.ParamTags(`name:"nlp-base-url"`),
		),
	),
)
