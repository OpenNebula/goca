package goca

import (
	"fmt"
	"strings"
)

type TemplateBuilder struct {
	body []byte
}

func NewTemplateBuilder() *TemplateBuilder {
	return &TemplateBuilder{}
}

func (t *TemplateBuilder) String() string {
	return string(t.body)
}

func (t *TemplateBuilder) AddValue(param, val string) error {
	paramUpper := strings.ToUpper(param)

	var startToken string
	if len(t.body) == 0 {
		startToken = ""
	} else {
		startToken = "\n"
	}

	s := fmt.Sprintf("%s%s=\"%s\"", startToken, paramUpper, val)
	t.body = append(t.body, s...)
	return nil
}

func (t *TemplateBuilder) AddVector(param, val string) error {
	paramUpper := strings.ToUpper(param)

	var startToken string
	if len(t.body) == 0 {
		startToken = ""
	} else {
		startToken = "\n"
	}

	valSplit := strings.Split(val, "\n")
	val = strings.Join(valSplit, ",\n")

	s := fmt.Sprintf("%s%s=[%s ]", startToken, paramUpper, val)
	t.body = append(t.body, s...)
	return nil
}

// var template, vector *goca.TemplateBuilder

// template = goca.NewTemplateBuilder()

// template.AddValue("cpu", "1")
// template.AddValue("memory", "64")

// vector = goca.NewTemplateBuilder()
// vector.AddValue("image_id", "119")
// vector.AddValue("image_id", "119")
// template.AddVector("disk", vector.String())

// vector = goca.NewTemplateBuilder()
// vector.AddValue("image_id", "119")
// vector.AddValue("image_id", "119")
// template.AddVector("disk", vector.String())

// fmt.Println(template)
