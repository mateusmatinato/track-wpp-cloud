package commands

type ProcessRequest struct {
	Message string
	Contact Contact
}

type Contact struct {
	Number string
	Name   string
}
