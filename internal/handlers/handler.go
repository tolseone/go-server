package handlers

import "github.com/julienschmidt/httprouter"

type Handler interface {
	Reqister(router *httprouter.Router)
}
