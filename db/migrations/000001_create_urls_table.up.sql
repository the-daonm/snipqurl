CREATE TABLE urls (
  id bigserial primary key,
  original_url text not null,
  short_code varchar(10) unique,
  created_at timestamp default now(),
  clicked integer default 0
);
