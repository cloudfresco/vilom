package msgcontrollers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/cloudfresco/vilom/testhelpers"
)

func TestShowTopic(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenstring := LoginUser()

	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://localhost:8000/v0.1/topics/44b2e674-7031-4487-be96-60093bfe8ac3", nil)
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
	expected := string(`{"id":1,"id_s":"44b2e674-7031-4487-be96-60093bfe8ac3","topic_name":"Floptical Question","topic_desc":"Floptical Question","num_messages":1,"category_id":2,"user_id":1,"statusc":1,"created_at":"2019-07-23T10:04:26Z","updated_at":"2019-07-23T10:04:26Z","created_day":204,"created_week":30,"created_month":7,"created_year":2019,"updated_day":204,"updated_week":30,"updated_month":7,"updated_year":2019,"Messages":[{"id":1,"id_s":"89193ec7-469e-4580-8bce-e68ceb5aa201","category_id":2,"topic_id":1,"user_id":1,"statusc":1,"created_at":"2019-07-23T10:04:26Z","updated_at":"2019-07-23T10:04:26Z","created_day":204,"created_week":30,"created_month":7,"created_year":2019,"updated_day":204,"updated_week":30,"updated_month":7,"updated_year":2019,"MessageTexts":[{"id":1,"mtext":"Hi. I am looking into buying a Floptical Drive, and was wondering what experience people have with the drives from Iomega, PLI, MASS MicroSystems, or Procom. These seem to be the main drives on the market. Any advice? Also, I heard about some article in MacWorld about Flopticals. Could someone post a summary, if they have it? Thanks in advance","category_id":2,"topic_id":1,"message_id":1,"user_id":1,"statusc":1,"created_at":"2019-07-23T10:04:26Z","updated_at":"2019-07-23T10:04:26Z","created_day":204,"created_week":30,"created_month":7,"created_year":2019,"updated_day":204,"updated_week":30,"updated_month":7,"updated_year":2019}],"MessageAttachments":[{"id":1,"mattach":"mattach","category_id":2,"topic_id":1,"message_id":1,"user_id":1,"statusc":1,"created_at":"2019-07-23T10:04:26Z","updated_at":"2019-07-23T10:04:26Z","created_day":204,"created_week":30,"created_month":7,"created_year":2019,"updated_day":204,"updated_week":30,"updated_month":7,"updated_year":2019}],"Mtext":"","Mattach":""}]}` + "\n")

	if w.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			w.Body.String(), expected)
		return
	}

}

func TestGetTopicByName(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenstring := LoginUser()

	w := httptest.NewRecorder()

	data := []byte(`{"topic_name" : "Floptical Question"}`)
	req, err := http.NewRequest("POST", "http://localhost:8000/v0.1/topics/topicbyname", bytes.NewBuffer(data))

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
	expected := string(`{"id":1,"id_s":"44b2e674-7031-4487-be96-60093bfe8ac3","topic_name":"Floptical Question","topic_desc":"Floptical Question","num_messages":1,"category_id":2,"user_id":1,"statusc":1,"created_at":"2019-07-23T10:04:26Z","updated_at":"2019-07-23T10:04:26Z","created_day":204,"created_week":30,"created_month":7,"created_year":2019,"updated_day":204,"updated_week":30,"updated_month":7,"updated_year":2019,"Messages":null}` + "\n")
	if w.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			w.Body.String(), expected)
		return
	}
}

func TestUpdateTopic(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenstring := LoginUser()

	w := httptest.NewRecorder()

	data := []byte(`{"topic_name" : "topic2", "topic_desc" : "topic2 description"}`)

	req, err := http.NewRequest("PUT", "http://localhost:8000/v0.1/topics/44b2e674-7031-4487-be96-60093bfe8ac3", bytes.NewBuffer(data))

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

func TestDeleteTopic(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenstring := LoginUser()

	w := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", "http://localhost:8000/v0.1/topics/44b2e674-7031-4487-be96-60093bfe8ac3", nil)

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
