SET SCHEMA 'places4all';

INSERT INTO template (name,description,created_date)
VALUES ('Modelo Teste','template teste 1','2016-12-25T14:47:26.680513Z'); -- 3


INSERT INTO audit(id_property, id_auditor, id_template,rating,observation,created_date ) VALUES
    (1,4,1,0,'audit 1','2016-12-25T14:47:26.680513Z'),
    (1,4,1,0,'audit 1','2016-12-26T14:47:26.680513Z');
