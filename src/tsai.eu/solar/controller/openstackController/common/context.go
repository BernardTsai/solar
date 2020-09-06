package common

import (
	"errors"
	"net/http"
	"net/url"
	"crypto/tls"
)

//------------------------------------------------------------------------------

// Context is a runtime information base for process
type Context struct {
  Request                 *Request
  Response                *Response
  TenantConfiguration     *TenantConfiguration
  TenantEndpoint          *TenantEndpoint
  TenantRelationshipState *RelationshipState
	Channel                 chan *Context
}

//------------------------------------------------------------------------------

// CreateContext creates a runtime information base from a request and response
func CreateContext(request *Request, response *Response) *Context {
  return &Context{
    Request:                 request,
    Response:                response,
    TenantConfiguration:     nil,
    TenantEndpoint:          nil,
    TenantRelationshipState: nil,
		Channel:                 make(chan *Context),
  }
}

//------------------------------------------------------------------------------

// GetRelationshipState determines a specific relationship
func (context *Context) GetRelationshipState(relationship string) (*RelationshipState){
  for _,relationshipState := range context.Request.Relationships {
    if relationshipState.Relationship == relationship {
      return &relationshipState
    }
  }

  return nil
}

//------------------------------------------------------------------------------

// SetStatus set the status and code of a response
func (context *Context) SetStatus(code int, status string) {
  context.Response.Code   = code
  context.Response.Status = status
}

//------------------------------------------------------------------------------

// DoGet conducts a get request to the OpenStack API
func (context *Context) DoGet(service string, suffix string) (response *http.Response, err error) {
	// determine tenant relationship
	context.TenantRelationshipState = context.GetRelationshipState("tenant")
	if context.TenantRelationshipState == nil {
		return nil, errors.New("Undefined tenant relationship")
	}

	// determine tenant endpoint, token and proxy
	endpoint, token, proxy := GetEndpointTokenAndProxy(context.TenantRelationshipState.Endpoint, service)
	if endpoint == "" {
		return nil, errors.New("Unable to parse tenant endpoint")
	}

	// define transport parameters
	var transport *http.Transport

	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return nil, errors.New("Invalid proxy setting")
		}

		transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(proxyURL),

		}
	} else {
		transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true},}
	}

	// issue a get request
  client    := &http.Client{Transport: transport}
  req, _    := http.NewRequest("GET", endpoint + suffix, nil)
  req.Header.Set("X-Auth-Token", token)
  return client.Do(req)
}

//------------------------------------------------------------------------------
