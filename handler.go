package main_handler

import (
    "fmt"
    "time"
    "math/rand"
    "net/http"
    "strconv"
    "appengine"
    "appengine/datastore"
)

const (
    c_err_missing_data   = 0
    c_err_db_put_fail    = 1
    c_err_db_query_fail  = 2
)

type RegNode struct {
    UserAgent string
    DeviceId  string
    Status    string
    Battery   int
    Cycles    int
    Date      time.Time
}

func RegRootKey(c appengine.Context) *datastore.Key {
    return datastore.NewKey(c, "Registrations", "iot_registrations", 0, nil)
}

func init() {
    rand.Seed(24);
    http.HandleFunc("/", top_handler)
    http.HandleFunc("/reg", reg_handler)
    http.HandleFunc("/lst", lst_handler)
}

func top_handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    c := appengine.NewContext(r)

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
    
    did := r.URL.Query().Get("did")
    if len(did) == 0 {
        rp_Error(w, c_err_missing_data)
        return;
    }
    
    battery, _ := strconv.Atoi(r.URL.Query().Get("bat"))
    cycles, _ := strconv.Atoi(r.URL.Query().Get("cyc"))
   
    reg := RegNode {
        UserAgent : r.UserAgent()[0:11],
        DeviceId  : did,
        Battery   : battery,
        Cycles    : cycles,
        Status    : r.URL.Query().Get("sta"),
        Date      : time.Now(),
    }
    
    c := appengine.NewContext(r)

    _, err := datastore.Put(
            c,
            datastore.NewIncompleteKey(c, "RegNode", RegRootKey(c)),
            &reg)

    if err != nil {
        rp_Error(w, c_err_db_put_fail)
        return
    }

    rp_RegOk(w, rand.Intn(9))
}

func lst_handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    c := appengine.NewContext(r)
    q := datastore.NewQuery("RegNode").Ancestor(RegRootKey(c)).Order("-Date").Limit(20)
    registrations := make([]RegNode, 0, 20)
    _, err := q.GetAll(c, &registrations)
    if err != nil {
        rp_Error(w, c_err_db_query_fail)
        return
    }
    
    for ix, v := range registrations {
        fmt.Fprintln(w, ix, "->", v)
    }
    
    fmt.Fprintln(w, ".done.")
}

func rp_Error(w http.ResponseWriter, code int32) {
    fmt.Fprintln(w, "@ error=", code); 
}

func rp_RegOk(w http.ResponseWriter, led_level int) {
    tt := time.Now().UTC()
    fmt.Fprintln(w, "IoT 0001");
    fmt.Fprintln(w, tt)
    fmt.Fprintln(w, "ledbar", led_level)
}

