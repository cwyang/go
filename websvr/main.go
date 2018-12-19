package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

func getEnv(k, fallback string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	return fallback
}

func getListenPort() string {
	port := getEnv("PORT", "80")
	return ":" + port
}

func logSetup() {
	log.Printf("ListenPort: %s\n", getListenPort())
}
func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
}
func handleRedirect(w http.ResponseWriter, r *http.Request) {
	r.URL.Host="http://172.16.202.30" + getListenPort()
	title := r.URL.Path[len("/r/"):]
	code, _ := strconv.Atoi(title)
	http.Redirect(w, r, r.URL.Host, code)
}
func handleQ(w http.ResponseWriter, r *http.Request) {
	const tpl = `
<form action="/r/{{.Code}}" method={{.Method}}>
<input type="submit" value="{{.Method}}-{{.Code}}">
</form>`
	t, _ := template.New("webpage").Parse(tpl)
	fmt.Fprintf(w, "<html>")
	t.Execute(w, struct {Code int; Method string}{Code: 301, Method: "GET"})
	t.Execute(w, struct {Code int; Method string}{Code: 302, Method: "GET"})
	t.Execute(w, struct {Code int; Method string}{Code: 303, Method: "GET"})
	t.Execute(w, struct {Code int; Method string}{Code: 307, Method: "GET"})
	t.Execute(w, struct {Code int; Method string}{Code: 308, Method: "GET"})
	t.Execute(w, struct {Code int; Method string}{Code: 301, Method: "POST"})
	t.Execute(w, struct {Code int; Method string}{Code: 302, Method: "POST"})
	t.Execute(w, struct {Code int; Method string}{Code: 303, Method: "POST"})
	t.Execute(w, struct {Code int; Method string}{Code: 307, Method: "POST"})
	t.Execute(w, struct {Code int; Method string}{Code: 308, Method: "POST"})
	fmt.Fprintf(w, "</html>")
}

func main() {
	logSetup()
	http.HandleFunc("/q", handleQ)
	http.HandleFunc("/", handle)
	http.HandleFunc("/r/", handleRedirect)
	err := http.ListenAndServe(getListenPort(), nil)
	if (err != nil) {
		log.Fatal("ListenAndServe: ", err)
	}
}
