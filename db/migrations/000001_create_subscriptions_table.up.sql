-- Used with golang-migrate. This script will run on start-up locally to do database migrations

CREATE TABLE if not exists  subscription (
                                             id SERIAL PRIMARY KEY,
                                             name character varying(255) NOT NULL UNIQUE,
                                             created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                             updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


-- INSERT INTO subscription (id, name) VALUES
--                                         (1, 'Pay as you go'),
--                                         (2, 'Unlimited'),
--                                         (3, 'Pay Monthly');