SET SCHEMA 'places4all';


INSERT INTO entity(id_country, name, email, username, password, image, banned_date, banned, reason, mobilephone, telephone, created_date) VALUES (
    (SELECT id FROM country where name='Portugal'),
        'Carlos fransisco', 'Cf3223@gmail.com', 'CarlosFcisco23', '111111', NULL, NULL, NULL, NULL, NULL, '22 32319234', '2016-11-25T14:48:26.680513Z'), -- 1
     ((SELECT id FROM country where name='Portugal'),
         'Antonio Coelho', 'AC33232@gmail.com', 'Antoenho23', '222222', NULL, NULL, NULL, NULL, NULL, '22 3223134', '2016-11-25T14:47:26.680513Z'),  -- 2
     ((SELECT id FROM country where name='Portugal'),
         'Carlos Coelho', 'CC33231@gmail.com', 'CAntoenho23', '222222', NULL, NULL, NULL, NULL, NULL, '22 3223134', '2016-11-25T14:47:26.680513Z');
--insert auditor
INSERT INTO auditor(id_entity) VALUES
  (2), --Carlos Francisco
  (3), -- antonioa Coelho
  (4); -- carlos Coelho

INSERT INTO template (name,description,created_date) VALUES ('Modelo Teste',
                                                             'template teste 1','2016-12-25T14:47:26.680513Z');


INSERT INTO audit(id_property, id_auditor, id_template,rating,observation,created_date ) VALUES
    (1,1,1,2,'audit 1','2016-12-25T14:47:26.680513Z'),
    (1,2,1,2,'audit 1','2016-12-26T14:47:26.680513Z');