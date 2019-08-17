package auth

type Auth interface {
	Validate(string) bool
}
