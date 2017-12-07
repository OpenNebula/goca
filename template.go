package goca

import (
	"errors"
)

// TemplateRoot is the XML root node
const TemplateRoot = "VMTEMPLATE"

// Template represents an OpenNebula Template
type Template struct {
	XMLResource
	ID   uint
	Name string
}

// TemplatePool represents an OpenNebula TemplatePool
type TemplatePool struct {
	XMLResource
}

// CreateTemplate allocates a new template. It returns the new template ID.
func CreateTemplate(template string) (uint, error) {
	response, err := client.Call("one.template.allocate", template)
	if err != nil {
		return 0, err
	}

	return uint(response.BodyInt()), nil
}

// NewTemplatePool returns a template pool. A connection to OpenNebula is
// performed.
func NewTemplatePool(args ...int) (*TemplatePool, error) {
	var who, start, end int

	switch len(args) {
	case 0:
		who = PoolWhoMine
		start = -1
		end = -1
	case 3:
		who = args[0]
		start = args[1]
		end = args[2]
	default:
		return nil, errors.New("Wrong number of arguments")
	}

	response, err := client.Call("one.templatepool.info", who, start, end)
	if err != nil {
		return nil, err
	}

	templatepool := &TemplatePool{XMLResource{body: response.Body()}}

	return templatepool, err

}

// NewTemplate finds a template object by ID. No connection to OpenNebula.
func NewTemplate(id uint) *Template {
	return &Template{ID: id}
}

// NewTemplateFromName finds a template object by name. It connects to
// OpenNebula to retrieve the pool, but doesn't perform the Info() call to
// retrieve the attributes of the template.
func NewTemplateFromName(name string) (*Template, error) {
	templatePool, err := NewTemplatePool()
	if err != nil {
		return nil, err
	}

	id, err := templatePool.GetIDFromName(name, "/VMTEMPLATE_POOL/VMTEMPLATE")
	if err != nil {
		return nil, err
	}

	return NewTemplate(id), nil
}

// Info connects to OpenNebula and fetches the information of the Template
func (template *Template) Info() error {
	response, err := client.Call("one.template.info", template.ID)
	template.body = response.Body()
	return err
}

// Delete will remove the template from OpenNebula.
func (template *Template) Delete() error {
	_, err := client.Call("one.template.delete", template.ID)
	return err
}

// Instantiate will instantiate the template
func (template *Template) Instantiate(name string, pending bool, extra string) (uint, error) {
	response, err := client.Call("one.template.instantiate", template.ID, name, pending, extra)

	if err != nil {
		return 0, err
	}

	return uint(response.BodyInt()), nil
}

// Update will modify the template. If appendTemplate is 0, it will
// replace the whole template. If its 1, it will merge.
func (template *Template) Update(tpl string, appendTemplate int) error {
	_, err := client.Call("one.template.update", template.ID, tpl, appendTemplate)
	return err
}
