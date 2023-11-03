package api_test

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Reisender/go-api"
	"github.com/Reisender/go-api/middleware"
)

func TestNewClient(t *testing.T) {
	wantBody := `{                                                                                                                                                                                                      
  "links": [                                                                                                                                                                                                          
    {                                                                                                                                                                                                                 
      "rel": "self",                                                                                                                                                                                                  
      "uri": "/v3.0/users?limit=1"                                                                                                                                                                                    
    },                                                                                                                                                                                                                
    {                                                                                                                                                                                                                 
      "rel": "next",                                                                                                                                                                                                  
      "uri": "/v3.0/users?limit=1&starting_after=58da8c63d7dc0ca0680003ed"                                                                                                                                            
    }                                                                                                                                                                                                                 
  ],                                                                                                                                                                                                                  
  "data": [                                                                                                                                                                                                           
    {                                                                                                                                                                                                                 
      "uri": "/v3.0/users/58da8c63d7dc0ca0680003ed"                                                                                                                                                                   
    }                                                                                                                                                                                                                 
  ]                                                                                                                                                                                                                   
}`

	// make a mock Do func that checks the req URL and returns the want body
	mockDo := middleware.NewMock(func(req *http.Request) (*http.Response, error) {

		want := "localhost/v1/foo"
		if req.URL.String() != want {
			t.Errorf("request URL: want '%s' got '%s'", want, req.URL.String())
		}

		buf := bytes.NewBufferString(fmt.Sprintf("HTTP/1.1 200\n\n%s", wantBody))
		return http.ReadResponse(bufio.NewReader(buf), req)
	})

	ctx := context.Background()

	// create the client with a mock Do func defined above
	c := api.NewClient("localhost", "/v1", 0, mockDo)

	// build our URL from the client
	// this should add in the "base" of the client
	url, err := c.NewURL("/foo")
	if err != nil {
		t.Error(err)
		return
	}

	resp, err := c.Get(ctx, url.String())
	if err != nil {
		t.Error(err)
		return
	}

	if resp == nil {
		t.Error("response is nil")
		return
	}

	gotBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Error(err)
		return
	}

	if wantBody != string(gotBody) {
		t.Errorf("\nwant '%s'\ngot '%s'\n", wantBody, gotBody)
	}
}
