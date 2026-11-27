package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testServer(t *testing.T) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(withCORS(newMux(getDataset())))
	t.Cleanup(server.Close)
	return server
}

func getJSON(t *testing.T, url string, target any) *http.Response {
	t.Helper()
	response, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	t.Cleanup(func() { response.Body.Close() })
	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		t.Fatalf("decode %s: %v", url, err)
	}
	return response
}

func TestEntitiesEndpoint(t *testing.T) {
	server := testServer(t)
	var entities []Entity
	response := getJSON(t, server.URL+"/api/entities", &entities)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.StatusCode)
	}
	if len(entities) != len(getDataset()) {
		t.Fatalf("got %d entities, want %d", len(entities), len(getDataset()))
	}
	if entities[0].History[0].Density == 0 {
		t.Fatal("density was not calculated for the first entity")
	}
}

func TestEntitiesTypeFilter(t *testing.T) {
	server := testServer(t)
	var cities []Entity
	getJSON(t, server.URL+"/api/entities?type=city", &cities)
	if len(cities) == 0 {
		t.Fatal("expected at least one city")
	}
	for _, city := range cities {
		if city.Type != "city" {
			t.Fatalf("entity %q has type %q, want city", city.ID, city.Type)
		}
	}
}

func TestEntitiesSearch(t *testing.T) {
	server := testServer(t)
	var entities []Entity
	getJSON(t, server.URL+"/api/entities?q=%D0%9C%D0%BE%D1%81%D0%BA%D0%B2%D0%B0", &entities) // "Москва"
	if len(entities) == 0 {
		t.Fatal("search by name returned no results")
	}
	for _, entity := range entities {
		if entity.Name != "Москва" && entity.Country != "Москва" {
			t.Fatalf("unexpected search result: %q", entity.Name)
		}
	}
}

func TestEntityByID(t *testing.T) {
	server := testServer(t)
	var entity Entity
	response := getJSON(t, server.URL+"/api/entities/msk", &entity)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.StatusCode)
	}
	if entity.Name != "Москва" {
		t.Fatalf("name = %q, want Москва", entity.Name)
	}
}

func TestEntityNotFound(t *testing.T) {
	server := testServer(t)
	response, err := http.Get(server.URL + "/api/entities/unknown")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", response.StatusCode)
	}
}

func TestMetaEndpoint(t *testing.T) {
	server := testServer(t)
	var meta struct {
		Years     []int `json:"years"`
		Count     int   `json:"count"`
		Countries int   `json:"countries"`
		Cities    int   `json:"cities"`
	}
	getJSON(t, server.URL+"/api/meta", &meta)
	if meta.Count != meta.Countries+meta.Cities {
		t.Fatalf("count %d != countries %d + cities %d", meta.Count, meta.Countries, meta.Cities)
	}
	if len(meta.Years) < 2 {
		t.Fatalf("expected several years, got %v", meta.Years)
	}
}
