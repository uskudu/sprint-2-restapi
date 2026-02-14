drop table if exists tasks;

create table tasks (
    id int generated always as identity primary key,
    title varchar(255) not null,
    description text,
    completed boolean default false,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);

insert into tasks (title, description, completed) values
    ('learn go', 'finish base course', false),
    ('find God', 'come back after you find God', true),
    ('finish this project', 'today', false);