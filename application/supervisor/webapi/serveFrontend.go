package webapi

import "net/http"

func RunFrontendServer() {
	http.Handle("/", http.FileServer(http.Dir("./web/out")))
	http.ListenAndServe(":3000", nil)
}
