package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type HomeSettings_20231218_134042 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &HomeSettings_20231218_134042{}
	m.Created = "20231218_134042"

	migration.Register("HomeSettings_20231218_134042", m)
}

// Run the migrations
func (m *HomeSettings_20231218_134042) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE IF NOT EXISTS "home_pages_setting_table" (
		"page_setting_id" serial NOT NULL PRIMARY KEY,
		"section" text NOT NULL DEFAULT '' ,
		"data_type" varchar(255) NOT NULL DEFAULT '' ,
		"unique_code" text NOT NULL DEFAULT '' ,
		"setting_data" text NOT NULL,
		"created_date" timestamp with time zone NOT NULL,
		"updated_date" timestamp with time zone,
		"created_by" integer NOT NULL DEFAULT 0 ,
		"updated_by" integer NOT NULL DEFAULT 0
	);`)

}

// Reverse the migrations
func (m *HomeSettings_20231218_134042) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
