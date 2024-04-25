package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"elahi-arman.github.com/example-http-server/internal/datastore"
	"github.com/julienschmidt/httprouter"
)

type JSON struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// This represents a single hyperlink.
type HATEOAS struct {
	Rel    string `json:"rel"`
	Href   string `json:"href"`
	Method string `json:"method"`
	Header []JSON `json:"header"`
	Body   []JSON `json:"body"`
}

// This wraps a data payload with HATEOAS links.
type Resource struct {
	ResponseData interface{} `json:"response"`
	HyperLinks   []HATEOAS   `json:"links"`
}

// GetLink is the function called when a user makes a request to retrieve a certain link
func (s *serverImpl) GetLink(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// ps are the parameters attached to this route. the paramter to ByName()
	// must match the name of the link from main.go
	linkId := ps.ByName("link")

	// do some preemptive error checking
	if linkId == "" {
		fmt.Println("GetLink: no linkId provided")
		w.WriteHeader(400)
		return
	}

	// access the datastore attached to the server and try to fetch the link
	link, err := s.linkStore.GetLink(linkId)
	if errors.Is(err, &datastore.NotFoundError{}) {
		fmt.Printf("GetLink: no entry for linkId=%s\n", linkId)
		w.WriteHeader(404)
		return
	}

	// return a 302 to redirect users
	fmt.Printf("GetLink: found link for linkId=%s, redirecting to url=%s", link.Id, link.Url)

	// w.Header().Add("Location", link.Url) // the location header is the destination URL
	// w.WriteHeader(302)                   // 302 informs the client to read the Location header for a redirection

	resource := Resource{
		ResponseData: link,
		HyperLinks: []HATEOAS{
			{"self", "/l/" + linkId, "GET", []JSON{}, []JSON{}},
			{"getUserLinks", "/userlinks", "GET", []JSON{{Key: "user", Value: link.Owner}}, []JSON{}},
			{"createLinks", "/api/links", "POST", []JSON{{Key: "content-type", Value: "application/json"}}, []JSON{{Key: "url", Value: "{webpage-url}"}, {Key: "owner", Value: link.Owner}}},
			{"deleteUserLinks", "/api/delete/" + link.Owner + "/" + linkId, "DELETE", []JSON{}, []JSON{}},
		},
	}

	jsonbody, err := json.Marshal(resource)
	if err != nil {
		// If encoding fails send an HTTP 500 Internal Server Error.
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonbody)

}

// createLinkParams represents the structure of the request body to
// a CreateLink function call
type createLinkParams struct {
	Url string `json:"url"`
	// temporary, eventually we'll replace this by retrieving from context
	Owner string `json:"owner"`
}

func (s *serverImpl) CreateLink(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// retrieve the value of the content-type header, if none is specified
	// the request should be rejected
	contentType := r.Header.Get("content-type")
	if contentType == "" {
		fmt.Println("CreateLink: no content-type header is sent")
		w.WriteHeader(400) // the status message will automatically be filled in
		return
	}

	var url string
	var owner string
	if strings.Contains(contentType, "json") {
		// read the body of the request
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("CreateLink: error while reading body of request %v\n", err)
			w.WriteHeader(400)
			return
		}

		// convert the request body into json
		lp := &createLinkParams{}
		err = json.Unmarshal(body, lp)
		if err != nil {
			fmt.Printf("CreateLink: error while unmarshalling err=%v. \n body=%s\n", err, body)
			w.WriteHeader(400)
			return
		}

		url = lp.Url
		owner = lp.Owner
	} else if strings.Contains(contentType, "form") {

		// when dealing with form data, call ParseForm to trigger parsing
		// then r.Form will have a map of the form values
		r.ParseForm()
		if formUrl, ok := r.Form["url"]; !ok || len(formUrl) == 0 || formUrl[0] == "" {
			fmt.Println("CreateLink: url key is not part of form data")
			w.Header().Add("Location", fmt.Sprintf("/public?error=%s", "cannot create a link without a url"))
			w.WriteHeader(303)
			return
		} else {
			url = formUrl[0]
		}

		if formOwner, ok := r.Form["owner"]; !ok || len(formOwner) == 0 || formOwner[0] == "" {
			fmt.Println("CreateLink: owner key is not part of form data")
			w.Header().Add("Location", fmt.Sprintf("/public?error=%s", "cannot create a link without an owner"))
			w.WriteHeader(303)
			return
		} else {
			owner = formOwner[0]
		}
	}

	// call the datastore function
	link, err := s.linkStore.CreateLink(url, owner)
	if err != nil {
		fmt.Printf("CreateLink: error while creating a link err=%v\n", err)
		w.WriteHeader(500)
		return
	}

	// redirect users
	// w.Header().Add("Location", fmt.Sprintf("/public?link=%s", link.Id))
	// w.WriteHeader(303)

	resource := Resource{
		ResponseData: link,
		HyperLinks: []HATEOAS{
			{"self", "/api/links", "POST", []JSON{{Key: "content-type", Value: "application/json"}}, []JSON{{Key: "url", Value: "{webpage-url}"}, {Key: "owner", Value: link.Owner}}},
			{"getLink", "/l/" + link.Id, "GET", []JSON{}, []JSON{}},
			{"getUserLinks", "/userlinks", "GET", []JSON{{Key: "user", Value: link.Owner}}, []JSON{}},
			{"deleteUserLinks", "/api/delete/" + link.Owner + "/" + link.Id, "DELETE", []JSON{}, []JSON{}},
		},
	}
	jsonbody, err := json.Marshal(resource)
	if err != nil {
		// If encoding fails send an HTTP 500 Internal Server Error.
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonbody)

}

// getUserLinksParams represents the structure of the request body to
// a getUserLinks function call
type getUserLinksParams struct {
	URLs  []string `json:"urls"`
	Owner string   `json:"owner"`
}

// read a header / body to get a user
// return a list of links in json format where Owner == user passed in
func (s *serverImpl) GetUserLinks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	user := r.Header.Get("User")

	if user == "" {
		fmt.Println("GetUserLinks: no user provided")
		w.WriteHeader(400)
		return
	}

	links, err := s.linkStore.GetUserLinks(user)

	if errors.Is(err, &datastore.NotFoundError{}) {
		fmt.Printf("GetUserLinks: No user found by name =%s\n", user)
		w.WriteHeader(204)
		return
	}

	// Extract URLs from each link in the array
	var urls []string
	for _, link := range links {
		urls = append(urls, link.Url)
	}
	fmt.Printf("GetUserLinks: found links for user=%s is %v ", user, urls)
	//  Json encode the body
	jsonbody, err := json.Marshal(getUserLinksParams{
		URLs:  urls,
		Owner: user,
	})
	if err != nil {
		// If encoding fails send an HTTP 500 Internal Server Error.
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonbody))

}
func (s *serverImpl) DeleteLink(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	linkId := ps.ByName("link")
	owner := ps.ByName("user")

	if linkId == "" || owner == "" {
		fmt.Println("DeleteLink: no link or user provided")
		w.WriteHeader(400)
		return
	}

	err := s.linkStore.DeleteLink(linkId, owner)

	if errors.Is(err, &datastore.NotFoundError{}) {
		fmt.Printf("DeleteLink: No Links were deleted")
		w.Write([]byte("No Links were deleted"))
		return
	}

	w.Write([]byte("Delete Link Completed"))
}
