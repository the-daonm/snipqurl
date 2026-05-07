CREATE TABLE urls (
  id bigserial primary key,
  original_url text not null,
  short_code varchar(10) unique,
  clicks integer default 0,
  created_at timestamp default now()
);
