package lib

type Route struct {
	Method  string
	Pattern string
	Handler Handler
}

type Routes []Route
