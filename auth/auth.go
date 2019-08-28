package auth

//Auth interface for validation struct
type Auth interface {
	Validate(string) bool
}
