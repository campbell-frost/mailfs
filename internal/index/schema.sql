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
