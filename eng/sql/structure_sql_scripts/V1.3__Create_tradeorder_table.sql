CREATE TABLE trade_order
(
    trade_order_id uuid PRIMARY KEY NOT NULL,
    trade_id       bigint           NOT NULL,
    maker_order_id uuid             NOT NULL,
    taker_order_id uuid             NOT NULL,
    side           VARCHAR(8)       NOT NULL,
    size           decimal          NOT NULL,
    price          decimal          NOT NULL,
    product_id     VARCHAR(12)      NOT NULL,
    sequence       bigint           NOT NULL,
    time           TIME             NOT NULL,
    UNIQUE (trade_id, product_id)
);

CREATE INDEX idx_product_id_sequence ON trade_order (product_id, sequence DESC);