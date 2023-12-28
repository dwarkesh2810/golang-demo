package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Multilanguagelable_20231226_101918 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Multilanguagelable_20231226_101918{}
	m.Created = "20231226_101918"

	migration.Register("Multilanguagelable_20231226_101918", m)
}

// Run the migrations
func (m *Multilanguagelable_20231226_101918) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`  CREATE TABLE IF NOT EXISTS "multi_language_lable" (
        "lable_id" serial NOT NULL PRIMARY KEY,
        "lable_code" text NOT NULL DEFAULT '' ,
        "language_value" text NOT NULL DEFAULT '' ,
        "language_code" text NOT NULL DEFAULT '' ,
        "section" text NOT NULL DEFAULT '' ,
        "created_by" integer NOT NULL DEFAULT 0 ,
        "updated_by" integer NOT NULL DEFAULT 0 ,
        "created_date" timestamp with time zone NOT NULL,
        "updated_date" timestamp with time zone
    );`)

}

// Reverse the migrations
func (m *Multilanguagelable_20231226_101918) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
