package mock_base_types

import "context"

type (
	SecurityString interface{ SecurityString() string }
	String         interface{ String() string }
	Named          interface{ Name() string }

	DefaultSetter            interface{ SetDefault() }
	Initializer              interface{ Init() }
	ValidatedInitializer     interface{ Init() error }
	InitializerWith          interface{ Init(context.Context) }
	ValidatedInitializerWith interface{ Init(context.Context) error }
)
