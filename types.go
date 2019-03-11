package gin_restful

type Resource struct {
	Prefix string
}

type ApplicationError struct {
	Message string
	Status  int
}

func (e ApplicationError) Error() string {
	return e.Message
}
