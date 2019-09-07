package renderer

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/hairyhenderson/gomplate"
	"github.com/imdario/mergo"
)

type (
	Renderer interface {
		Render() (*bytes.Buffer, error)
	}

	Options struct {
		TemplateReaders map[string]io.Reader
		ValueReaders    map[string][]io.Reader
		RootNamespace   string
		LeftDelim       string
		RightDelim      string
		Name            string
	}

	Values map[string]interface{}

	renderer struct {
		templateReaders map[string]io.Reader
		valueReaders    map[string][]io.Reader
		leftDelim       string
		rightDelim      string
		name            string
	}
)

func New(opt *Options) Renderer {
	return &renderer{
		templateReaders: opt.TemplateReaders,
		valueReaders:    opt.ValueReaders,
		leftDelim:       opt.LeftDelim,
		rightDelim:      opt.RightDelim,
		name:            opt.Name,
	}
}

func (r *renderer) Render() (*bytes.Buffer, error) {
	var content []string
	out := new(bytes.Buffer)
	root := Values{}

	for _, reader := range r.templateReaders {
		res, err := ioutil.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		content = append(content, string(res))
	}

	for name, readers := range r.valueReaders {
		vals := Values{}
		for _, reader := range readers {
			v := Values{}
			res, err := ioutil.ReadAll(reader)
			if err != nil {
				return nil, err
			}
			err = yaml.Unmarshal([]byte(string(res)), &v)
			if err != nil {
				return nil, err
			}
			mergo.Merge(&vals, v)
		}
		root[name] = vals

	}

	tmpl, err := template.New(r.name).Funcs(gomplate.Funcs(nil)).Delims(r.leftDelim, r.rightDelim).Parse(strings.Join(content, "\n"))
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(out, root)
	if err != nil {
		return nil, err
	}
	return out, nil
}
