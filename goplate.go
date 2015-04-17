package goplate

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/blackss2/gohtml"
	"github.com/ditashi/jsbeautifier-go/jsbeautifier"
	"github.com/gorilla/css/scanner"
	//"html/template"
	"io"
	"os"
	"strings"
)

type PlateLoader struct {
	plateHash map[string]*Plate
	ctrlHash  map[string]*Controller
}

func NewPlateLoader() *PlateLoader {
	return &PlateLoader{
		plateHash: make(map[string]*Plate),
		ctrlHash:  make(map[string]*Controller),
	}
}

func (this *PlateLoader) LoadFile(filepath string) {
	docFile, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	this.Load(docFile)
}

func (this *PlateLoader) LoadString(src string) {
	var srcBuffer bytes.Buffer
	srcBuffer.WriteString(src)
	this.Load(&srcBuffer)
}

func (this *PlateLoader) Load(src io.Reader) {
	doc, err := goquery.NewDocumentFromReader(src)
	if err != nil {
		panic(err)
	}

	if doc.Find("plate > script").Length() > 0 {
		panic("plate not allowed exist script tag in plate's body")
	}
	doc.Find("plate").Each(func(idx int, jPlate *goquery.Selection) {
		plate := NewPlate()
		if name, has := jPlate.Attr("name"); !has {
			panic("plate name missing")
		} else {
			plate.Name = name
		}
		if _, has := this.plateHash[plate.Name]; has {
			panic("plate collision")
		}
		this.plateHash[plate.Name] = plate
		plate.jNode = jPlate
		jPlate.Remove()

		plate.jNode.Find("css").Each(func(idx int, jCss *goquery.Selection) {
			css := NewCss(int64(len(plate.cssList) + 1))
			plate.cssList = append(plate.cssList, css)
			css.Parse(jCss.Text())
			jCss.Remove()

			if inherit, has := jCss.Attr("inherit"); has && inherit == "true" {
				css.inherit = true
			}
		})
		for _, css := range plate.cssList {
			if !css.inherit {
				css.Apply(plate.jNode)
			}
		}
	})
}

func (this *PlateLoader) ApplyFile(filepath string) string {
	docFile, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	return this.Apply(docFile)
}

func (this *PlateLoader) ApplyString(src string) string {
	var srcBuffer bytes.Buffer
	srcBuffer.WriteString(src)
	return this.Apply(&srcBuffer)
}

func (this *PlateLoader) Apply(src io.Reader) string {
	doc, err := goquery.NewDocumentFromReader(src)
	if err != nil {
		panic(err)
	}

	for _, p := range this.plateHash {
		this.replacePlate(p, doc.Find("body"))
	}

	htmlStr, err := doc.Html()
	if err != nil {
		panic(err)
	}
	var htmlBuffer bytes.Buffer
	htmlBuffer.WriteString(htmlStr)
	htmlBuffer.WriteString("<script>")
	options := jsbeautifier.DefaultOptions()
	for _, ctrl := range this.ctrlHash {
		jsStr := ctrl.String()
		jsFormatted, err := jsbeautifier.Beautify(&jsStr, options)
		if err != nil {
			panic(err)
		}
		htmlBuffer.WriteString(jsFormatted)
		htmlBuffer.WriteString("\n")
	}
	htmlBuffer.WriteString("</script>")
	htmlBuffer.WriteString("<style>")
	for _, plate := range this.plateHash {
		for _, css := range plate.cssList {
			htmlBuffer.WriteString(css.String())
			htmlBuffer.WriteString("\n")
		}
	}
	htmlBuffer.WriteString("</style>")
	return gohtml.Format(htmlBuffer.String())
}

type Plate struct {
	Name    string
	jNode   *goquery.Selection
	cssList []*Css
}

func NewPlate() *Plate {
	return &Plate{
		cssList: make([]*Css, 0),
	}
}

