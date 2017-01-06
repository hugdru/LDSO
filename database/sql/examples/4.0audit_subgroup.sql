
SET SCHEMA 'places4all';

INSERT INTO audit_subgroup(id_audit,id_subgroup) VALUES(1,1),(1,2),(1,3);

INSERT INTO audit_criterion(id_audit,id_criterion,value) VALUES
                (1,1,10), -- 1
                (1,2,30), -- 2
                (1,3,20); -- 3

-- UPDATE audit_criterion
-- SET value = ac.value
-- FROM (
--   SELECT * FROM unnest(array[1, 2, 3], array[10, 30, 20])
-- ) AS ac(id_criterion, value)
-- WHERE audit_criterion.id_audit = 1 AND audit_criterion.id_criterion = ac.id_criterion;
