package gosumo

type ErrPostingLogs struct {
	Message string
}

func (e ErrPostingLogs) Error() string {
	return e.Message
}

type ErrParsingLogs struct {
	Message string
}

func (e ErrParsingLogs) Error() string {
	return e.Message
}
