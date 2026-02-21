package payindb

var (
	queryByRUT = `
	SELECT
		"RUT", balance
	FROM 
		balance
	WHERE
		"RUT" = :rut
	`

	queryUpdateByRUT = `
	INSERT INTO balance (rut, balance)
	VALUES (:rut, :balance)
	ON CONFLICT (rut)
	DO UPDATE SET balance = balance.balance + EXCLUDED.balance;
	
	`

	queryAudit = `
	INSERT INTO WEBHOOK ("ID_MERCADO_PAGO",topic, state)
	VALUES (:id_mercado_pago, :topic, :state);`
)
