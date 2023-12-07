create table keys
(
    id         varchar(36) not null primary key,
    account_id varchar(36) foreign key references accounts(id),
    name       varchar(200),
    type       varchar(100),
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp
)
go

