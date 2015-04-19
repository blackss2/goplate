package goplate

import (
	"fmt"
	"github.com/blackss2/goplate"
	"github.com/revel/revel"
	"net/http"
	"html/template"
	"os"
	"path/filepath"
	"time"
)

var (
	MainRenderer *Renderer
)

type Renderer struct {
	loader *goplate.PlateLoader
	templateSet *template.Template
	isInit bool
}

func Render(c *revel.Controller, v ...interface{}) revel.Result {
	MainRenderer.Refresh()
	/*
	filepath := revel.TemplatePaths[0] + "/" + c.Name + "/" + c.MethodName + c.Request.Format
	result := MainRenderer.loader.ApplyFile(filepath)
	return c.RenderHtml(result)
	*/
	c.Response.Status = http.StatusOK

	template, err := MainRenderer.templateSet.Template(templatePath)
	if err != nil {
		return c.RenderError(err)
	}

	return &RenderTemplateResult{
		Template:   template,
		RenderArgs: c.RenderArgs,
	}
}

func (this *Renderer) Refresh() {
	//TODO : mutex lock MainRenderer
	if !MainRenderer.isInit {
		this.loader = goplate.NewPlateLoader()
		if len(revel.TemplatePaths) > 0 {
			path := revel.TemplatePaths[0] //TEMP
			rootPath := fmt.Sprintf("%s/goplates", filepath.Dir(path))
			filepath.Walk(rootPath, func(path string, f os.FileInfo, err error) error {
				if filepath.Ext(path) == ".htm" || filepath.Ext(path) == ".html" {
					this.loader.LoadFile(path)
					result := MainRenderer.loader.ApplyFile(path)
					MainRenderer.templateSet.Parse(result)
				}
				return nil
			})
			MainRenderer.isInit = true
		}
	}
}

func init() {
	MainRenderer = &Renderer{
		templateSet: template.New(templateName).Funcs(TemplateFuncs),
	}
	go func() {
		tick := time.Tick(5 * time.Second)
		select {
		case <-tick:
			MainRenderer.isInit = false
			MainRenderer.Refresh() //TEMP
		}
	}()
}
