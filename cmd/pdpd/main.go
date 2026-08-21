package main

import (
        "flag"
        "fmt"
        "log"
        "os"
        "os/signal"
        "syscall"

        abacpdp "github.com/LYH2263/go-abacpdp"
        "github.com/LYH2263/go-abacpdp/internal/api"
)

func main() {
        addr := flag.String("addr", ":8112", "listen address")
        web := flag.String("web", "web", "static web directory")
        flag.Parse()

        pdp := abacpdp.New()
        defer pdp.Close()
        if err := abacpdp.LoadSampleDocs(pdp); err != nil {
                log.Fatalf("load demo: %v", err)
        }

        srv := api.New(pdp, *web)
        go func() {
                log.Printf("pdpd listening on %s", *addr)
                if err := srv.ListenAndServe(*addr); err != nil {
                        log.Printf("server stopped: %v", err)
                }
        }()

        ch := make(chan os.Signal, 1)
        signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
        <-ch
        fmt.Println("shutting down")
}
