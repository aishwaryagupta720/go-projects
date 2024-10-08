package server

import (
	"net/http"

	"elahi-arman.github.com/example-http-server/internal/datastore"
	"github.com/julienschmidt/httprouter"
)

// The initialisation of Server and its methods
type Server interface {
	GetLink(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	CreateLink(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

// Implementation of the Server which was defined in the interface above
type serverImpl struct {
	// database to store links
	linkStore datastore.LinkStorer
}

var _ Server = (*serverImpl)(nil)

func NewServer(ls datastore.LinkStorer) *serverImpl {
	return &serverImpl{
		linkStore: ls,
	}
}
