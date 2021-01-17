package domain

type Server interface {
	Scheme() string
	ListenAndServe()
}
