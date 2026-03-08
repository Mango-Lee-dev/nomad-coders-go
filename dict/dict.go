package dict

import "errors"

//	Dictionary type
type Dictionary map[string]string

var errNotFound = errors.New("Not found")

// Search for a word in the dictionary
func (d Dictionary) Search(word string) (string, error) {
	value, exists := d[word]
	if exists {
		return value, nil
	}
	return "", errNotFound
}