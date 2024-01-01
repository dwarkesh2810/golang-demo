package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type AfterDeleteUserTrigger_20231229_170814 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AfterDeleteUserTrigger_20231229_170814{}
	m.Created = "20231229_170814"

	migration.Register("AfterDeleteUserTrigger_20231229_170814", m)
}

// Run the migrations
func (m *AfterDeleteUserTrigger_20231229_170814) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update

	m.SQL(`CREATE OR REPLACE TRIGGER after_delete_users_trigger
AFTER DELETE ON users
FOR EACH ROW
EXECUTE FUNCTION after_delete_user();`)
}

// Reverse the migrations
func (m *AfterDeleteUserTrigger_20231229_170814) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