func (this *PlateLoader) replacePlate(plate *Plate, jTarget *goquery.Selection) {
	jTarget.Find(plate.Name).Each(func(idx int, jPlate *goquery.Selection) {
		/*
			argList := make([]string, 0)
			this.Find("arg").Each(func(idx int, arg *goquery.Selection) {
				argList = append(argList, arg.Text())
				arg.Remove()
			})
		*/

		jClone := plate.jNode.Clone()
		for _, p := range this.plateHash {
			if p.Name != plate.Name {
				this.replacePlate(p, jClone)
			}
		}
		/*
			htmlStr, err := clone.Html()
			if err != nil {
				panic(err)
			}
			tp, err := template.New(p.Name).Delims("[[", "]]").Parse(htmlStr)
			if err != nil {
				panic(err)
			}
			var buffer bytes.Buffer
			tp.Execute(&buffer, &map[string]interface{}{
				"args": argList,
			})

			executedHtml := buffer.String()
			clone.ReplaceWithHtml(executedHtml)
		*/
		for _, attr := range jPlate.Nodes[0].Attr {
			jClone.Children().SetAttr(attr.Key, attr.Val)
		}
		jPlate.Children().Each(func(idx int, jChild *goquery.Selection) {
			jClone.Children().AppendSelection(jChild)
		})

		for _, css := range plate.cssList {
			if css.inherit {
				css.Apply(jClone)
			}
		}
		jClone.Find("script").Each(func(idx int, jScript *goquery.Selection) {
			jParent := jScript.Parent()
			jScript.Remove()

			var ctrl *Controller
			if ctrlName, has := jParent.Attr("ng-controller"); !has {
				ctrl = NewController(int64(len(this.ctrlHash) + 1)) //TEMP
				ctrlName = fmt.Sprintf("Ctrl_%d", ctrl.Id)
				jParent.SetAttr("ng-controller", ctrlName)
				this.ctrlHash[ctrlName] = ctrl
			} else {
				ctrl = this.ctrlHash[ctrlName]
			}

			if injectStr, has := jScript.Attr("inject"); has {
				injectList := strings.Split(injectStr, ",")
				for _, v := range injectList {
					inject := strings.TrimSpace(v)
					ctrl.InjectHash[inject] = true
				}
			}

			event := ctrl.NewEventHandler()
			if eventStr, has := jScript.Attr("event"); has {
				idx := strings.Index(eventStr, "(")
				if idx >= 0 {
					event.Type = eventStr[:idx]
					event.Args = eventStr[idx:]
					ngEventName := "ng-" + event.Type
					eventHandlerStr := fmt.Sprintf("EventHandler_%d_%d%s", ctrl.Id, event.Id, event.Args)
					if attrStr, has := jParent.Attr(ngEventName); has {
						jParent.SetAttr(ngEventName, fmt.Sprintf("%s; %s", attrStr, eventHandlerStr))
					} else {
						jParent.SetAttr(ngEventName, eventHandlerStr)
					}
				} else {
					event.Type = eventStr
				}
			}
			event.Body = jScript.Text()
		})

		jPlate.ReplaceWithSelection(jClone.Children())
	})
}

type Controller struct {
	Id            int64
	EventHandlers []*EventHandler
	InjectHash    map[string]bool
}

type EventHandler struct {
	Id   int64
	Type string
	Body string
	Args string
}

func NewController(Id int64) *Controller {
	return &Controller{
		Id:            Id,
		EventHandlers: make([]*EventHandler, 0),
		InjectHash:    make(map[string]bool),
	}
}

func (this *Controller) NewEventHandler() *EventHandler {
	event := &EventHandler{
		Id: int64(len(this.EventHandlers) + 1), //TEMP
	}
	this.EventHandlers = append(this.EventHandlers, event)
	return event
}

func (this *Controller) String() string {
	var buffer bytes.Buffer
	var injectBuffer bytes.Buffer
	for k, _ := range this.InjectHash {
		injectBuffer.WriteString(", ")
		injectBuffer.WriteString(k)
	}
	injectStr := injectBuffer.String()
	buffer.WriteString(fmt.Sprintf(`function Ctrl_%d($scope, $element%s) {`, this.Id, injectStr))
	for _, event := range this.EventHandlers {
		if len(event.Type) > 0 {
			buffer.WriteString("\n")
			buffer.WriteString(fmt.Sprintf("$scope.EventHandler_%d_%d = function%s {", this.Id, event.Id, event.Args))
		}
		buffer.WriteString(event.Body)
		if len(event.Type) > 0 {
			buffer.WriteString("}")
			buffer.WriteString("\n")
		}
	}
	buffer.WriteString(`}`)
	return buffer.String()
}

type CssContext struct {
	Id       int64
	selector string
	bodyList []string
	child    []*CssContext
	parent   *CssContext
}

func (this *CssContext) Selector() string {
	var buffer bytes.Buffer
	idx := strings.Index(this.selector, "&")
	if idx >= 0 {
		buffer.WriteString(this.selector[idx+1:])
	} else {
		buffer.WriteString(" ")
		buffer.WriteString(this.selector)
	}
	return buffer.String()
}

