package goca

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/kolo/xmlrpc"
	"launchpad.net/xmlpath"
)

var (
	client *oneClient
)

type oneClient struct {
	token        string
	xmlrpcClient *xmlrpc.Client
}

type response struct {
	status bool
	body   string
}

type xmlResource struct {
	body string
}

type xmlIter struct {
	iter *xmlpath.Iter
}

type xmlNode struct {
	node *xmlpath.Node
}

func init() {
	err := SetClient()
	if err != nil {
		log.Fatal(err)
	}
}

func Client() *oneClient {
	return client
}

func SetClient(args ...string) error {
	var auth_token string
	var one_auth_path string

	if len(args) == 1 {
		auth_token = args[0]
	} else {
		one_auth_path = os.Getenv("ONE_AUTH")
		if one_auth_path == "" {
			one_auth_path = os.Getenv("HOME") + "/.one/one_auth"
		}

		token, err := ioutil.ReadFile(one_auth_path)
		if err == nil {
			auth_token = strings.TrimSpace(string(token))
		} else {
			return err
		}
	}

	one_xmlrpc := os.Getenv("ONE_XMLRPC")
	if one_xmlrpc == "" {
		one_xmlrpc = "http://localhost:2633/RPC2"
	}

	xmlrpcClient, err := xmlrpc.NewClient(one_xmlrpc, nil)
	if err != nil {
		log.Fatal(err)
	}

	client = &oneClient{
		token:        auth_token,
		xmlrpcClient: xmlrpcClient,
	}

	return nil
}

func SystemVersion() string {
	response, err := client.Call("one.system.version")
	if err != nil {
		log.Fatal(err)
	}

	return response.String()
}

func (c *oneClient) Call(method string, args ...interface{}) (*response, error) {
	result := []interface{}{}

	xmlArgs := make([]interface{}, len(args)+1)

	xmlArgs[0] = c.token
	copy(xmlArgs[1:], args[:])

	err := c.xmlrpcClient.Call(method, xmlArgs, &result)
	if err != nil {
		log.Fatal(err)
	}

	var ok bool

	status, ok := result[0].(bool)
	if ok == false {
		log.Fatal("Unexpected XML-RPC response. Expected: Index 0 Boolean ")
	}

	body, ok := result[1].(string)
	if ok == false {
		log.Fatal("Unexpected XML-RPC response. Expected: Index 0 String ")
	}

	// TODO: errCode? result[2]

	r := &response{status, body}

	if status == false {
		err = errors.New(body)
	}

	return r, err
}

func (r *response) String() string {
	return r.body
}

func (r *xmlResource) Body() string {
	return r.body
}

func (r *xmlResource) XPath(xpath string) (string, bool) {
	path := xmlpath.MustCompile(xpath)
	b := bytes.NewBufferString(r.Body())

	root, _ := xmlpath.Parse(b)

	return path.String(root)
}

func (r *xmlResource) XPathIter(xpath string) *xmlIter {
	path := xmlpath.MustCompile(xpath)
	b := bytes.NewBufferString(string(r.Body()))

	root, _ := xmlpath.Parse(b)

	return &xmlIter{iter: path.Iter(root)}
}

func (r *xmlResource) GetIdFromName(name string, xpath string) (uint, error) {
	var id int
	var match bool = false

	iter := r.XPathIter(xpath)
	for iter.Next() {
		node := iter.Node()

		n, _ := node.XPathNode("NAME")
		if n == name {
			if match {
				return 0, errors.New("Multiple resources with that name.")
			}

			idString, _ := node.XPathNode("ID")
			id, _ = strconv.Atoi(idString)
			match = true
		}
	}

	if match {
		return uint(id), nil
	} else {
		return 0, errors.New("Resource not found.")
	}
}

func (i *xmlIter) Next() bool {
	return i.iter.Next()
}

func (i *xmlIter) Node() *xmlNode {
	return &xmlNode{node: i.iter.Node()}
}

func (n *xmlNode) XPathNode(xpath string) (string, bool) {
	path := xmlpath.MustCompile(xpath)
	return path.String(n.node)
}
