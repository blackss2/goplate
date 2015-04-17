package goplate

import (
	"fmt"
	"github.com/blackss2/goplate"
	"github.com/revel/revel"
	"path/filepath"
	"os"
)

var (
	MainRenderer *Renderer
)

type Renderer struct {
	loader *goplate.PlateLoader
}

func Render(c *revel.Controller, v ...interface{}) revel.Result {
	MainRenderer.Refresh()
	filepath := c.Request.RequestURI
	result := MainRenderer.loader.ApplyFile(filepath)
	return c.RenderHtml(result)
}

func (this *Renderer) Refresh() {
	this.loader = goplate.NewPlateLoader()
	path := revel.TemplatePaths[0] //TEMP
	rootPath := fmt.Sprintf("%s/goplates", filepath.Dir(path))
	filepath.Walk(rootPath, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == ".htm" || filepath.Ext(path) == ".html" {
			this.loader.LoadFile(path)
		}
		return nil
	})
}

func init() {
	MainRenderer = &Renderer{}
}