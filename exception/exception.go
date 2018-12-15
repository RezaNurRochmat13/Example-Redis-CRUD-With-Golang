package exception

type Exception interface {
	Error() string
}

type MessageException struct {
	message string
}

func GlobalMessageException(message error) *MessageException {
	return &MessageException{
		message: message,
	}
}
