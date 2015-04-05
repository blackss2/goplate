package goplate

import (
	"fmt"
	"net/http"
	"github.com/revel/revel"
)

type Renderer struct {
	controller *revel.Controller
	v []interface{}
}

func (r Renderer) Apply(req *revel.Request, resp *revel.Response) {
	/*
	//TODO
	1. Load html from "www" folter
	2. If [path].js is exist, add this to scriptsList
	3. If [path].css is exist, add this to cssList
	4. Append scriptList, cssList to end of body
	*/
	resp.WriteHeader(http.StatusOK, "text/html")
	resp.Out.Write([]byte(fmt.Sprintln(r.controller.Name + "/" + r.controller.MethodType.Name + "." + r.controller.Request.Format)))
}

func Render(c *revel.Controller, v ...interface{}) revel.Result {
	r := new(Renderer)
	r.controller = c
	r.v = v
	return r
}
