package model

import "github.com/satori/go.uuid"

type Books struct {
	UUIDBooks        uuid.UUID
	BooksName        string
	BooksPublisher   string
	BooksWriter      string
	BooksDescription string
}
