SET SCHEMA 'places4all';

INSERT INTO address(id_country, address_line1, address_line2, address_line3, town_city, county, postcode, latitude, longitude) VALUES
  ((SELECT id FROM country WHERE name='Portugal'), 'Via Futebol Clube do Porto', NULL, NULL, 'Porto', NULL, '4350-415 Porto', 41.1606111, -8.5820467); -- 1

INSERT INTO property(id_address, name, details, created_date) VALUES
  (1, 'Estádio do Dragão', 'The Estádio do Dragão is an all-seater football stadium located in Porto, Portugal, and the current home ground of Futebol Clube do Porto. It has a capacity of 52,000, making it the second largest football ground in Portugal.', '2016-11-25T14:46:26.672053Z'); -- 1

INSERT INTO tag(name) VALUES
  ('futebol'), -- 1
  ('super dragões'), -- 2
  ('estádio'); -- 3

INSERT INTO property_tag(id_property, id_tag) VALUES
  (1, 1),
  (1, 2),
  (1, 3);

INSERT INTO entity(id_country, name, email, username, password, image, banned_date, banned, reason, mobilephone, telephone, created_date) VALUES
  ((SELECT id FROM country where name='Portugal'), 'FC Porto', 'fcporto@fcporto.pt', 'fcporto', '16384$8$1$4cddfcdf22ffed7b834ffa53bb1fe257$c6d9547d6d7737845a31dfdc01488ba945d5086ba54ce0aa360e025b3d0423c9', NULL, NULL, NULL, NULL, NULL, '22 557 0400', '2016-11-25T14:46:26.680513Z'); -- 1

INSERT INTO entity (id_country, name, email, username, password, image, banned_date, banned, reason, mobilephone, telephone, created_date)
VALUES
  ((SELECT id
    FROM country
    WHERE name = 'Portugal'), 'Super Administrador',
                              'superadmin@places4all.com',
                              'superadmin',
                              '16384$8$1$4cddfcdf22ffed7b834ffa53bb1fe257$c6d9547d6d7737845a31dfdc01488ba945d5086ba54ce0aa360e025b3d0423c9',
                              NULL, NULL, NULL, NULL, NULL, NULL,
   '2016-11-25T14:46:26.680513Z'); -- 2

INSERT INTO entity (id_country, name, email, username, password, image, banned_date, banned, reason, mobilephone, telephone, created_date)
VALUES
  ((SELECT id
    FROM country
    WHERE name = 'Portugal'), 'Local Administrador',
                              'localadmin@places4all.com',
                              'localadmin',
                              '16384$8$1$4cddfcdf22ffed7b834ffa53bb1fe257$c6d9547d6d7737845a31dfdc01488ba945d5086ba54ce0aa360e025b3d0423c9',
                              NULL, NULL, NULL, NULL, NULL, NULL,
   '2016-11-25T14:46:26.680513Z'); -- 2

INSERT INTO entity (id_country, name, email, username, password, image, banned_date, banned, reason, mobilephone, telephone, created_date)
VALUES
  ((SELECT id
    FROM country
    WHERE name = 'Portugal'), 'Auditor',
                              'auditor@places4all.com',
                              'auditor',
                              '16384$8$1$4cddfcdf22ffed7b834ffa53bb1fe257$c6d9547d6d7737845a31dfdc01488ba945d5086ba54ce0aa360e025b3d0423c9',
                              NULL, NULL, NULL, NULL, NULL, NULL,
   '2016-11-25T14:46:26.680513Z'); -- 2

INSERT INTO client(id_entity) VALUES
  (1);

INSERT INTO superadmin(id_entity) VALUES
  (2);

INSERT INTO localadmin(id_entity) VALUES
  (3);

INSERT INTO auditor(id_entity) VALUES
  (4);


INSERT INTO property_client(id_property, id_client) VALUES
  (1, 1);
