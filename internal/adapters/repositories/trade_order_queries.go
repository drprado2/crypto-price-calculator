package repositories

const (
	insertTradeOrderQuery = `
INSERT INTO trade_order
    (trade_order_id, trade_id, maker_order_id, taker_order_id, side, size, price, product_id, sequence, time) VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	getLastNTradesByProductQuery = `
SELECT size, price from trade_order where product_id = $1 order by sequence desc limit $2
`
)
