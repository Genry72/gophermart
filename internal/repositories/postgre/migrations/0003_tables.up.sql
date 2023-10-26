create table user_balance
(
    user_id         bigserial               not null
        constraint user_balance_users_user_id_fk
            references users,
    accruals        numeric   default 0     not null,
    drawal          numeric   default 0     not null,
    current_balance numeric   default 0     not null,
    last_update     timestamp default now() not null
);

comment on table user_balance is 'Актуальный баланс пользователей';
comment on column user_balance.user_id is 'Уникальный идентификатор пользователя';
comment on column user_balance.accruals is 'Сумма начисленных баллов';
comment on column user_balance.drawal is 'Сумма списанных баллов';
comment on column user_balance.current_balance is 'Актуальный баланс пользователя';
comment on column user_balance.last_update is 'Дата последнего обновления';

