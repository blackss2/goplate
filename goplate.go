package goplate

import (
	"net/http"
	"github.com/revel/revel"
)

type MyHtml string

func (r MyHtml) Apply(req *revel.Request, resp *revel.Response) {
	resp.WriteHeader(http.StatusOK, "text/html")
	resp.Out.Write([]byte(r))
}
