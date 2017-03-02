CREATE TABLE t_news (
    id int not null AUTO_INCREMENT,
    name varchar(100),
    url varchar(100),
    inserted_date datetime,
    primary key (id)
);