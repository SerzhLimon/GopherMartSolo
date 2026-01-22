package accrual

const (
	querySelectOrders = `SELECT * FROM orders WHERE status IN ('NEW', 'PROCESSING')`

	queryUpdateFailed = `UPDATE orders SET status = $1, accrual = $2, processed_at = $4 WHERE number = $3`

	queryUpdate = `UPDATE orders SET status = $1, processed_at = $3 WHERE number = $2`

	queryUpdateBalance = `
		UPDATE user_balance 
		SET current = current + $2, 
			updated_at = $3 
		WHERE user_id = $1
	`

	queryLogOperation = `INSERT INTO balance_operations (user_id, amount, "type", order_number, processed_at) VALUES ($1, $2, $3, $4, $5)`
)
