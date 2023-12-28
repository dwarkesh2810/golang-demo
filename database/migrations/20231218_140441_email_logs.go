package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type EmailLogs_20231218_140441 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &EmailLogs_20231218_140441{}
	m.Created = "20231218_140441"

	migration.Register("EmailLogs_20231218_140441", m)
}

// Run the migrations
func (m *EmailLogs_20231218_140441) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE IF NOT EXISTS "public"."email_logs" (
		"LogId" bigserial NOT NULL,
		"emailTo" text NOT NULL DEFAULT '',
		"name" text NOT NULL DEFAULT '',
		"subject" text NOT NULL DEFAULT '',
		"body" text NOT NULL DEFAULT '',
		"status" text NOT NULL DEFAULT '',
		PRIMARY KEY ("LogId")
	);`)
}

// Reverse the migrations
func (m *EmailLogs_20231218_140441) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
