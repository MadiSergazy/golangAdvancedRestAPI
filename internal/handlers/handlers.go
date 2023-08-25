package handlers

import "github.com/julienschmidt/httprouter"

// all of the entities in the project that have handlers have to implement this Method of interface
type Handler interface {
	Register(router *httprouter.Router)
}
