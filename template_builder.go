package goca

import (
	"fmt"
	"strings"
)

type Template struct {
	elements []TemplateElement
}

type TemplateElement interface {
	String() string
}

type TemplatePair struct {
	key   string
	value string
}

type TemplateVector struct {
	key   string
	pairs []TemplatePair
}

func NewTemplate() *Template {
	return &Template{}
}

func (t *Template) NewVector(key string) *TemplateVector {
	vector := &TemplateVector{key: key}
	t.elements = append(t.elements, vector)
	return vector
}

func (t *Template) String() string {
	s := ""
	endToken := "\n"

	for i, element := range t.elements {
		if i == len(t.elements)-1 {
			endToken = ""
		}
		s += element.String() + endToken
	}

	return s
}

func (t *TemplatePair) String() string {
	return fmt.Sprintf("%s=\"%s\"", t.key, t.value)
}

func (t *TemplateVector) String() string {
	s := fmt.Sprintf("%s=[\n", strings.ToUpper(t.key))

	endToken := ",\n"
	for i, pair := range t.pairs {
		if i == len(t.pairs)-1 {
			endToken = ""
		}

		s += fmt.Sprintf("    %s%s", pair.String(), endToken)

	}
	s += " ]"

	return s
}

func (t *Template) AddValue(key, val string) error {
	pair := &TemplatePair{strings.ToUpper(key), val}
	t.elements = append(t.elements, pair)

	return nil
}

func (t *TemplateVector) AddValue(key, val string) error {
	pair := TemplatePair{strings.ToUpper(key), val}
	t.pairs = append(t.pairs, pair)

	return nil
}
