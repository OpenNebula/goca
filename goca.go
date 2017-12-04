package goca

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/kolo/xmlrpc"
	"gopkg.in/xmlpath.v2"
)

var (
	client *oneClient
)

const (
	// PoolWhoMine to list resources that belong to the user that performs the
	// query.
	PoolWhoMine = -3

	// PoolWhoAll to list all the resources seen by the user that performs the
	// query.
	PoolWhoAll = -2

	// PoolWhoGroup to list all the resources that belong to the group that performs
	// the query.
	PoolWhoGroup = -1
)

// OneConfig contains the information to communicate with OpenNebula
type OneConfig struct {
	// Token is the authentication string. In the format of <user>:<password>
	Token string

	// XmlrpcURL contains OpenNebula's XML-RPC API endpoint. Defaults to
	// http://localhost:2633/RPC2
	XmlrpcURL string
}

type oneClient struct {
	token             string
	xmlrpcClient      *xmlrpc.Client
	xmlrpcClientError error
}

type response struct {
	status  bool
	body    string
	bodyInt int
}

// Resource implements an OpenNebula Resource methods. *XMLResource implements
// all these methods
type Resource interface {
	Body() string
	XPath(string) (string, bool)
	XPathIter(string) *XMLIter
	GetIDFromName(string, string) (uint, error)
}

// XMLResource contains an XML body field. All the resources in OpenNebula are
// of this kind.
type XMLResource struct {
	body string
}

// XMLIter is used to iterate over XML xpaths in an object.
type XMLIter struct {
	iter *xmlpath.Iter
}

// XMLNode represent an XML node.
type XMLNode struct {
	node *xmlpath.Node
}

// Initializes the client variable, used as a singleton
func init() {
	err := SetClient(NewConfig("", "", ""))
	if err != nil {
		log.Fatal(err)
	}
}

// NewConfig returns a new OneConfig object with the specified user, password,
// and xmlrpcURL
func NewConfig(user string, password string, xmlrpcURL string) OneConfig {
	var authToken string
	var oneAuthPath string

	oneXmlrpc := xmlrpcURL

	if user == "" && password == "" {
		oneAuthPath = os.Getenv("ONE_AUTH")
		if oneAuthPath == "" {
			oneAuthPath = os.Getenv("HOME") + "/.one/one_auth"
		}

		token, err := ioutil.ReadFile(oneAuthPath)
		if err == nil {
			authToken = strings.TrimSpace(string(token))
		} else {
			authToken = ""
		}
	} else {
		authToken = user + ":" + password
	}

	if oneXmlrpc == "" {
		oneXmlrpc = os.Getenv("ONE_XMLRPC")
		if oneXmlrpc == "" {
			oneXmlrpc = "http://localhost:2633/RPC2"
		}
	}

	config := OneConfig{
		Token:     authToken,
		XmlrpcURL: oneXmlrpc,
	}

	return config
}

// SetClient assigns a value to the client variable
func SetClient(conf OneConfig) error {

	xmlrpcClient, xmlrpcClientError := xmlrpc.NewClient(conf.XmlrpcURL, nil)

	client = &oneClient{
		token:             conf.Token,
		xmlrpcClient:      xmlrpcClient,
		xmlrpcClientError: xmlrpcClientError,
	}

	return nil
}

// SystemVersion returns the current OpenNebula Version
func SystemVersion() (string, error) {
	response, err := client.Call("one.system.version")
	if err != nil {
		return "", err
	}

	return response.Body(), nil
}

// Call is an XML-RPC wrapper. It returns a pointer to response and an error.
func (c *oneClient) Call(method string, args ...interface{}) (*response, error) {
	var (
		ok bool

		status  bool
		body    string
		bodyInt int64
	)

	if c.xmlrpcClientError != nil {
		return nil, fmt.Errorf("Unitialized client. Token: '%s', xmlrpcClient: '%s'", c.token, c.xmlrpcClientError)
	}

	result := []interface{}{}

	xmlArgs := make([]interface{}, len(args)+1)

	xmlArgs[0] = c.token
	copy(xmlArgs[1:], args[:])

	err := c.xmlrpcClient.Call(method, xmlArgs, &result)
	if err != nil {
		log.Fatal(err)
	}

	status, ok = result[0].(bool)
	if ok == false {
		log.Fatal("Unexpected XML-RPC response. Expected: Index 0 Boolean")
	}

	body, ok = result[1].(string)
	if ok == false {
		bodyInt, ok = result[1].(int64)
		if ok == false {
			log.Fatal("Unexpected XML-RPC response. Expected: Index 0 Int or String")
		}
	}

	// TODO: errCode? result[2]

	r := &response{status, body, int(bodyInt)}

	if status == false {
		err = errors.New(body)
	}

	return r, err
}

// Body accesses the body of the response
func (r *response) Body() string {
	return r.body
}

// BodyInt accesses the body of the response, if it's an int.
func (r *response) BodyInt() int {
	return r.bodyInt
}

// Body accesses the body of an XMLResource
func (r *XMLResource) Body() string {
	return r.body
}

// XPath returns the string pointed at by xpath, for an XMLResource
func (r *XMLResource) XPath(xpath string) (string, bool) {
	path := xmlpath.MustCompile(xpath)
	b := bytes.NewBufferString(r.Body())

	root, _ := xmlpath.Parse(b)

	return path.String(root)
}

// XPathIter returns an XMLIter object pointed at by the xpath
func (r *XMLResource) XPathIter(xpath string) *XMLIter {
	path := xmlpath.MustCompile(xpath)
	b := bytes.NewBufferString(string(r.Body()))

	root, _ := xmlpath.Parse(b)

	return &XMLIter{iter: path.Iter(root)}
}

// GetIDFromName finds the a resource by ID by looking at an xpath contained
// in that resource
func (r *XMLResource) GetIDFromName(name string, xpath string) (uint, error) {
	var id int
	var match = false

	iter := r.XPathIter(xpath)
	for iter.Next() {
		node := iter.Node()

		n, _ := node.XPathNode("NAME")
		if n == name {
			if match {
				return 0, errors.New("multiple resources with that name")
			}

			idString, _ := node.XPathNode("ID")
			id, _ = strconv.Atoi(idString)
			match = true
		}
	}

	if !match {
		return 0, errors.New("resource not found")
	}

	return uint(id), nil
}

// Next moves on to the next resource
func (i *XMLIter) Next() bool {
	return i.iter.Next()
}

// Node returns the XMLNode
func (i *XMLIter) Node() *XMLNode {
	return &XMLNode{node: i.iter.Node()}
}

// XPathNode returns an XMLNode pointed at by xpath
func (n *XMLNode) XPathNode(xpath string) (string, bool) {
	path := xmlpath.MustCompile(xpath)
	return path.String(n.node)
}
