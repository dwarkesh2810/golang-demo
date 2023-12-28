package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Englishlanguagelable_20231226_101842 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Englishlanguagelable_20231226_101842{}
	m.Created = "20231226_101842"

	migration.Register("Englishlanguagelable_20231226_101842", m)
}

// Run the migrations
func (m *Englishlanguagelable_20231226_101842) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`
	CREATE TABLE IF NOT EXISTS "english_language_lable" (
        "lang_id" serial NOT NULL PRIMARY KEY,
        "language_code" text NOT NULL DEFAULT '' ,
        "language_value" text NOT NULL DEFAULT '' ,
        "lable_code" text NOT NULL DEFAULT ''  UNIQUE,
        "section" text NOT NULL DEFAULT '' ,
        "created_by" integer NOT NULL DEFAULT 0 ,
        "updated_by" integer NOT NULL DEFAULT 0 ,
        "created_date" timestamp with time zone NOT NULL,
        "updated_date" timestamp with time zone
    );`)

}

// Reverse the migrations
func (m *Englishlanguagelable_20231226_101842) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
