
SET SCHEMA 'places4all';

INSERT INTO audit_subgroup(id_audit,id_subgroup) VALUES(1,1),(1,2),(1,3);

INSERT INTO audit_criterion(id_audit,id_criterion,value) VALUES
                (1,1,10),
                (1,2,30),
                (1,3,20);

