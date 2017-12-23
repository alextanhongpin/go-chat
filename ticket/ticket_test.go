package ticket

import (
	"log"
	"strings"
	"testing"
	"time"
)

func TestNewTicket(t *testing.T) {
	id := "abc123"
	tic := New(id, 1)

	want := id
	got := tic.ID

	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}
func TestExpireTicket(t *testing.T) {

	tic := New("abc123", 1*time.Second)
	ticketStr, err := Sign(tic)
	if err != nil {
		t.Error(err.Error())
	}
	// Sleep for 2 seconds to ensure ticket is expired
	time.Sleep(2 * time.Second)

	_, err = Verify(ticketStr)
	log.Println(err)
	want := "token is expired by"
	got := err.Error()
	if !strings.Contains(got, want) {
		t.Errorf(`want "%v", got "%v"`, want, got)
	}
}
