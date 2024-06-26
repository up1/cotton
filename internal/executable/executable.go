package executable

import (
	"cotton/internal/capture"
	"cotton/internal/execution"
	"cotton/internal/httphelper"
	"cotton/internal/logger"
	"cotton/internal/template"
	"cotton/internal/variable"
	"errors"
	"slices"
)

type ExecutableOptions struct {
	Logger          logger.Logger
	RequestParser   httphelper.RequestParser
	InsecureRequest bool
	// ClockWrapper    clock.ClockWrapper
}

// For setups and teardowns
type Executable struct {
	options  *ExecutableOptions
	title    string
	reqRaw   string
	captures []*capture.Capture
	// stopwatch    *stopwatch.Stopwatch
	// ellapsedTime *stopwatch.EllapsedTime
}

func New(title, reqRaw string, options *ExecutableOptions) *Executable {
	return &Executable{
		options:  options,
		title:    title,
		reqRaw:   reqRaw,
		captures: []*capture.Capture{},
		// ellapsedTime: nil,
		// stopwatch:    stopwatch.New(options.ClockWrapper),
	}
}

func (ex *Executable) SetTitle(title string) {
	ex.title = title
}

func (ex *Executable) Title() string {
	if ex.title == "" {
		return "Untitled"
	}
	return ex.title
}

func (ex *Executable) RawRequest() string {
	return ex.reqRaw
}

func (ex *Executable) Captures() []*capture.Capture {
	// return clone of captures
	clones := []*capture.Capture{}
	for _, cap := range ex.captures {
		clones = append(clones, cap.Clone())
	}
	return clones
}

func (ex *Executable) AddCapture(capture *capture.Capture) {
	ex.captures = append(ex.captures, capture.Clone())
}

func (ex *Executable) Clone() *Executable {
	capturesClone := []*capture.Capture{}
	for _, cap := range ex.captures {
		capturesClone = append(capturesClone, cap.Clone())
	}

	return &Executable{
		options:  ex.options,
		title:    ex.title,
		reqRaw:   ex.reqRaw,
		captures: capturesClone,
	}
}

func (ex *Executable) Execute(initialVars *variable.Variables) (*execution.Execution, error) {
	// ex.stopwatch.Start()
	// defer func() { ex.ellapsedTime = ex.stopwatch.Stop() }()

	if ex.reqRaw == "" {
		return nil, errors.New("no callable request")
	}

	reqTemplate := template.New(ex.reqRaw)
	compiledRequest := reqTemplate.Apply(initialVars)

	request, err := ex.options.RequestParser.Parse(compiledRequest)
	if err != nil {
		return nil, err
	}

	ex.options.Logger.PrintExecutableTitle(ex.Title())
	ex.options.Logger.PrintRequest(compiledRequest)
	resp, err := request.Do(ex.options.InsecureRequest)
	if err != nil {
		return nil, err
	}

	vars := variable.New()
	for _, cap := range ex.captures {
		value, err := resp.ValueOf(cap.Selector)
		if err != nil {
			return nil, err
		}
		vars.Set(cap.Name, value)
	}

	return &execution.Execution{
		Variables: initialVars.MergeWith(vars),
	}, nil
}

func (ex *Executable) SimilarTo(anotherEx *Executable) bool {
	return ex.title == anotherEx.title &&
		ex.reqRaw == anotherEx.reqRaw &&
		slices.EqualFunc(ex.captures, anotherEx.Captures(), func(c1, c2 *capture.Capture) bool {
			return c1.SimilarTo(c2)
		})
}
