package webapi

import "net/http"

func RunFrontendServer() {
	http.Handle("/", http.FileServer(http.Dir("./web/out")))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
