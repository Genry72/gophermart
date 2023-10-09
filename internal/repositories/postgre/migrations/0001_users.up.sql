create table if not exists users
(
    user_id       bigserial primary key   not null,
    username      varchar                 not null,
    password_hash varchar                 not null,
    email         varchar,
    first_name    varchar,
    last_name     varchar,
    phone         varchar,
    created_at    timestamp default now() not null,
    updated_at    timestamp default now() not null,
    deleted_at    timestamp
);

comment on table users is 'Пользователи Gopher Mart';
comment on column users.user_id is 'Уникальный идентификатор пользователя';
comment on column users.username is 'Логин пользователя';
comment on column users.password_hash is 'Пароль пользователя';
comment on column users.email is 'Электронная почта пользователя';
comment on column users.first_name is 'Имя';
comment on column users.last_name is 'Фамилия';
comment on column users.phone is 'Номер телефона';
comment on column users.created_at is 'Дата создания';
comment on column users.updated_at is 'Дата обновления';
comment on column users.deleted_at is 'Дата удаления';