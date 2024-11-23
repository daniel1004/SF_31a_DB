DROP TABLE IF EXISTS posts,authors;
CREATE TABLE authors (
                         id BIGSERIAL PRIMARY KEY,
                         name TEXT
);
CREATE TABLE posts(
                      id BIGSERIAL PRIMARY KEY,
                      author_id BIGINT REFERENCES authors(id),
                      title TEXT NOT NULL,
                      context TEXT,
                      created_at bigint
);

