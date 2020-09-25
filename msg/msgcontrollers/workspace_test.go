package msgcontrollers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/cloudfresco/vilom/testhelpers"
)

func TestGetWorkspaces(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenstring := LoginUser()

	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://localhost:8000/v0.1/workspaces", bytes.NewBuffer([]byte(`{"limit": 20}, "cursor": ""`)))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+tokenstring)

	mux.ServeHTTP(w, req)

	resp := w.Result()
	// Check the status code is what we expect.
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code %d", resp.StatusCode)
		return
	}
	expected := string(`{"Workspaces":[{"id":1,"id_s":"1bd1888a-dbfe-4510-a7ad-a98f69fd0a6b","workspace_name":"Performance Portable Transmitter","workspace_desc":"Performance Portable Transmitter","num_chd":1,"user_id":1,"statusc":1,"created_at":"2019-07-23T10:04:26Z","updated_at":"2019-07-23T10:04:26Z","created_day":204,"created_week":30,"created_month":7,"created_year":2019,"updated_day":204,"updated_week":30,"updated_month":7,"updated_year":2019,"Channels":null}],"next_cursor":"MA=="}` + "\n")

	if w.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			w.Body.String(), expected)
		return
	}
}

func TestGetWorkspaceWithChannels(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenstring := LoginUser()

	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://localhost:8000/v0.1/workspaces/1c29bf3a-4684-499c-a519-2c348aa13246", nil)

	if err != nil {
		t.Error(err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+tokenstring)
	mux.ServeHTTP(w, req)

	resp := w.Result()
	// Check the status code is what we expect.
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code %d", resp.StatusCode)
		return
	}
	expected := string(`{"id":2,"id_s":"1c29bf3a-4684-499c-a519-2c348aa13246","workspace_name":"Drive","workspace_desc":"Drive","num_channels":1,"levelc":1,"parent_id":1,"user_id":1,"statusc":1,"created_at":"2019-07-23T10:04:26Z","updated_at":"2019-07-23T10:04:26Z","created_day":204,"created_week":30,"created_month":7,"created_year":2019,"updated_day":204,"updated_week":30,"updated_month":7,"updated_year":2019,"Channels":[{"id":1,"id_s":"44b2e674-7031-4487-be96-60093bfe8ac3","channel_name":"Floptical Question","channel_desc":"Floptical Question","num_messages":1,"workspace_id":2,"user_id":1,"statusc":1,"created_at":"2019-07-23T10:04:26Z","updated_at":"2019-07-23T10:04:26Z","created_day":204,"created_week":30,"created_month":7,"created_year":2019,"updated_day":204,"updated_week":30,"updated_month":7,"updated_year":2019,"Messages":null}]}` + "\n")
	if w.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			w.Body.String(), expected)
		return
	}
}

func TestGetTopLevelWorkspaces(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenstring := LoginUser()

	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://localhost:8000/v0.1/workspaces/topworkspaces", nil)

	if err != nil {
		t.Error(err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+tokenstring)
	mux.ServeHTTP(w, req)

	resp := w.Result()
	// Check the status code is what we expect.
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code %d", resp.StatusCode)
		return
	}

	expected := string(`[{"id":1,"id_s":"1bd1888a-dbfe-4510-a7ad-a98f69fd0a6b","workspace_name":"Performance Portable Transmitter","workspace_desc":"Performance Portable Transmitter","num_chd":1,"user_id":1,"statusc":1,"created_at":"2019-07-23T10:04:26Z","updated_at":"2019-07-23T10:04:26Z","created_day":204,"created_week":30,"created_month":7,"created_year":2019,"updated_day":204,"updated_week":30,"updated_month":7,"updated_year":2019,"Channels":null}]` + "\n")

	if w.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			w.Body.String(), expected)
		return
	}
}

func TestGetChildWorkspaces(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenstring := LoginUser()

	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://localhost:8000/v0.1/workspaces/1bd1888a-dbfe-4510-a7ad-a98f69fd0a6b/chdn", nil)

	if err != nil {
		t.Error(err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+tokenstring)
	mux.ServeHTTP(w, req)

	resp := w.Result()
	// Check the status code is what we expect.
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code %d", resp.StatusCode)
		return
	}

	expected := string(`[{"id":2,"id_s":"1c29bf3a-4684-499c-a519-2c348aa13246","workspace_name":"Drive","workspace_desc":"Drive","num_channels":1,"levelc":1,"parent_id":1,"user_id":1,"statusc":1,"created_at":"2019-07-23T10:04:26Z","updated_at":"2019-07-23T10:04:26Z","created_day":204,"created_week":30,"created_month":7,"created_year":2019,"updated_day":204,"updated_week":30,"updated_month":7,"updated_year":2019,"Channels":null}]` + "\n")

	if w.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			w.Body.String(), expected)
		return
	}
}

func TestGetParentWorkspace(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenstring := LoginUser()

	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://localhost:8000/v0.1/workspaces/1c29bf3a-4684-499c-a519-2c348aa13246/getparent", nil)

	if err != nil {
		t.Error(err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+tokenstring)
	mux.ServeHTTP(w, req)

	resp := w.Result()
	// Check the status code is what we expect.
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code %d", resp.StatusCode)
		return
	}

	expected := string(`{"id":1,"id_s":"1bd1888a-dbfe-4510-a7ad-a98f69fd0a6b","workspace_name":"Performance Portable Transmitter","workspace_desc":"Performance Portable Transmitter","num_chd":1,"user_id":1,"statusc":1,"created_at":"2019-07-23T10:04:26Z","updated_at":"2019-07-23T10:04:26Z","created_day":204,"created_week":30,"created_month":7,"created_year":2019,"updated_day":204,"updated_week":30,"updated_month":7,"updated_year":2019,"Channels":null}` + "\n")

	if w.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			w.Body.String(), expected)
		return
	}
}

func TestUpdateWorkspace(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenstring := LoginUser()

	w := httptest.NewRecorder()

	data := []byte(`{"workspace_name" : "workspace3", "workspace_desc" : "workspace3 description"}`)

	req, err := http.NewRequest("PUT", "http://localhost:8000/v0.1/workspaces/1c29bf3a-4684-499c-a519-2c348aa13246", bytes.NewBuffer(data))

	if err != nil {
		t.Error(err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+tokenstring)
	mux.ServeHTTP(w, req)

	resp := w.Result()
	// Check the status code is what we expect.
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code %d", resp.StatusCode)
		return
	}

	expected := string(`"Updated Successfully"` + "\n")

	if w.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			w.Body.String(), expected)
		return
	}
}

func TestDeleteWorkspace(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenstring := LoginUser()

	w := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", "http://localhost:8000/v0.1/workspaces/1c29bf3a-4684-499c-a519-2c348aa13246", nil)

	if err != nil {
		t.Error(err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+tokenstring)
	mux.ServeHTTP(w, req)

	resp := w.Result()
	// Check the status code is what we expect.
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code %d", resp.StatusCode)
		return
	}

	expected := string(`"Deleted Successfully"` + "\n")

	if w.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			w.Body.String(), expected)
		return
	}
}
