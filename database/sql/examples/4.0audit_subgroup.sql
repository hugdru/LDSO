
SET SCHEMA 'places4all';

INSERT INTO audit_subgroup(id_audit,id_subgroup) VALUES(1,1),(1,2),(1,3);

INSERT INTO audit_criterion(id_audit,id_criterion,value,observation) VALUES
                (1,1,10,'criterio 1'),
                (1,2,30,'criterio 2'),
                (1,3,20,'criterio 3');

