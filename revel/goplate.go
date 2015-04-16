package goplate

import (
	"fmt"
	"net/http"
	"github.com/revel/revel"
	"github.com/blackss2/goplate"
	"path/filepath"
)

type Renderer struct {
	controller *revel.Controller
	v []interface{}
}

var (
	MainRenderer *Renderer
)

func (r Renderer) Apply(req *revel.Request, resp *revel.Response) {
	/*
	//TODO
	1. Load html from "www" folter
	*/
	//
	resp.WriteHeader(http.StatusOK, "text/html")
	resp.Out.Write([]byte(fmt.Sprintln(r.controller.Name + "/" + r.controller.MethodType.Name + "." + r.controller.Request.Format)))
}

func Render(c *revel.Controller, v ...interface{}) revel.Result {
	r := new(Renderer)
	r.controller = c
	r.v = v
	return r
}

func init() {
	
}