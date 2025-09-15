-- Create a 'snippets' table.
create table snippets (
  id integer not null primary key auto_increment,
  title varchar(100) not null,
  content text not null,
  created datetime not null,
  expires datetime not null
);

-- Add an index on the created column.
create index idx_snippets_created on snippets (created);

create table users (
  id integer not null primary key auto_increment,
  name varchar(255) not null,
  email varchar(255) not null,
  hashed_password char(60) not null,
  created datetime not null
);

alter table users add constraint users_uc_email unique (email);

insert into users (name, email, hashed_password, created) values (
  'Alice Jones',
  'alice@example.com',
  '$2a$10$.fCMCuuOpdPbBElL.yZbgeAHBgc77cBWI.yfMJyB2uHmIVnI8C79u',
  now()
);
