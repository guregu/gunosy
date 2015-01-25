package daihinmin

type X string

type WelcomeReply struct {
	X
	Name string
	Sesh sesh
}

type ErrorReply struct {
	X
	Wrt string
	Msg string
}
