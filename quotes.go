package quotes

import (
	// "fmt"
	// "html/template"
	// "net/http"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"golang.org/x/net/context"
	// "google.golang.org/appengine"
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

func quoteKey(c context.Context) *datastore.Key {
	return datastore.NewKey(c, "Quote", "default_quote", 0, nil)
}

// func root(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprint(w, rootForm)
// }

// const rootForm = `
// 	<html>
// 		<body>
// 			<h1>Enter Author's quote</h1>
// 			<form action="/quote" method="post" accept-charset="utf-8">
// 				<input type="text" name="author" value="Author's Name: " id="author"/>
// 				<input type="text" name="quote" value="Write quote..." id="quote"/>
// 				<input type="submit" value="Submit Quote"/>
// 			</form>
// 		</body>
// 	</html>
// `

// var quoteTemplate = template.Must(template.New("quote").Parse(quoteTemplateHTML))

// func quote(w http.ResponseWriter, r *http.Request) {
// 	c := appengine.NewContext(r)
// 	q1 := Quote{
// 		Author:  r.FormValue("author"),
// 		Message: r.FormValue("quote"),
// 	}

// 	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Quote", quoteKey(c)), &q1)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	var q2 Quote
// 	err1 := datastore.Get(c, key, &q2)
// 	if err1 != nil {
// 		http.Error(w, err1.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	err2 := quoteTemplate.Execute(w, q2.Message)
// 	if err2 != nil {
// 		http.Error(w, err2.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

// const quoteTemplateHTML = `
// 	<html>
// 		<body>
// 			<p>You wrote:</p>
// 			<pre>{{html .}}</pre>
// 		</body>
// 	</html>
// `

// func fetch(w http.ResponseWriter, r *http.Request) {
// 	c := appengine.NewContext(r)
// 	q := datastore.NewQuery("Quote").Ancestor(quoteKey(c)).Limit(10)

// 	quotes := make([]Quote, 0, 10)
// 	_, err := q.GetAll(c, &quotes)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	err1 := fetchTemplate.Execute(w, quotes)
// 	if err1 != nil {
// 		http.Error(w, err1.Error(), http.StatusInternalServerError)
// 	}
// }

// var fetchTemplate = template.Must(template.New("fetch").Parse(`
// <html>
//   <head>
//     <title>All the Quotes</title>
//   </head>
//   <body>
//     <h2>All the Quotes</h2>
//     {{range .}}
//       {{with .Author}}
//         <p><b>{{.}}</b> wrote:</p>
//       {{else}}
//         <p>An anonymous person wrote:</p>
//       {{end}}
//       <pre>{{.Message}}</pre>
//     {{end}}
//   </body>
// </html>
// `))
