package exception

// Global Exception
func GlobalException(message error) {
	if message != nil {
		panic(message)
	}
}
