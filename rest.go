//  Copyright (c) 2015 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the
//  License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing,
//  software distributed under the License is distributed on an "AS
//  IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
//  express or implied. See the License for the specific language
//  governing permissions and limitations under the License.

package cbft

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/couchbaselabs/cbgt"
	"github.com/couchbaselabs/cbgt/rest"
)

func InitStaticRouter(staticDir, staticETag string) *mux.Router {
	hfsStaticX := http.FileServer(assetFS())

	router := mux.NewRouter()
	router.StrictSlash(true)

	router.Handle("/",
		http.RedirectHandler("/staticx/index.html", 302))
	router.Handle("/index.html",
		http.RedirectHandler("/staticx/index.html", 302))
	router.Handle("/static/partials/index/list.html",
		http.RedirectHandler("/staticx/partials/index/list.html", 302))

	router = rest.InitStaticRouter(router,
		staticDir, staticETag, []string{
			"/indexes",
			"/nodes",
			"/monitor",
			"/manage",
			"/logs",
			"/debug",
		}, http.RedirectHandler("/staticx/index.html", 302))

	router.PathPrefix("/staticx/").Handler(
		http.StripPrefix("/staticx/", hfsStaticX))

	return router
}

func myAssetDir(name string) ([]string, error) {
	a, err := AssetDir(name)
	if err == nil {
		return a, err
	}

	return rest.AssetDir(name)
}

func myAsset(name string) ([]byte, error) {
	b, err := Asset(name)
	if err == nil {
		return b, err
	}

	return rest.Asset(name)
}

// NewRESTRouter creates a mux.Router initialized with the REST
// API and web UI routes.  See also InitStaticRouter if you need finer
// control of the router initialization.
func NewRESTRouter(versionMain string, mgr *cbgt.Manager,
	staticDir, staticETag string, mr *cbgt.MsgRing) (
	*mux.Router, map[string]rest.RESTMeta, error) {
	return rest.InitRESTRouter(InitStaticRouter(staticDir, staticETag),
		versionMain, mgr, staticDir, staticETag, mr,
		myAssetDir, myAsset)
}
