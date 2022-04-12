CREATE TABLE if not exists subscription_type
(
    id         SERIAL PRIMARY KEY,
    name       character varying(255)                             NOT NULL UNIQUE,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

INSERT INTO subscription_type (id, name)
VALUES (1, 'Limited'),
       (2, 'Unlimited');


CREATE TABLE if not exists subscription_action
(
    name        character varying(255) NOT NULL UNIQUE PRIMARY KEY ,
    description character varying(255) NOT NULL UNIQUE,
    unit        character varying(255) NOT NULL UNIQUE
);


INSERT INTO subscription_action (name, description, unit)
VALUES ('API Call', 'User interaction with public API Gateway', 'HTTP Requests');

-- create foreign key to subscription_type
CREATE TABLE if not exists subscription_account
(
    id                SERIAL PRIMARY KEY,
    account_holder_id int                    NOT NULL,
    name              character varying(255) NOT NULL UNIQUE,
    threshold         int                    NOT NULL,
    action            character varying(255) NOT NULL UNIQUE,
    product           varchar                NOT NULL,
    type              int                    NOT NULL UNIQUE,
    last_processed    timestamp with time zone,
    next_run          timestamp with time zone,
    CONSTRAINT action FOREIGN KEY (action) REFERENCES subscription_action (name),
    CONSTRAINT type FOREIGN KEY (type) REFERENCES subscription_type (id)
);


INSERT INTO subscription_account (id, account_holder_id, name, threshold, action, type, product, next_run)
VALUES (1, 123, 'Limited', 10, 'API Call', 1,  'Content', NOW() + INTERVAL '1 DAY');


-- Make composite key of
CREATE TABLE if not exists subscription_account_user_log
(
    GUID                    int,
    subscription_account_id int                                                NOT NULL,
    action_type             varchar                                            NOT NULL,
    interaction_at          timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (GUID, action_type, interaction_at)
--     CONSTRAINT account_holder_id FOREIGN KEY (subscription_account_id) REFERENCES subscription_account (account_holder_id)
);
