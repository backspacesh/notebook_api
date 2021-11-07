CREATE TABLE articles(
    id serial primary key,
    article_header varchar(50) not null,
    article_text text not null,
    author_id integer not null references users(id),
    creating_date date not null
);