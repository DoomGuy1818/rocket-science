-- +goose Up
create table if not exists orders (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL,
    total_price NUMERIC(12, 2),
    transaction_id UUID,
    payment_method VARCHAR(30) NOT NULL,
    status VARCHAR(30) NOT NULL
);

create table if not exists order_items(
    id UUID PRIMARY KEY NOT NULL,
    order_id UUID references orders(id) ON DELETE CASCADE,
    part_id UUID NOT NULL,
    unit_price NUMERIC(12, 2)
);

-- +goose Down
drop table if exists orders;

drop table if exists order_items;
