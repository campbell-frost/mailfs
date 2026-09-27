create table if not exists files (
    id text primary key,
    filename text not null,
    temp_path text not null,
    size integer not null,
    status text not null,
    created_at text not null
);

create table if not exists chunks (
    file_id text not null references files(id),
    seq integer not null,
    size integer not null,
    sha256 text not null,
    ref text not null,
    primary key (file_id, seq)
);

create table if not exists users (
    id integer primary key,
    sub text not null unique,
    email text not null,
    name text not null,
    created_at text not null
);

create table if not exists sessions (
    token_hash blob primary key,
    user_id integer not null references users(id) on delete cascade,
    created_at text not null,
    expires_at text not null
);
