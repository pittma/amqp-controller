package controller

type NotBoundError struct {
    name string
}

func (e NotBoundError) Error() string {
    return e.name
}

type NotConnectedError struct {
    name string
}

func (e NotConnectedError) Error() string {
    return e.name
}
