package mydict

import "errors"

// Doctionary map
type Dictionary map[string]string

var errNotFound = errors.New("not found")
var errWordExists = errors.New("that word already exists")
var errCantUpdate = errors.New("can not update, non-existing word")

// Search for a word
func (d Dictionary) Search(word string) (string, error) {
	v, exists := d[word]
	if exists {
		return v, nil
	} else {
		return "", errNotFound
	}
}

// Add a word to dictionary
func (d Dictionary) Add(word, def string) error {
	_, err := d.Search(word)

	switch err {
	case errNotFound:
		d[word] = def
	case nil:
		return errWordExists
	}

	return nil
}

// Update definition of a word
func (d Dictionary) Update(word, def string) error {
	_, err := d.Search(word)

	switch err {
	case nil:
		d[word] = def
	case errNotFound:
		return errCantUpdate
	}

	return nil
}

// Delete a word
func (d Dictionary) Delete(word string) {
	delete(d, word)
}
