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

type YouJoined struct {
	X
	Chan string
}

type UserJoinPartReply struct {
	X
	Chan string
	User string
}

type GameInfo struct {
	X     `json:",omitempty"`
	ID    string
	Name  string
	Users []string
	Owner string
}
