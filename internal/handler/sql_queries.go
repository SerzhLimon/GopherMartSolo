package handler

const (
	queryInsertWithReturning = `
		SELECT user_id, current, withdrawn 
		FROM user_balance 
		WHERE user_id = $1
	`

	queryUpdateUserBalance = `
		UPDATE user_balance
		SET current = $2,
			withdrawn = $3,
			updated_at = $4
		WHERE user_id = $1
	`

	queryInsertBalanceOperations = `
		INSERT INTO balance_operations (user_id, type, amount, order_number, processed_at)
		VALUES ($1, 'withdraw', $2, $3, $4)
	`
	querySelectUserID = `SELECT user_id FROM orders WHERE number = $1`

	queryCreateNewOrder = `INSERT INTO orders (number, user_id, status, uploaded_at) VALUES ($1, $2, $3, $4)`

	queryCreateUser = `INSERT INTO users (login, password) VALUES ($1, $2)`

	queryIsRegistred = `SELECT COUNT(login) FROM users WHERE login = $1`

	queryAddBalance = `INSERT INTO user_balance (user_id, current, withdrawn, updated_at, uploaded_at) 
	VALUES ($1, $2, $3, $4, $5)`

	queryLogin = `SELECT id FROM users WHERE login = $1 AND password = $2`

	queryGetWithdrawal = `SELECT * FROM balance_operations WHERE user_id = $1 AND type = 'withdraw'  ORDER BY processed_at DESC`

	queryGetOrder = `SELECT number, status, accrual, uploaded_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC`

	queryGetBalance = `SELECT user_id, current, withdrawn FROM user_balance WHERE user_id = $1`
)
