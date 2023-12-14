/*
 * Сервис по обмену вещами Steam
 *
 * API for exchanging virtual items
 *
 * API version: 0.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello World!</h1>")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/api/v1/",
		Index,
	},

	Route{
		"ItemsGet",
		strings.ToUpper("Get"),
		"/api/v1/items",
		ItemsGet,
	},

	Route{
		"ItemsItemIdDelete",
		strings.ToUpper("Delete"),
		"/api/v1/items/{item_id}",
		ItemsItemIdDelete,
	},

	Route{
		"ItemsItemIdGet",
		strings.ToUpper("Get"),
		"/api/v1/items/{item_id}",
		ItemsItemIdGet,
	},

	Route{
		"ItemsPost",
		strings.ToUpper("Post"),
		"/api/v1/items",
		ItemsPost,
	},

	Route{
		"ItemsItemIdTradesGet",
		strings.ToUpper("Get"),
		"/api/v1/items/{item_id}/trades",
		ItemsItemIdTradesGet,
	},

	Route{
		"TradesGet",
		strings.ToUpper("Get"),
		"/api/v1/trades",
		TradesGet,
	},

	Route{
		"TradesPost",
		strings.ToUpper("Post"),
		"/api/v1/trades",
		TradesPost,
	},

	Route{
		"TradesTradeIdDelete",
		strings.ToUpper("Delete"),
		"/api/v1/trades/{trade_id}",
		TradesTradeIdDelete,
	},

	Route{
		"TradesTradeIdGet",
		strings.ToUpper("Get"),
		"/api/v1/trades/{trade_id}",
		TradesTradeIdGet,
	},

	Route{
		"TradesTradeIdPut",
		strings.ToUpper("Put"),
		"/api/v1/trades/{trade_id}",
		TradesTradeIdPut,
	},

	Route{
		"UsersUserIdTradesGet",
		strings.ToUpper("Get"),
		"/api/v1/users/{user_id}/trades",
		UsersUserIdTradesGet,
	},

	Route{
		"UsersGet",
		strings.ToUpper("Get"),
		"/api/v1/users",
		UsersGet,
	},

	Route{
		"UsersPost",
		strings.ToUpper("Post"),
		"/api/v1/users",
		UsersPost,
	},

	Route{
		"UsersUserIdDelete",
		strings.ToUpper("Delete"),
		"/api/v1/users/{user_id}",
		UsersUserIdDelete,
	},

	Route{
		"UsersUserIdGet",
		strings.ToUpper("Get"),
		"/api/v1/users/{user_id}",
		UsersUserIdGet,
	},

	Route{
		"UsersUserIdPut",
		strings.ToUpper("Put"),
		"/api/v1/users/{user_id}",
		UsersUserIdPut,
	},
}