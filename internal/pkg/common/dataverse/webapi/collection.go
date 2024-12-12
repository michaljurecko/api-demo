package webapi

type Collection[T any] struct {
	Items []*T `json:"value"`
}
