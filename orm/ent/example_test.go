package ent

import (
	"context"
	"entgo.io/ent/dialect"
	"gitee.com/geektime-geekbang/geektime-go/orm/ent/ent"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"testing"
)

func TestEntCURD(t *testing.T) {
	// Create an ent.Client with in-memory SQLite database.
	client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	ctx := context.Background()
	// Run the automatic migration tool to create all schema resources.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	err = client.User.Create().Exec(context.Background())
}
