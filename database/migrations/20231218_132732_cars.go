package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Cars_20231218_132732 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Cars_20231218_132732{}
	m.Created = "20231218_132732"

	migration.Register("Cars_20231218_132732", m)
}

// Run the migrations
func (m *Cars_20231218_132732) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE IF NOT EXISTS "car" (
		"id" serial NOT NULL PRIMARY KEY,
		"car_name" text NOT NULL DEFAULT '' ,
		"car_image" text,
		"modified_by" text NOT NULL DEFAULT '' ,
		"model" text NOT NULL DEFAULT '' ,
		"car_type" text NOT NULL DEFAULT '' ,
		"ctreated_date" timestamp with time zone,
		"updated_at" timestamp with time zone
	);`)
}

// Reverse the migrations
func (m *Cars_20231218_132732) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
