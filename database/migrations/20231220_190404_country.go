package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Country_20231220_190404 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Country_20231220_190404{}
	m.Created = "20231220_190404"

	migration.Register("Country_20231220_190404", m)
}

// Run the migrations
func (m *Country_20231220_190404) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE "mod_country_master" (
		"country_id" bigint,
		"name" text,
		"iso3" text,
		"iso2" text,
		"numeric_code" bigint,
		"phone_code" text,
		"capital" text,
		"currency" text,
		"currency_name" text,
		"currency_symbol" text,
		"tld" text,
		"native" text,
		"region" text,
		"region_id" text,
		"subregion" text,
		"subregion_id" text,
		"nationality" text,
		"timezones" text,
		"latitude" text,
		"longitude" text,
		"emoji" text,
		"emojiU" text
	);  `)

}

// Reverse the migrations
func (m *Country_20231220_190404) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
