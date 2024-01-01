package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type DeleteAuditLogs_20240101_121028 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &DeleteAuditLogs_20240101_121028{}
	m.Created = "20240101_121028"

	migration.Register("DeleteAuditLogs_20240101_121028", m)
}

// Run the migrations
func (m *DeleteAuditLogs_20240101_121028) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE OR REPLACE PROCEDURE delete_old_records_procedure()
	LANGUAGE plpgsql
	AS $$
	DECLARE
		record_date TIMESTAMP;
	BEGIN
		FOR record_date IN (SELECT created_date FROM audit_logs) LOOP
			IF EXTRACT(day FROM (NOW() - record_date))::int > 45 THEN
				DELETE FROM audit_logs WHERE created_date = record_date;
			END IF;
		END LOOP;
	END;
	$$;`)

}

// Reverse the migrations
func (m *DeleteAuditLogs_20240101_121028) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
