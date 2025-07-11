package domain

type URLService interface {
	Shorten(originalURL string) (string, error)
	Redirect(shortURL string) (string, error)
}
