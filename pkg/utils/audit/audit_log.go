package helpers

import (
	"fmt"
	"os"
	"time"

	custom_log "admin-phone-shop-api/pkg/custom_log"
	sql "admin-phone-shop-api/pkg/sql"

	"github.com/jmoiron/sqlx"
)

func AddMemeberAuditLog(admin_id float64, audit_context string, audit_desc string, audit_type_id int, admin_agent string, admin_name string, ip string, by_id float64, db_pool *sqlx.DB) (*bool, error) {

	orderSeqName := "tbl_admin_audit_id_seq"
	orderVal, err := sql.GetSeqNextVal(orderSeqName, db_pool)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch next order value: %w", err)
	}

	var query = `
		INSERT INTO tbl_admin_audit (
			id, 
			admin_id, 
			admin_audit_context, 
			admin_audit_desc, 
			audit_type_id, 
			admin_agent, 
			operator, 
			ip, 
			status_id, 
			"order", 
			created_by, 
			created_at
		) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)`
	app_timezone := os.Getenv("TIME_ZONE")
	location, err := time.LoadLocation(app_timezone)
	if err != nil {
		return nil, fmt.Errorf("failed to load location: %w", err)
	}
	local_now := time.Now().In(location)
	_, err = db_pool.Exec(
		query,
		*orderVal,
		admin_id,
		audit_context,
		audit_desc,
		audit_type_id,
		admin_agent,
		admin_name,
		ip,
		1,
		*orderVal,
		by_id,
		local_now,
	)
	if err != nil {
		custom_log.NewCustomLog("admin_create_failed", err.Error(), "error")
		return nil, err
	}
	state := true
	return &state, nil
}