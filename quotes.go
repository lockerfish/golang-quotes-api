package quotes

import (
	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Quote struct {
	UID     *datastore.Key `json:"uid" datastore:"-"`
	Author  string         `json:"author"`
	Message string         `json:"message"`
}

type QuotesAPI struct{}

type Quotes struct {
	Quotes []Quote `json:"quotes"`
}

type AddRequest struct {
	Author  string
	Message string
}

func (QuotesAPI) Add(c context.Context, r *AddRequest) (*Quote, error) {
	k := datastore.NewIncompleteKey(c, "Quote", quoteKey(c))

	t := &Quote{Author: r.Author, Message: r.Message}

	k, err := datastore.Put(c, k, t)
	if err != nil {
		return nil, err
	}
	t.UID = k
	return t, nil
}

func (QuotesAPI) List(c context.Context) (*Quotes, error) {
	quotes := []Quote{}
	// If we omitted the .Ancestor from this query there would be
	// a slight chance that Quote that had just been written would not
	// show up in a query.
	keys, err := datastore.NewQuery("Quote").Ancestor(quoteKey(c)).GetAll(c, &quotes)
	if err != nil {
		return nil, err
	}

	for i, k := range keys {
		quotes[i].UID = k
	}
	return &Quotes{quotes}, nil
}

func quoteKey(c context.Context) *datastore.Key {
	return datastore.NewKey(c, "Quote", "default_quote", 0, nil)
}

func init() {
	// register the quotes API with cloud endpoints.
	api, err := endpoints.RegisterService(QuotesAPI{}, "quotesService", "v1", "Quotes API", true)
	if err != nil {
		panic(err)
	}

	// adapt the name, method, and path for each method.
	info := api.MethodByName("List").Info()
	info.Name, info.HTTPMethod, info.Path = "getQuotes", "GET", "quotesService"

	info = api.MethodByName("Add").Info()
	info.Name, info.HTTPMethod, info.Path = "addQuote", "POST", "quotesService"

	// start handling cloud endpoint requests.
	endpoints.HandleHTTP()
}
