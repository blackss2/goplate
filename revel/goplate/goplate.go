package goplate

import (
	"fmt"
	"github.com/blackss2/goplate"
	"github.com/revel/revel"
	"os"
	"path/filepath"
	"strings"
)

var (
	MainGoplateLoader *GoplateLoader
	MainWatcher       *revel.Watcher
)

type GoplateLoader struct {
	*goplate.PlateLoader
	paths []string
}

func NewGoplateLoader() *GoplateLoader {
	return &GoplateLoader{
		paths: make([]string, 0, 1),
	}
}

func mkDirIfNotExist(path string) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(path, 0777)
		}
	}
}

func (this *GoplateLoader) Refresh() *revel.Error {
	this.PlateLoader = goplate.NewPlateLoader()
	viewsPath := revel.TemplatePaths[0]

	for _, rootPath := range MainGoplateLoader.paths {
		filepath.Walk(rootPath, func(path string, f os.FileInfo, err error) error {
			if filepath.Ext(path) == ".htm" || filepath.Ext(path) == ".html" {
				this.LoadFile(path)
				relPath, err := filepath.Rel(rootPath, path)
				if err != nil {
					panic(err)
				}
				if strings.HasPrefix(relPath, "views\\") || strings.HasPrefix(relPath, "views/") {
					outputPath := fmt.Sprintf("%s/%s", viewsPath, relPath[5:])
					mkDirIfNotExist(filepath.Dir(outputPath))
					file, err := os.Create(outputPath)
					if err != nil {
						panic(err)
					}
					if file != nil {
						result := this.ApplyFile(path)
						file.WriteString(result)
						file.Close()
					}
				}
			}
			return nil
		})
	}
	return nil
}

func (this *GoplateLoader) WatchDir(info os.FileInfo) bool {
	return !strings.HasPrefix(info.Name(), ".")
}

func (this *GoplateLoader) WatchFile(basename string) bool {
	return !strings.HasPrefix(basename, ".")
}

var revelWatchFilter func(c *revel.Controller, fc []revel.Filter)
var WatchFilter = func(c *revel.Controller, fc []revel.Filter) {
	if MainWatcher != nil {
		err := MainWatcher.Notify()
		if err != nil {
			c.Result = c.RenderError(err)
			return
		}
	}
	revelWatchFilter(c, fc)
}

func init() {
	revel.OnAppStart(func() {
		MainGoplateLoader = NewGoplateLoader()
		viewsPath := revel.TemplatePaths[0]
		appPath := fmt.Sprintf("%s/goplates", filepath.Dir(viewsPath))
		MainGoplateLoader.paths = append(MainGoplateLoader.paths, appPath)
		revel.TemplateDelims = "[[ ]]"

		if revel.Config.BoolDefault("watch", true) {
			MainWatcher = revel.NewWatcher()
			revelWatchFilter = revel.WatchFilter
			revel.WatchFilter = WatchFilter
			revel.Filters = append([]revel.Filter{WatchFilter}, revel.Filters...)
		}
		if MainWatcher != nil && revel.Config.BoolDefault("watch.templates", true) {
			MainWatcher.Listen(MainGoplateLoader, MainGoplateLoader.paths...)
		}
	})
}
