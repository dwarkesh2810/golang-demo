package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Auditlogs_20231225_101008 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Auditlogs_20231225_101008{}
	m.Created = "20231225_101008"

	migration.Register("Auditlogs_20231225_101008", m)
}

// Run the migrations
func (m *Auditlogs_20231225_101008) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE IF NOT EXISTS "audit_logs" (
        "log_id" serial NOT NULL PRIMARY KEY,
        "user_id" bigint CHECK("user_id" >= 0) NOT NULL DEFAULT 0 ,
        "action" text NOT NULL DEFAULT '' ,
        "user_ip" text NOT NULL DEFAULT '' ,
        "discription" text NOT NULL DEFAULT '' ,
        "end_points" text NOT NULL DEFAULT '' ,
        "created_date" timestamp with time zone NOT NULL
    );`)
}

// Reverse the migrations
func (m *Auditlogs_20231225_101008) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
