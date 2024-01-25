CREATE TABLE people_info (
    id bigserial not null primary key,
    name varchar not null,
    surname varchar not null,
    pathronymic varchar,
    age int not null,
    sex varchar not null,
    nation varchar not null
);
CREATE USER dev WITH PASSWORD '12345';
GRANT ALL ON people_info TO dev;
GRANT USAGE, SELECT ON SEQUENCE people_info_id_seq TO dev;
