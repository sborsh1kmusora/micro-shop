-- +goose Up
create table orders (
    uuid text not null,
    user_uuid text not null,
    items_uuids text[] not null,
    transaction_uuid text not null,
    total_price float not null,
    payment_method text not null,
    status text not null
);

-- +goose Down
drop table orders;