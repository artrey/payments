INSERT INTO users(id, login, password, roles)
VALUES (1, 'admin', '$2a$10$ctPFhgJh.YIE21AA0OGl5er3p9f3XsAwkmTXnP2I7BxCpQbr1QAg2', '{"ADMIN", "USER"}'),
       (2, 'user', '$2a$10$ctPFhgJh.YIE21AA0OGl5er3p9f3XsAwkmTXnP2I7BxCpQbr1QAg2', '{"USER"}');

INSERT INTO payments(id, senderId, amount)
VALUES ('9f3ff07d-4d4e-4825-8a92-526df6dc7545', 1, 10000),
       ('c763b586-f7d2-4f38-a2fe-70b55ebc654a', 2, 50000),
       ('eb476191-0ad5-48f4-8c31-aad1eeacbb9f', 1, 100000),
       ('b11861e8-e566-417e-8b4c-3155997e1c61', 1, 50000),
       ('9bc02e9f-12fd-4e19-bcc9-14bc98f79bfd', 2, 50000),
       ('786e2125-12cd-4680-9dde-bce77f67efb8', 2, 100000),
       ('1e3cf754-c734-4519-968c-fbba7a409410', 2, 500000);
