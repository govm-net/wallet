package main

import (
	"expvar"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Route 单个http路由信息的结构体
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes 路由列表
type Routes []Route

// NewRouter 创建http路由
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

	router.Methods("GET").Path("/api/v1/").Name("Index").HandlerFunc(Index)
	router.Handle("/debug/vars", Logger(expvar.Handler(), "expvar"))
	router.Methods("GET").Name("static").Handler(http.FileServer(http.Dir("./static/")))

	return router
}

// Index api接口的默认的处理函数
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!\n")
	for _, route := range routes {
		fmt.Fprintf(w, "    %s: %s %s\n", route.Name, route.Method, route.Pattern)
	}
}

var routes = Routes{
	Route{
		"AccountGet",
		strings.ToUpper("Get"),
		"/api/v1/{chain}/account",
		AccountGet,
	},

	Route{
		"TransactionMovePost",
		strings.ToUpper("Post"),
		"/api/v1/{chain}/transaction/move",
		TransactionMovePost,
	},

	Route{
		"TransactionTransferPost",
		strings.ToUpper("Post"),
		"/api/v1/{chain}/transaction/transfer",
		TransactionTransferPost,
	},

	Route{
		"TransactionVotePost",
		strings.ToUpper("Post"),
		"/api/v1/{chain}/transaction/vote",
		TransactionVotePost,
	},

	Route{
		"TransactionVoteDelete",
		strings.ToUpper("Delete"),
		"/api/v1/{chain}/transaction/vote",
		TransactionVoteDelete,
	},

	Route{
		"TransactionMinerPost",
		strings.ToUpper("Post"),
		"/api/v1/{chain}/transaction/miner",
		TransactionMinerPost,
	},

	Route{
		"TransactionRunAppPost",
		strings.ToUpper("Post"),
		"/api/v1/{chain}/transaction/app/run",
		TransactionRunAppPost,
	},

	Route{
		"TransactionAppLifePost",
		strings.ToUpper("Post"),
		"/api/v1/{chain}/transaction/app/life",
		TransactionAppLifePost,
	},

	Route{
		"TransactionAppInfoGet",
		strings.ToUpper("Get"),
		"/api/v1/{chain}/transaction/app/info",
		proxyHTTP,
	},

	Route{
		"TransactionInfoGet",
		strings.ToUpper("Get"),
		"/api/v1/{chain}/transaction/info",
		proxyHTTP,
	},

	Route{
		"BlockInfoGet",
		strings.ToUpper("Get"),
		"/api/v1/{chain}/block/info",
		proxyHTTP,
	},

	Route{
		"DataRead",
		strings.ToUpper("Get"),
		"/api/v1/{chain}/data",
		proxyHTTP,
	},

	Route{
		"DataNextKeyGet",
		strings.ToUpper("Get"),
		"/api/v1/{chain}/data/visit",
		proxyHTTP,
	},

	Route{
		"AdminsGet",
		strings.ToUpper("Get"),
		"/api/v1/{chain}/admins",
		AdminsGet,
	},

	Route{
		"AdminInfoGet",
		strings.ToUpper("Get"),
		"/api/v1/{chain}/admin",
		AdminInfoGet,
	},
	Route{
		"VersionGet",
		strings.ToUpper("Get"),
		"/api/v1/version",
		VersionGet,
	},
}
