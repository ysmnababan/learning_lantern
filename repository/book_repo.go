package repository

type BookRepo interface {
	GetAllBooks() error
}

func (r *Repo) GetAllBooks() error {
	return nil
}
