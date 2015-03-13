package main_handler

import (
    "fmt"
    "time"
    "math/rand"
    "net/http"
    "appengine"
)

func init() {
    rand.Seed(24);
    http.HandleFunc("/", top_handler)
    http.HandleFunc("/reg", reg_handler)
}

func top_handler(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")

    fmt.Fprintf(w, "= carlos.pizano@gmail.com \n")
    fmt.Fprintf(w, "= appengine app id = %q\n", appengine.AppID(c))
    fmt.Fprintf(w, "= appengine version = %q\n", appengine.VersionID(c))
    
    huri :=  "http://" + r.Host
    fmt.Fprintf(w, "= host is %s\n", huri)

    name, index := appengine.BackendInstance(c)
    fmt.Fprintf(w, "= appengine backendinstance = %q, %d\n", name, index)
}

func reg_handler(w http.ResponseWriter, r *http.Request) {   
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    tt := time.Now().UTC()
    fmt.Fprintln(w, "IoT 0001");
    fmt.Fprintln(w, tt)
    fmt.Fprintln(w, "ledbar", rand.Intn(9))
}

