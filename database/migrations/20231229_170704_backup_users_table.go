package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type BackupUsersTable_20231229_170704 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &BackupUsersTable_20231229_170704{}
	m.Created = "20231229_170704"

	migration.Register("BackupUsersTable_20231229_170704", m)
}

// Run the migrations
func (m *BackupUsersTable_20231229_170704) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE IF NOT EXISTS "backup_users" (
		"user_id" integer NOT NULL DEFAULT 0 ,
		"first_name" text NOT NULL DEFAULT '' ,
		"last_name" text NOT NULL DEFAULT '' ,
		"email" text NOT NULL DEFAULT '' ,
		"phone_number" text NOT NULL DEFAULT '' ,
		"password" text NOT NULL DEFAULT '' ,
		"isverified" integer NOT NULL DEFAULT 0 ,
		"otp_code" text NOT NULL DEFAULT '' ,
		"role" text NOT NULL DEFAULT '' ,
		"country_id" integer NOT NULL DEFAULT 0 ,
		"created_date" timestamp with time zone NOT NULL,
		"updated_date" timestamp with time zone NOT NULL,
		"delete_from_user" timestamp with time zone NOT NULL
	);`)
}

// Reverse the migrations
func (m *BackupUsersTable_20231229_170704) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
