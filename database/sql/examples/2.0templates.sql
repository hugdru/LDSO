SET SCHEMA 'places4all';

INSERT INTO template(name, created_date) VALUES
  ('Modelo 1', '2016-11-25T14:46:26.680513Z'), -- 1
  ('Modelo 2', '2016-11-25T14:46:26.680513Z'); -- 2

INSERT INTO maingroup(id_template, name, weight, created_date) VALUES
  (1, 'Grupo 1 (M1)', 40, '2016-11-25T14:46:26.680513Z'), -- 1
  (1, 'Grupo 2 (M1)', 60, '2016-11-25T14:46:26.680513Z'), -- 2
  (2, 'Grupo 3 (M2)', 100, '2016-11-25T14:46:26.680513Z'); -- 3

INSERT INTO subgroup(id_maingroup, name, weight, created_date) VALUES
  (1, 'Subgrupo 1 (G1)', 25, '2016-11-25T14:46:26.680513Z'), -- 1
  (1, 'Subgrupo 2 (G1)', 42, '2016-11-25T14:46:26.680513Z'), -- 2
  (1, 'Subgrupo 3 (G1)', 33, '2016-11-25T14:46:26.680513Z'), -- 3
  (2, 'Subgrupo 4 (G2)', 40, '2016-11-25T14:46:26.680513Z'), -- 4
  (2, 'Subgrupo 5 (G2)', 60, '2016-11-25T14:46:26.680513Z'), -- 5
  (3, 'Subgrupo 6 (G3)', 100, '2016-11-25T14:46:26.680513Z'); -- 6

INSERT INTO legislation(name, description, url) VALUES
  ('lei 1', NULL, NULL), -- 1
  ('lei 2', NULL, NULL), -- 2
  ('lei 3', NULL, NULL); -- 3

INSERT INTO criterion(id_subgroup, id_legislation, name, weight, created_date) VALUES
  (1, 1, 'Critério 1 (S1)', 11, '2016-11-25T14:46:26.680513Z'), -- 1
  (1, NULL, 'Critério 2 (S1)', 11, '2016-11-25T14:46:26.680513Z'), -- 2
  (1, 1, 'Critério 3 (S1)', 78, '2016-11-25T14:46:26.680513Z'), -- 3
  (2, NULL, 'Critério 4 (S2)', 33, '2016-11-25T14:46:26.680513Z'), -- 4
  (2, 1, 'Critério 5 (S2)', 7, '2016-11-25T14:46:26.680513Z'), -- 5
  (2, NULL, 'Critério 6 (S2)', 25, '2016-11-25T14:46:26.680513Z'), --6
  (2, 2, 'Critério 7 (S2)', 15, '2016-11-25T14:46:26.680513Z'), -- 7
  (2, NULL, 'Critério 8 (S2)', 1, '2016-11-25T14:46:26.680513Z'), -- 8
  (2, 2, 'Critério 9 (S2)', 19, '2016-11-25T14:46:26.680513Z'), -- 9
  (3, NULL, 'Critério 10 (S3)', 64, '2016-11-25T14:46:26.680513Z'), -- 10
  (3, 3, 'Critério 11 (S3)', 36, '2016-11-25T14:46:26.680513Z'), -- 11
  (4, NULL, 'Critério 12 (S4)', 100, '2016-11-25T14:46:26.680513Z'), -- 12
  (5, 3, 'Critério 13 (S5)', 9, '2016-11-25T14:46:26.680513Z'), -- 13
  (5, NULL, 'Critério 14 (S5)', 9, '2016-11-25T14:46:26.680513Z'), -- 14
  (5, 3, 'Critério 15 (S5)', 82, '2016-11-25T14:46:26.680513Z'), -- 15
  (6, NULL, 'Critério 16 (S6)', 11, '2016-11-25T14:46:26.680513Z'), -- 16
  (6, 1, 'Critério 17 (S6)', 11, '2016-11-25T14:46:26.680513Z'), -- 17
  (6, 2, 'Critério 18 (S6)', 11, '2016-11-25T14:46:26.680513Z'), -- 18
  (6, 3, 'Critério 19 (S6)', 11, '2016-11-25T14:46:26.680513Z'), -- 19
  (6, NULL, 'Critério 20 (S6)', 56, '2016-11-25T14:46:26.680513Z'); -- 20

INSERT INTO accessibility(name) VALUES
  ('Auditiva'), -- 1
  ('Cognitiva'), -- 2
  ('Física'), -- 3
  ('Visual'); -- 4

INSERT INTO criterion_accessibility(id_criterion, id_accessibility, weight) VALUES
  (1, 1, 20),
  (1, 2, 30),
  (1, 3, 30),
  (1, 4, 20),
  (2, 1, 15),
  (2, 2, 20),
  (2, 3, 55),
  (2, 4, 10),
  (3, 1, 11),
  (3, 2, 11),
  (3, 3, 11),
  (3, 4, 67),
  (4, 1, 20),
  (4, 2, 20),
  (4, 3, 60),
  (4, 4, 0),
  (5, 1, 13),
  (5, 2, 32),
  (5, 3, 55),
  (5, 4, 0),
  (6, 1, 33),
  (6, 2, 33),
  (6, 3, 34),
  (6, 4, 0),
  (7, 1, 55),
  (7, 2, 40),
  (7, 3, 5),
  (7, 4, 0),
  (8, 1, 98),
  (8, 2, 1),
  (8, 3, 1),
  (8, 4, 0),
  (9, 1, 77),
  (9, 2, 13),
  (9, 3, 10),
  (9, 4, 0),
  (10, 1, 19),
  (10, 2, 75),
  (10, 3, 6),
  (10, 4, 0),
  (11, 1, 88),
  (11, 2, 4),
  (11, 3, 8),
  (11, 4, 0),
  (12, 1, 66),
  (12, 2, 24),
  (12, 3, 10),
  (12, 4, 0),
  (13, 1, 17),
  (13, 2, 17),
  (13, 3, 66),
  (13, 4, 0),
  (14, 1, 55),
  (14, 2, 33),
  (14, 3, 12),
  (14, 4, 0),
  (15, 1, 54),
  (15, 2, 6),
  (15, 3, 40),
  (15, 4, 0),
  (16, 1, 9),
  (16, 2, 8),
  (16, 3, 83),
  (16, 4, 0),
  (17, 1, 15),
  (17, 2, 15),
  (17, 3, 70),
  (17, 4, 0),
  (18, 1, 76),
  (18, 2, 14),
  (18, 3, 10),
  (18, 4, 0),
  (19, 1, 32),
  (19, 2, 23),
  (19, 3, 45),
  (19, 4, 0),
  (20, 1, 20),
  (20, 2, 30),
  (20, 3, 30),
  (20, 4, 20);

UPDATE template
SET closed_date = '2017-01-03T22:30:53.083371Z', closed = TRUE
WHERE id IN (1, 2);
