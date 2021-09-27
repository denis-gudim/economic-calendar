package data

import "fmt"

type LanguagesRepository struct {
	baseRepository
}

func (r *LanguagesRepository) GetAll() (languages []Language, err error) {
	db, err := r.createConnection()

	if err != nil {
		return nil, fmt.Errorf("get all languages create connection error: %w", err)
	}

	defer db.Close()

	err = db.Select(&languages, "SELECT code, name, native_name FROM languages")

	if err != nil {
		return nil, fmt.Errorf("get all languages select error: %w", err)
	}

	return
}
