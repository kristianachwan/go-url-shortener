package shortener

type RedirectService interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
	Delete(code string) error
}