func (this *CssContext) String() string {
	var buffer bytes.Buffer
	selectorList := make([]string, 0)
	selectorList = append(selectorList, this.Selector())
	current := this.parent
	for current != nil {
		selectorList = append(selectorList, current.Selector())
		current = current.parent
	}
	for i := len(selectorList) - 1; i >= 0; i-- {
		buffer.WriteString(selectorList[i])
	}
	buffer.WriteString(" {\n")
	buffer.WriteString(strings.Join(this.bodyList, "\n"))
	buffer.WriteString("\n}\n")
	for _, c := range this.child {
		buffer.WriteString(c.String())
	}
	return buffer.String()
}

func (this *CssContext) visit(visitor func(string, *CssContext)) {
	var buffer bytes.Buffer
	selectorList := make([]string, 0)
	selectorList = append(selectorList, this.Selector())
	current := this.parent
	for current != nil {
		selectorList = append(selectorList, current.Selector())
		current = current.parent
	}
	for i := len(selectorList) - 1; i >= 0; i-- {
		buffer.WriteString(selectorList[i])
	}
	selector := buffer.String()

	visitor(selector, this)

	for _, c := range this.child {
		c.visit(visitor)
	}
}

func NewCssContext(Id int64) *CssContext {
	return &CssContext{
		Id:       Id,
		bodyList: make([]string, 0),
		child:    make([]*CssContext, 0),
	}
}

type Css struct {
	Id        int64
	root      *CssContext
	inherit   bool
	contextId int64
}

func NewCss(Id int64) *Css {
	css := &Css{
		Id: Id,
	}
	css.root = NewCssContext(css.contextId)
	css.contextId++
	return css
}

func (this *Css) Parse(src string) {
	s := scanner.New(src)
	current := this.root
	var lineBuffer bytes.Buffer
	for {
		token := s.Next()
		if token.Type == scanner.TokenError {
			panic(token.String())
		}
		if token.Type == scanner.TokenEOF {
			break
		}
		switch token.Type {
		case scanner.TokenIdent:
		case scanner.TokenAtKeyword:
		case scanner.TokenString:
		case scanner.TokenHash:
		case scanner.TokenNumber:
		case scanner.TokenPercentage:
		case scanner.TokenDimension:
		case scanner.TokenURI:
		case scanner.TokenUnicodeRange:
		case scanner.TokenCDO:
		case scanner.TokenCDC:
		case scanner.TokenS:
			if strings.Contains(token.Value, "\n") {
				current.bodyList = append(current.bodyList, strings.TrimSpace(lineBuffer.String()))
				lineBuffer.Reset()
				continue
			}
		case scanner.TokenComment:
		case scanner.TokenFunction:
		case scanner.TokenIncludes:
		case scanner.TokenDashMatch:
		case scanner.TokenPrefixMatch:
		case scanner.TokenSuffixMatch:
		case scanner.TokenSubstringMatch:
		case scanner.TokenChar:
			switch token.Value {
			case "{":
				newOne := NewCssContext(this.contextId)
				this.contextId++
				newOne.parent = current
				newOne.selector = strings.TrimSpace(lineBuffer.String())
				lineBuffer.Reset()
				current.child = append(current.child, newOne)
				current = newOne
				continue
			case "}":
				current = current.parent
				continue
			}
		case scanner.TokenBOM:
		}
		lineBuffer.WriteString(token.Value)
	}
}

func (this *Css) String() string {
	var buffer bytes.Buffer
	for _, c := range this.root.child {
		buffer.WriteString(c.String())
	}
	return buffer.String()
}

func (this *Css) Visit(visitor func(string, *CssContext)) {
	for _, c := range this.root.child {
		c.visit(visitor)
	}
}

func (this *Css) Apply(jPlate *goquery.Selection) {
	this.Visit(func(selector string, ctx *CssContext) {
		defer func() {
			if err := recover(); err != nil {
				//fmt.Println(err)
			}
		}()
		class := fmt.Sprintf("genclass_%d_%d", this.Id, ctx.Id)
		jPlate.Find(selector).Each(func(idx int, jObj *goquery.Selection) {
			if !jObj.HasClass(class) {
				jObj.AddClass(class)
				ctx.selector = fmt.Sprintf("%s.%s", ctx.selector, class)
			}
		})
	})
}
