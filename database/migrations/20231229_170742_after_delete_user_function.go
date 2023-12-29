package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type AfterDeleteUserFunction_20231229_170742 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AfterDeleteUserFunction_20231229_170742{}
	m.Created = "20231229_170742"

	migration.Register("AfterDeleteUserFunction_20231229_170742", m)
}

// Run the migrations
func (m *AfterDeleteUserFunction_20231229_170742) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE OR REPLACE FUNCTION after_delete_user()
RETURNS TRIGGER AS $$
BEGIN
     INSERT INTO backup_users(user_id,first_name,last_name,email,phone_number,password,isverified,otp_code,role,country_id,created_date,updated_date,delete_from_user)
     VALUES(OLD.user_id,OLD.first_name,OLD.last_name,OLD.email,OLD.phone_number,OLD.password,OLD.isverified,OLD.otp_code,OLD.role,OLD.country_id,OLD.created_date,OLD.updated_date,CURRENT_TIMESTAMP);
     RETURN NULL;
END;
$$ LANGUAGE plpgsql;`)
}

// Reverse the migrations
func (m *AfterDeleteUserFunction_20231229_170742) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
