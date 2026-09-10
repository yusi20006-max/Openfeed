package provider

type Provider interface {
	LoadChannel(name string) ([]byte, error)
	ClientFetchHTMLBefore(name string, beforeID int) ([]byte, error)
}

var Default Provider
