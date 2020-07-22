package namespace

import (
	"errors"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/responses"
)

type Response struct {
	Personal []Namespace
	Other    []Namespace
	Shared   []Namespace
}

const responseName = "NAMESPACE"

func (r *Response) Handle(resp imap.Resp) error {
	name, fields, ok := imap.ParseNamedResp(resp)
	if !ok || name != responseName {
		return responses.ErrUnhandled
	}

	return r.Parse(fields)
}

func (r *Response) Parse(fields []interface{}) error {
	if len(fields) != 3 {
		return errors.New("namespace: expected three response parts")
	}

	for i, f := range fields {
		var res []Namespace
		if f == nil {
			continue
		}

		list, ok := f.([]interface{})
		if !ok {
			return errors.New("namespace: expected list")
		}
		if len(list) == 0 {
			return errors.New("namespace: empty namespace list not allowed")
		}
		for _, ns := range list {
			nsList, ok := ns.([]interface{})
			if !ok {
				return errors.New("namespace: expected list for namespace")
			}
			if len(nsList) < 2 {
				return errors.New("namespace: missing namespace arguments")
			}
			prefix, ok := nsList[0].(string)
			if !ok {
				return errors.New("namespace: prefix should be a string")
			}
			delimiter, ok := nsList[1].(string)
			if !ok {
				return errors.New("namespace: delimiter should be a string")
			}

			res = append(res, Namespace{Prefix: prefix, Delimiter: delimiter})
		}

		switch i {
		case 0:
			r.Personal = res
		case 1:
			r.Other = res
		case 2:
			r.Shared = res
		}
	}

	return nil
}

func (r Response) WriteTo(w *imap.Writer) error {
	fields := make([]interface{}, 0, 4)
	fields = append(fields, imap.RawString("NAMESPACE"))
	for _, nsSet := range [...][]Namespace{r.Personal, r.Other, r.Shared} {
		if len(nsSet) == 0 {
			fields = append(fields, nil)
			continue
		}
		setFields := make([]interface{}, 0, len(nsSet))
		for _, ns := range nsSet {
			setFields = append(setFields, []interface{}{
				ns.Prefix,
				ns.Delimiter,
			})
		}
		fields = append(fields, setFields)
	}
	return imap.NewUntaggedResp(fields).WriteTo(w)
}
