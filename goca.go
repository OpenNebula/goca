package goca

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/kolo/xmlrpc"
	"launchpad.net/xmlpath"
)

var (
	client *Client
)

func init() {
	var err error
	client, err = NewClient()
	if err != nil {
		log.Fatal(err)
	}
}

type Client struct {
	token        string
	xmlrpcClient *xmlrpc.Client
}

type Response struct {
	Status bool
	Body   string
}

type XML string

func (r *Response) String() string {
	return fmt.Sprintf("Status: %v\nBody:\n%v\n", r.Status, r.Body)
}

func (r *Response) Debug() error {
	fmt.Println(r)
	return nil
}

func NewClient() (*Client, error) {
	token, err := ioutil.ReadFile("/var/lib/one/.one/one_auth")
	if err != nil {
		return nil, err
	}

	xmlrpcClient, err := xmlrpc.NewClient("http://localhost:2633/RPC2", nil)
	if err != nil {
		return nil, err
	}

	c := &Client{
		token:        strings.TrimSpace(string(token)),
		xmlrpcClient: xmlrpcClient,
	}

	return c, nil
}

func (c *Client) SystemVersion() (r *Response, err error) {
	r, err = c.Call("one.system.version")
	return
}

func (c *Client) Call(method string, args ...interface{}) (r *Response, err error) {
	result := []interface{}{}

	xmlArgs := make([]interface{}, len(args)+1)

	xmlArgs[0] = c.token
	copy(xmlArgs[1:], args[:])

	err = c.xmlrpcClient.Call(method, xmlArgs, &result)
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

	r = &Response{Status: status, Body: body}

	if status == false {
		err = errors.New(body)
	}

	return
}

func (xml XML) XPath(xpath string) (string, bool) {
	path := xmlpath.MustCompile(xpath)
	b := bytes.NewBufferString(string(xml))

	root, _ := xmlpath.Parse(b)

	return path.String(root)
}

func (xml XML) XPathIter(xpath string) *xmlpath.Iter {
	path := xmlpath.MustCompile(xpath)
	b := bytes.NewBufferString(string(xml))

	root, _ := xmlpath.Parse(b)

	return path.Iter(root)
}
