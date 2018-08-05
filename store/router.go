package store

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var controller = &Controller{Repository: Repository{}}

type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Authentication",
		"POST",
		"/get_token",
		controller.GetToken,
	},
	Route{
		"ConvertImage",
		"POST",
		"/convertimage",
		controller.ConvertImageWithFileBuffer,
	},
	Route{
		"AddUser",
		"POST",
		"/adduser",
		controller.AddUser,
	},
	Route{
		"ConvertImages",
		"POST",
		"/convertimages",
		controller.ConvertImageFromApi,
	},
	/*Route{
		"ConvertImages",
		"POST",
		"/convertimages",
		controller.ConvertImages,
	},*/
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		log.Println(route.Name)
		handler = route.HandleFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
