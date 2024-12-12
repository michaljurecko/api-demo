package webapi

import (
	"encoding/json"
	"errors"
	"strconv"
)

// Lookup is generic type to represents relations between entities.
// T is type of the referenced entity.
type Lookup[T any] struct {
	// id of the referenced entity, e.g., 11bb11bb-cc22-dd33-ee44-55ff55ff55ff
	id string
	// contentID is used if the entity was not saved yet, and it is created in the same ChangeSet as the target entity.
	// See:
	// https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/execute-batch-operations-using-web-api#reference-uris-in-request-body
	contentID int
}

type LookupInterface interface {
	ID() string
	ContentIDRef() *string
}

func LookupValue[T any](id string) Lookup[T] {
	return Lookup[T]{id: id}
}

func (l Lookup[T]) ID() string {
	if l.contentID > 0 {
		panic(errors.New("ID is not know, entity was not saved yet"))
	}
	return l.id
}

func (l Lookup[T]) ContentIDRef() *string {
	if l.contentID > 0 {
		ref := "$" + strconv.Itoa(l.contentID)
		return &ref
	}
	return nil
}

func (l *Lookup[T]) SetID(id string) {
	if l.contentID > 0 {
		panic(errors.New("cannot set 'id', 'contentIDRef' is set"))
	}
	l.id = id
}

func (l *Lookup[T]) SetContentID(id int) {
	if l.id != "" {
		panic(errors.New("cannot set 'contentIDRef', 'id' is set"))
	}
	l.contentID = id
}

func (l *Lookup[T]) Clear() {
	l.id = ""
	l.contentID = 0
}

// MarshalJSON implements the Marshaler interface for Lookup[T].
func (l Lookup[T]) MarshalJSON() ([]byte, error) {
	if l.contentID > 0 {
		return json.Marshal(l.ContentIDRef())
	}

	return json.Marshal(l.id)
}

// UnmarshalJSON implements the Unmarshaler interface for Lookup[T].
func (l *Lookup[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &l.id)
}
