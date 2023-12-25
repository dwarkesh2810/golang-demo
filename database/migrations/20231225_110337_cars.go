package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Cars_20231225_110337 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Cars_20231225_110337{}
	m.Created = "20231225_110337"

	migration.Register("Cars_20231225_110337", m)
}

// Run the migrations
func (m *Cars_20231225_110337) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE IF NOT EXISTS "car" (
		"id" serial NOT NULL PRIMARY KEY,
		"car_name" text NOT NULL DEFAULT '' ,
		"car_image" text,
		"modified_by" text NOT NULL DEFAULT '' ,
		"model" text NOT NULL DEFAULT '' ,
		"car_type" text NOT NULL DEFAULT '' ,
		"ctreated_date" timestamp with time zone,
		"updated_date" timestamp with time zone,
		"created_by" int not null,
		"updated_by" int
	);`)
}

// Reverse the migrations
func (m *Cars_20231225_110337) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
