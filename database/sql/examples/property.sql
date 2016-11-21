SET SCHEMA 'places4all';

INSERT INTO address(id, id_country, address_line1, address_line2, address_line3, town_city, county, postcode, latitude, longitude) VALUES
  (1, (SELECT id FROM country WHERE name='Portugal'), 'Via Futebol Clube do Porto', NULL, NULL, 'Porto', NULL, '4350-415 Porto', 41.1606111, -8.5820467);

INSERT INTO property(id, id_address, name, details, created) VALUES
  (1, 1, 'Estádio do Dragão', 'The Estádio do Dragão is an all-seater football stadium located in Porto, Portugal, and the current home ground of Futebol Clube do Porto. It has a capacity of 52,000, making it the second largest football ground in Portugal.', current_timestamp);

INSERT INTO tag(id, name) VALUES
  (1, 'futebol'),
  (2, 'super dragões'),
  (3, 'estádio');

INSERT INTO property_tag(id_property, id_tag) VALUES
  (1, 1),
  (1, 2),
  (1, 3);

INSERT INTO gallery(id, id_property, name, description, created) VALUES
  (1, 1, 'Estádio', 'Where we battle', current_timestamp);

INSERT INTO image(id, id_gallery, name, description, image_url, created) VALUES
  (1, 1, 'The inside', NULL, 'http://img11.deviantart.net/696b/i/2011/013/2/e/estadio_do_dragao_by_jopeg-d373f2n.jpg', current_timestamp),
  (2, 1, 'The outside', 'Dragonlair', 'https://3.bp.blogspot.com/-ZbjmUcASPhY/V3ZC1CySEGI/AAAAAAABJMQ/u8ADqFqPJRETL3crIYv6mk-1YGD5lKvAACLcB/s1600/Estadio-do-Dragao-Porto.jpg', current_timestamp);

INSERT INTO entity(id, id_country, name, email, username, password, image_url, banned, reason, mobilephone, telephone, created) VALUES
  (1, (SELECT id FROM country where name='Portugal'), 'FC Porto', 'fcporto@fcporto.pt', 'fcporto', 'fcporto', 'https://lh6.googleusercontent.com/-W_SiEiCY4jw/AAAAAAAAAAI/AAAAAAAAkro/7hjZq8nT5TY/s0-c-k-no-ns/photo.jpg', NULL, NULL, NULL, '22 557 0400', current_timestamp);

INSERT INTO client(id, id_entity) VALUES
  (1, 1);
