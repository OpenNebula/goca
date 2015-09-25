package goca

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/kolo/xmlrpc"
	"launchpad.net/xmlpath"
)

var (
	client *Client
)

type Client struct {
	token        string
	xmlrpcClient *xmlrpc.Client
}

type Response struct {
	Status bool
	Body   string
}

type XML string

func init() {
	err := SetClient()
	if err != nil {
		log.Fatal(err)
	}
}

func SystemVersion() string {
	response, err := client.Call("one.system.version")
	if err != nil {
		log.Fatal(err)
	}

	return response.Body
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
		return err
	}

	client = &Client{
		token:        auth_token,
		xmlrpcClient: xmlrpcClient,
	}

	return nil
}

func (r *Response) String() string {
	return fmt.Sprintf("Status: %v\nBody:\n%v\n", r.Status, r.Body)
}

func (r *Response) Debug() error {
	fmt.Println(r)
	return nil
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
