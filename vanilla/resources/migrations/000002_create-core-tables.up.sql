create table if not exists core_customers
(
    id         int primary key,
    name       varchar(100) not null,
    email      varchar(100) not null,
    phone      varchar(100) not null
);


insert into core_customers (id, name, email, phone) values (1, 'Arisha Barron', 'a.barron@someemail.test', '1234567890');
insert into core_customers (id, name, email, phone) values (2, 'Branden Gibson', 'b.gibson@someemail.test', '1234567890');
insert into core_customers (id, name, email, phone) values (3, 'Rhonda Church', 'r.church@someemail.test', '1234567890');
insert into core_customers (id, name, email, phone) values (4, 'Georgina Haze', 'g.Haze@someemail.test', '1234567890');


create table if not exists core_accounts
(
    id         int primary key,
    owner      int not null,
    balance    decimal(10, 2) default 0.00 not null
);

create table if not exists core_entries
(
    id         int primary key,
    account_id int not null,
    amount     decimal(10, 2) default 0.00 not null
);

create table if not exists core_transfers
(
    id         int primary key,
    from_account_id int not null,
    to_account_id int not null,
    amount     decimal(10, 2) default 0.00 not null
);
