package fixture // OMIT

// panics
// +getters
type ()

type (
	// +getters
	First struct{ value string }

	// ignored
	// +getters
	Second struct{ value string }
)
