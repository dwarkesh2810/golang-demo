package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type LanguageLable_20231218_134228 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &LanguageLable_20231218_134228{}
	m.Created = "20231218_134228"

	migration.Register("LanguageLable_20231218_134228", m)
}

// Run the migrations
func (m *LanguageLable_20231218_134228) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE IF NOT EXISTS "language_lable" (
		"lable_id" serial NOT NULL PRIMARY KEY,
		"lable_code" text NOT NULL DEFAULT '' ,
		"language_value" text NOT NULL DEFAULT '' ,
		"language_code" text NOT NULL DEFAULT '' ,
		"lang_id" integer NOT NULL DEFAULT 0 ,
		"section" text NOT NULL DEFAULT '' 
	);`)
}

// Reverse the migrations
func (m *LanguageLable_20231218_134228) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
