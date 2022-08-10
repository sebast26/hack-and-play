package es

type envelope struct {
	Event   string
	Version string
	Data    interface{}
}
