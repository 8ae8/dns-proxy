package domain

type ServerModel interface {
	ListenAndServe() error
}
