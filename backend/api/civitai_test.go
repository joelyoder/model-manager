package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func redirectToServer(srv *httptest.Server) func() {
	u, _ := url.Parse(srv.URL)
	orig := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		req.URL.Scheme = u.Scheme
		req.URL.Host = u.Host
		return http.DefaultTransport.RoundTrip(req)
	})
	return func() { http.DefaultClient.Transport = orig }
}

func TestFetchCivitModels(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		models := []CivitModel{{ID: 1, Name: "foo"}}
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/v1/models" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			if err := json.NewEncoder(w).Encode(models); err != nil {
				t.Fatalf("encode: %v", err)
			}
		}))
		defer srv.Close()
		cleanup := redirectToServer(srv)
		defer cleanup()

		got, err := FetchCivitModels("token")
		if err != nil {
			t.Fatalf("FetchCivitModels: %v", err)
		}
		if len(got) != 1 || got[0].ID != models[0].ID || got[0].Name != models[0].Name {
			t.Fatalf("unexpected result: %+v", got)
		}
	})

	t.Run("http error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()
		cleanup := redirectToServer(srv)
		defer cleanup()

		if _, err := FetchCivitModels("token"); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("bad json", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("not json"))
		}))
		defer srv.Close()
		cleanup := redirectToServer(srv)
		defer cleanup()

		if _, err := FetchCivitModels("token"); err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestFetchCivitModel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		model := CivitModel{ID: 2, Name: "bar"}
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/v1/models/2" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			json.NewEncoder(w).Encode(model)
		}))
		defer srv.Close()
		cleanup := redirectToServer(srv)
		defer cleanup()

		got, err := FetchCivitModel("token", 2)
		if err != nil {
			t.Fatalf("FetchCivitModel: %v", err)
		}
		if got.ID != model.ID || got.Name != model.Name {
			t.Fatalf("unexpected model: %+v", got)
		}
	})

	t.Run("http error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()
		cleanup := redirectToServer(srv)
		defer cleanup()

		if _, err := FetchCivitModel("token", 3); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("bad json", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("not json"))
		}))
		defer srv.Close()
		cleanup := redirectToServer(srv)
		defer cleanup()

		if _, err := FetchCivitModel("token", 4); err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestFetchModelVersion(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		version := VersionResponse{ID: 5, ModelID: 2, Name: "v"}
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/v1/model-versions/5" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			json.NewEncoder(w).Encode(version)
		}))
		defer srv.Close()
		cleanup := redirectToServer(srv)
		defer cleanup()

		got, err := FetchModelVersion("token", 5)
		if err != nil {
			t.Fatalf("FetchModelVersion: %v", err)
		}
		if got.ID != version.ID || got.Name != version.Name || got.ModelID != version.ModelID {
			t.Fatalf("unexpected version: %+v", got)
		}
	})

	t.Run("http error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer srv.Close()
		cleanup := redirectToServer(srv)
		defer cleanup()

		if _, err := FetchModelVersion("token", 6); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("bad json", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("not json"))
		}))
		defer srv.Close()
		cleanup := redirectToServer(srv)
		defer cleanup()

		if _, err := FetchModelVersion("token", 7); err == nil {
			t.Fatal("expected error")
		}
	})
}
