package api

import (
        "encoding/json"
        "net/http"
        "time"

        abacpdp "github.com/LYH2263/go-abacpdp"
)

type Server struct {
        pdp    *abacpdp.PDP
        webDir string
        mux    *http.ServeMux
}

func New(pdp *abacpdp.PDP, webDir string) *Server {
        s := &Server{pdp: pdp, webDir: webDir, mux: http.NewServeMux()}
        s.routes()
        return s
}

func (s *Server) routes() {
        s.mux.HandleFunc("/api/health", s.handleHealth)
        s.mux.HandleFunc("/api/stats", s.handleStats)
        s.mux.HandleFunc("/api/policy", s.handlePolicy)
        s.mux.HandleFunc("/api/evaluate", s.handleEvaluate)
        s.mux.HandleFunc("/api/load", s.handleLoad)
        if s.webDir != "" {
                s.mux.Handle("/", http.FileServer(http.Dir(s.webDir)))
        }
}

func (s *Server) Handler() http.Handler { return s.mux }

func (s *Server) ListenAndServe(addr string) error {
        srv := &http.Server{Addr: addr, Handler: s.mux, ReadHeaderTimeout: 5 * time.Second}
        return srv.ListenAndServe()
}

func writeJSON(w http.ResponseWriter, code int, v any) {
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        w.WriteHeader(code)
        _ = json.NewEncoder(w).Encode(v)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
        writeJSON(w, 200, map[string]any{"ok": !s.pdp.IsClosed()})
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
        writeJSON(w, 200, s.pdp.Stats())
}

func (s *Server) handlePolicy(w http.ResponseWriter, r *http.Request) {
        writeJSON(w, 200, s.pdp.PolicyView())
}

func (s *Server) handleEvaluate(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                http.Error(w, "method", 405)
                return
        }
        var bag abacpdp.AttrBag
        if err := json.NewDecoder(r.Body).Decode(&bag); err != nil {
                http.Error(w, err.Error(), 400)
                return
        }
        dec, err := s.pdp.EvaluateContext(r.Context(), bag)
        if err != nil {
                http.Error(w, err.Error(), 400)
                return
        }
        writeJSON(w, 200, dec)
}

func (s *Server) handleLoad(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                http.Error(w, "method", 405)
                return
        }
        var raw json.RawMessage
        if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
                http.Error(w, err.Error(), 400)
                return
        }
        if err := s.pdp.LoadJSON(raw); err != nil {
                http.Error(w, err.Error(), 400)
                return
        }
        writeJSON(w, 200, s.pdp.PolicyView())
}
