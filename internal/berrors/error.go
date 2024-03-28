package berrors

type Berror interface {
	Code() string
	Msg() string
}

type berrorStruct struct {
	c string
	s string
}

func New(code string, msg string) Berror {
	return &berrorStruct{
		c: code,
		s: msg,
	}
}

func (e *berrorStruct) Code() string {
	return e.c
}

func (e *berrorStruct) Msg() string {
	return e.s
}
