package msgcontrollers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/cloudfresco/vilom/testhelpers"
)

func TestGetMessage(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenstring := LoginUser()

	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://localhost:8000/v0.1/messages/89193ec7-469e-4580-8bce-e68ceb5aa201", nil)
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
	expected := string(`{"id":1,"id_s":"89193ec7-469e-4580-8bce-e68ceb5aa201","workspace_id":2,"channel_id":1,"user_id":1,"statusc":1,"created_at":"2019-07-23T10:04:26Z","updated_at":"2019-07-23T10:04:26Z","created_day":204,"created_week":30,"created_month":7,"created_year":2019,"updated_day":204,"updated_week":30,"updated_month":7,"updated_year":2019,"MessageTexts":[{"id":1,"mtext":"Hi. I am looking into buying a Floptical Drive, and was wondering what experience people have with the drives from Iomega, PLI, MASS MicroSystems, or Procom. These seem to be the main drives on the market. Any advice? Also, I heard about some article in MacWorld about Flopticals. Could someone post a summary, if they have it? Thanks in advance","workspace_id":2,"channel_id":1,"message_id":1,"user_id":1,"statusc":1,"created_at":"2019-07-23T10:04:26Z","updated_at":"2019-07-23T10:04:26Z","created_day":204,"created_week":30,"created_month":7,"created_year":2019,"updated_day":204,"updated_week":30,"updated_month":7,"updated_year":2019}],"MessageAttachments":[{"id":1,"mattach":"mattach","workspace_id":2,"channel_id":1,"message_id":1,"user_id":1,"statusc":1,"created_at":"2019-07-23T10:04:26Z","updated_at":"2019-07-23T10:04:26Z","created_day":204,"created_week":30,"created_month":7,"created_year":2019,"updated_day":204,"updated_week":30,"updated_month":7,"updated_year":2019}],"Mtext":"","Mattach":""}` + "\n")

	if w.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			w.Body.String(), expected)
		return
	}
}

func TestUpdateMessage(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenstring := LoginUser()

	w := httptest.NewRecorder()

	data := []byte(`{"m_text" : "Messagetext2"}`)

	req, err := http.NewRequest("PUT", "http://localhost:8000/v0.1/messages/89193ec7-469e-4580-8bce-e68ceb5aa201", bytes.NewBuffer(data))

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

func TestDeleteMessage(t *testing.T) {
	var err error
	err = testhelpers.LoadSQL(dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenstring := LoginUser()

	w := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", "http://localhost:8000/v0.1/messages/89193ec7-469e-4580-8bce-e68ceb5aa201", nil)

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
