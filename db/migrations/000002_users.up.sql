create table if not exists users(
	id uuid default uuid_generate_v4() primary key,
	name text check(length(name) >= 3 and length(name) <= 25) not null,
	surname text check(length(surname) >= 3 and length(surname) <= 25) not null,
	patronymic text check(length(patronymic) >= 3 and length(patronymic) <= 25),
	age integer check(age >= 1 and age <= 100) not null,
	gender varchar(7) check(gender in ('male', 'female')) not null,
	countries jsonb[] check(array_length(countries, 1) > 0 and array_length(countries, 1) = cardinality(countries)) not null,
	user_created_at timestamp with time zone default current_timestamp 
);