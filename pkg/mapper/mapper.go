package mapper

import (
	"github.com/jinzhu/copier"
)

// AutoMap automatically maps between structs with similar field names
// Uses reflection to copy fields from source to destination type
func AutoMap[T any](src interface{}) (*T, error) {
	var dst T
	err := copier.Copy(&dst, src)
	if err != nil {
		return nil, err
	}
	return &dst, nil
}

// MapSlice maps a slice of structs to another slice type
func MapSlice[T any](src interface{}) ([]T, error) {
	var dst []T
	err := copier.Copy(&dst, src)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

// MapTo maps from source to destination directly (for existing instances)
func MapTo(dst, src interface{}) error {
	return copier.Copy(dst, src)
}
