package captive

import (
    "flag"
    "fmt"
    "log"
    "net/http"
)

func run() {
    cert := flag.String("cert", "", "cert file")
    key := flag.String("key", "", "private key file")
    port := flag.Uint("port", 443, "listen port")

    flag.Parse()

    var lts bool

    if *cert == "" || *key == "" {
        lts = false
    } else {
        lts = true
    }

    http.HandleFunc("/", generate)

    switch lts {
    case false:
        log.Println("recommend use lts mode")
        log.Printf("listen on [::]:%d\n", *port)
        log.Fatal(http.ListenAndServe(fmt.Sprintf("[::]:%d", *port), nil))

    case true:
        log.Printf("listen on [::]:%d\n", *port)
        log.Fatal(http.ListenAndServeTLS(fmt.Sprintf("[::]:%d", *port), *cert, *key, nil))
    }
}

func generate(w http.ResponseWriter, _ *http.Request) {
    w.WriteHeader(http.StatusNoContent)
}
