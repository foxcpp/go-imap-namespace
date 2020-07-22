package namespace

import (
	"errors"

	"github.com/emersion/go-imap"
)

type Command struct{}

func (cmd Command) Command() *imap.Command {
	return &imap.Command{
		Name:      "NAMESPACE",
		Arguments: nil,
	}
}

func (cmd Command) Parse(fields []interface{}) error {
	if len(fields) != 0 {
		return errors.New("No arguments expected")
	}
	return nil
}
