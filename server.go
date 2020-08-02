package namespace

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/server"
)

type User interface {
	Namespaces() (personal, other, shared []Namespace, err error)
}

type Handler struct {
	Command
}

func (h *Handler) Handle(conn server.Conn) error {
	if conn.Context().User == nil {
		return server.ErrNotAuthenticated
	}

	resp := Response{}

	usr, ok := conn.Context().User.(User)
	if ok {
		var err error
		resp.Personal, resp.Other, resp.Shared, err = usr.Namespaces()
		if err != nil {
			return err
		}
		return conn.WriteResp(&resp)
	}

	info, err := conn.Context().User.ListMailboxes(false)
	if err != nil {
		return err
	}

	if len(info) != 0 {
		resp.Personal = []Namespace{{
			Prefix:    "",
			Delimiter: info[0].Delimiter,
		}}
	}

	return conn.WriteResp(&resp)
}

type extension struct{}

func NewExtension() server.Extension {
	return &extension{}
}

func (ext *extension) Capabilities(c server.Conn) []string {
	if c.Context().State&imap.AuthenticatedState == 0 {
		return nil
	}

	return []string{"NAMESPACE"}
}

func (ext *extension) Command(name string) server.HandlerFactory {
	if name == "NAMESPACE" {
		return func() server.Handler {
			return &Handler{}
		}
	}
	return nil
}
