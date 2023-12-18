package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type LanguageLableLang_20231218_134204 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &LanguageLableLang_20231218_134204{}
	m.Created = "20231218_134204"

	migration.Register("LanguageLableLang_20231218_134204", m)
}

// Run the migrations
func (m *LanguageLableLang_20231218_134204) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`
	CREATE TABLE IF NOT EXISTS "language_lable_lang" (
		"lang_id" serial NOT NULL PRIMARY KEY,
		"language_code" text NOT NULL DEFAULT '' ,
		"language_value" text NOT NULL DEFAULT '' ,
		"lable_code" text NOT NULL DEFAULT ''  UNIQUE,
		"section" text NOT NULL DEFAULT '' 
	);`)

}

// Reverse the migrations
func (m *LanguageLableLang_20231218_134204) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
