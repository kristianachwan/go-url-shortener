package shortener

type RedirectRepository interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
	Delete(code string) error
}
