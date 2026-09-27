package handler

import (
	"encoding/json"
	"itemsim-server/internal/application"
	"itemsim-server/internal/common/search/invindex"
	"itemsim-server/internal/config"
	"itemsim-server/internal/domain/soul"
	"itemsim-server/internal/infrastructure/repository/inmemory"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestSoulHandler_Search(t *testing.T) {
	dir := t.TempDir()
	data := `{"1":{"name":"Alpha soul","option":{"str":3}},"2":{"name":"Alpha soul","magnificent":true,"options":{"attackPower":{"attackPower":5}}},"3":{"name":"Beta soul","magnificent":true},"4":{"name":"Alpha soul","magnificent":false}}`
	if err := os.WriteFile(filepath.Join(dir, "soul.json"), []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	repository, err := inmemory.NewSoulRepository(&config.Config{ResourcesPath: dir})
	if err != nil {
		t.Fatal(err)
	}
	h := NewSoulHandler(application.NewSoulService(repository, invindex.NewSearcher[soul.Soul](repository.Count())))
	tests := []struct {
		name   string
		query  string
		status int
		ids    []int
	}{
		{"without filter", "query=Alpha", 200, []int{1, 2, 4}},
		{"magnificent", "query=Alpha&magnificent=true", 200, []int{2}},
		{"ordinary", "query=Alpha&magnificent=false", 200, []int{1, 4}},
		{"no matches", "query=Beta&magnificent=false", 200, []int{}},
		{"invalid filter", "query=Alpha&magnificent=invalid", 400, nil},
		{"missing query", "magnificent=true", 400, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.GET("/souls/search", h.Search)
			recorder := httptest.NewRecorder()
			e.ServeHTTP(recorder, httptest.NewRequest("GET", "/souls/search?"+tt.query, nil))
			if recorder.Code != tt.status {
				t.Fatalf("status = %d, want %d: %s", recorder.Code, tt.status, recorder.Body.String())
			}
			if tt.status != 200 {
				return
			}
			var results []application.SoulSearchResult
			if err := json.Unmarshal(recorder.Body.Bytes(), &results); err != nil {
				t.Fatal(err)
			}
			ids := make([]int, len(results))
			for i, result := range results {
				ids[i] = result.Id
			}
			if !reflect.DeepEqual(ids, tt.ids) {
				t.Errorf("ids = %v, want %v", ids, tt.ids)
			}
		})
	}
}
