package exception

type Exception interface {
	Error() string
}

type MessageException struct {
	message string
}

func GlobalMessageException(message string) *MessageException {
	return &MessageException{
		message: message,
	}
}
