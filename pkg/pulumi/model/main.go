package model

type Metadata struct {
	Name    string
	Team    string
	Env     string
	Cloud   string
	Account string
	Region  string
	Info    map[string]string
}
