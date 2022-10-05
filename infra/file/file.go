package file

import (
	"errors"
	"io/ioutil"
)

func ReadFile(name string) (data []byte, err error) {
	if len(name) == 0 {
		return nil, errors.New(`name of file is empty`)
	}

	return ioutil.ReadFile(name)
}
