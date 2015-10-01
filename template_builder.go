package goca

import (
	"fmt"
	"strings"
)

type TemplateBuilder struct {
	elements []TemplateBuilderElement
}

type TemplateBuilderElement interface {
	String() string
}

type TemplateBuilderPair struct {
	key   string
	value string
}

type TemplateBuilderVector struct {
	key   string
	pairs []TemplateBuilderPair
}

func NewTemplateBuilder() *TemplateBuilder {
	return &TemplateBuilder{}
}

func (t *TemplateBuilder) NewVector(key string) *TemplateBuilderVector {
	vector := &TemplateBuilderVector{key: key}
	t.elements = append(t.elements, vector)
	return vector
}

func (t *TemplateBuilder) String() string {
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

func (t *TemplateBuilderPair) String() string {
	return fmt.Sprintf("%s=\"%s\"", t.key, t.value)
}

func (t *TemplateBuilderVector) String() string {
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

func (t *TemplateBuilder) AddValue(key, val string) error {
	pair := &TemplateBuilderPair{strings.ToUpper(key), val}
	t.elements = append(t.elements, pair)

	return nil
}

func (t *TemplateBuilderVector) AddValue(key, val string) error {
	pair := TemplateBuilderPair{strings.ToUpper(key), val}
	t.pairs = append(t.pairs, pair)

	return nil
}
