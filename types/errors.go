package types

import (
	"errors"
)

var (
	ErrFileNotFound = errors.New("file not found")
)

//type FileNotFound string
//
//func (e FileNotFound) Error() string {
//	return fmt.Sprintf("File not found: %s", (string)(e))
//}
