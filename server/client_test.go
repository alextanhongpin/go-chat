package server_test

import (
	"log"
	"testing"

	"github.com/alextanhongpin/go-chat/server"
)

func TestMapper(t *testing.T) {
	mapper := server.NewMapper()
	mapper.Add("a", "b")

	if _, found := mapper.Data()["a"]["b"]; !found {
		t.Fatal("not found")
	}
	log.Println(mapper)
	mapper.Delete("a", "b")
	if _, found := mapper.Data()["a"]["b"]; found {
		t.Fatal("not found")
	}
	log.Println(mapper)
}
