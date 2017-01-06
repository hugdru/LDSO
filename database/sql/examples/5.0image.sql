SET SCHEMA 'places4all';

INSERT INTO remark(id_audit, id_criterion, observation) VALUES(1, 1, 'excelente');

UPDATE audit SET finished_date = '2016-12-26T15:47:26.680513Z' WHERE id = 1;

