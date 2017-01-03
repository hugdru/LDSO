SET SCHEMA 'places4all';


INSERT INTO entity(id_country, name, email, username, password, image, banned_date, banned, reason, mobilephone, telephone, created_date) VALUES (
    (SELECT id FROM country where name='Portugal'),
        'Carlos Francisco', 'Cf3223@gmail.com', 'CarlosFcisco23', '16384$8$1$4cddfcdf22ffed7b834ffa53bb1fe257$c6d9547d6d7737845a31dfdc01488ba945d5086ba54ce0aa360e025b3d0423c9', NULL, NULL, NULL, NULL, NULL, '22 32319234', '2016-11-25T14:48:26.680513Z'), -- 2
     ((SELECT id FROM country where name='Portugal'),
         'António Coelho', 'AC33232@gmail.com', 'Antoenho23', '16384$8$1$4cddfcdf22ffed7b834ffa53bb1fe257$c6d9547d6d7737845a31dfdc01488ba945d5086ba54ce0aa360e025b3d0423c9', NULL, NULL, NULL, NULL, NULL, '22 3223134', '2016-11-25T14:47:26.680513Z'),  -- 3
     ((SELECT id FROM country where name='Portugal'),
         'Carlos Coelho', 'CC33231@gmail.com', 'CAntoenho23', '16384$8$1$4cddfcdf22ffed7b834ffa53bb1fe257$c6d9547d6d7737845a31dfdc01488ba945d5086ba54ce0aa360e025b3d0423c9', NULL, NULL, NULL, NULL, NULL, '22 3223134', '2016-11-25T14:47:26.680513Z'); -- 4
--insert auditor
INSERT INTO auditor(id_entity) VALUES
  (2), -- Carlos Francisco
  (3), -- António Coelho
  (4); -- Carlos Coelho

INSERT INTO template (name,description,created_date)
VALUES ('Modelo Teste','template teste 1','2016-12-25T14:47:26.680513Z'); -- 3


INSERT INTO audit(id_property, id_auditor, id_template,rating,observation,created_date ) VALUES
    (1,2,1,0,'audit 1','2016-12-25T14:47:26.680513Z'),
    (1,3,1,0,'audit 1','2016-12-26T14:47:26.680513Z');
