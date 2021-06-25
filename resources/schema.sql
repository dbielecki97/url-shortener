create table urls
(
    code       varchar not null
        constraint urls_pk
            primary key,
    url        varchar not null,
    created_at timestamp default CURRENT_TIMESTAMP
);

alter table urls
    owner to postresuser;

create unique index urls_code_uindex
    on urls (code);

