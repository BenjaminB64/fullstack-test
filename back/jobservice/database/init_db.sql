create schema if not exists public;

create table if not exists jobs (
  id serial primary key,
  name varchar(255) not null,
  status varchar(50) not null check (status in ('pending', 'in_progress', 'completed', 'failed')) default 'pending',
  task_type varchar(50) not null check (task_type in ('get_chaban_delmas_bridge_status', 'get_weather')),

  created_at timestamp without time zone not null default current_timestamp,
  updated_at timestamp without time zone null default null,
  deleted_at timestamp without time zone null default null
);

create table if not exists weather_job_results (
  id serial primary key,
  job_id integer not null references jobs(id),

  latitude numeric(9,6) not null,
  longitude numeric(9,6) not null,

  temperature numeric(5,2) not null,
  relative_humidity numeric(5,2) not null,
  weather_wmo_code int not null,

  created_at timestamp without time zone not null default current_timestamp,
  updated_at timestamp without time zone null default null,
  deleted_at timestamp without time zone null default null
);

create table if not exists chaban_delmas_bridge_job_results (
  id serial primary key,
  job_id integer not null references jobs(id),

  close_time timestamp without time zone not null,
  open_time timestamp without time zone not null,

  created_at timestamp without time zone not null default current_timestamp,
  updated_at timestamp without time zone null default null,
  deleted_at timestamp without time zone null default null
);


