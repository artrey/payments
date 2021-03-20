-- табличка с пользователями и их паролями
CREATE TABLE users
(
    id       BIGSERIAL PRIMARY KEY,
    login    VARCHAR(40) NOT NULL UNIQUE,
    password VARCHAR     NOT NULL,
    roles    VARCHAR[] NOT NULL DEFAULT '{}',
    created  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- табличка с токенами (чтобы пользователь каждый раз не присылал пароль и логин)
CREATE TABLE tokens
(
    id      VARCHAR PRIMARY KEY,
    userId  BIGINT    NOT NULL REFERENCES users (id),
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- табличка с платежами
CREATE TABLE payments
(
    id       VARCHAR PRIMARY KEY,
    senderId BIGINT NOT NULL REFERENCES users (id),
    amount   BIGINT NOT NULL
);

