//nolint:forbidigo // fmt.Printf is allowed here
package webapi

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func dumpRequest(req *http.Request) {
	dump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("\n\nREQUEST %s %s\n--------------\n", req.Method, req.URL.String())
	fmt.Println(string(dump))
	fmt.Println("--------------")
}

func dumpResponse(req *http.Request, resp *http.Response) {
	dump, _ := httputil.DumpResponse(resp, true)
	//nolint:forbidigo
	fmt.Printf("\n\nRESPONSE %s %s\n--------------\n", req.Method, req.URL.String())
	fmt.Println(string(dump))
	fmt.Println("--------------")
}
