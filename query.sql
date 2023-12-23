create type todostatus as enum (
  'TODO', 
  'INPROGRESS', 
  'DONE'
);

drop table todo;

create table todo (
  id varchar(255) primary key,
  title varchar(255) not null,
  description varchar(255) not null,
  status todostatus not null,
  createdDate timestamp without time zone not null default (now() at time zone 'utc'),
  updatedDate timestamp without time zone,
  dueDate timestamp without time zone
);

insert into todo (id, title, description, status)
values ('x-dummy-todo-id', 'x-dummy-todo-title', 'x-dummy-todo-description', 'TODO');

select * from todo;