SET SCHEMA 'places4all';

CREATE TYPE evaluation_criterion_type AS
  (id INTEGER, w PERCENTAGE, v PERCENTAGE);

CREATE TYPE evaluation_subgroup_type AS
  (id INTEGER, w PERCENTAGE, criteria evaluation_criterion_type[]);

CREATE TYPE evaluation_maingroup_type AS
  (id INTEGER, w PERCENTAGE, subgroups evaluation_subgroup_type[]);

CREATE FUNCTION evaluation_procedure() RETURNS TRIGGER AS $$
DECLARE
  _evals evaluation_maingroup_type[];
  _maingroups_total REAL;
  _subgroups_total REAL;
  _criteria_total REAL;
  _maingroups_sum REAL;
  _subgroups_sum REAL;
  _criteria_sum REAL;
  _rating REAL;
  _m evaluation_maingroup_type;
  _s evaluation_subgroup_type;
  _c evaluation_criterion_type;
BEGIN
  _evals := ARRAY (
    SELECT (maingroup.id, maingroup.weight,
      array_agg((scac.id, scac.weight, scac.criteria)::evaluation_subgroup_type))::evaluation_maingroup_type
    FROM maingroup
    JOIN (
      SELECT subgroup.id, subgroup.weight, subgroup.id_maingroup,
        array_agg((criterion.id, criterion.weight, audit_criterion.value)::evaluation_criterion_type) AS criteria
      FROM audit_criterion
      JOIN criterion ON criterion.id = audit_criterion.id_criterion
      JOIN subgroup ON subgroup.id = criterion.id_subgroup
      WHERE audit_criterion.id_audit = NEW.id_audit
      GROUP BY subgroup.id
    ) scac ON scac.id_maingroup = maingroup.id
    GROUP BY maingroup.id
  );

  _rating := 0;
  _maingroups_total := 0;
  _maingroups_sum := 0;
  FOREACH _m IN ARRAY _evals
  LOOP
    _maingroups_total := _maingroups_total + _m.w;
    _subgroups_total := 0;
    _subgroups_sum := 0;
    FOREACH _s IN ARRAY _m.subgroups
    LOOP
      _subgroups_total := _subgroups_total + _s.w;
      _criteria_total := 0;
      _criteria_sum := 0;
      FOREACH _c IN ARRAY _s.criteria
      LOOP
        _criteria_sum := _criteria_sum + _c.w * _c.v;
        _criteria_total := _criteria_total + _c.w;
      END LOOP;
      _subgroups_sum := _subgroups_sum + _s.w * _criteria_sum / _criteria_total;
    END LOOP;
    _maingroups_sum := _maingroups_sum + _m.w * _subgroups_sum / _subgroups_total;
  END LOOP;
  _rating := _maingroups_sum / _maingroups_total;

  UPDATE audit
  SET rating = _rating::INTEGER
  WHERE audit.id = NEW.id_audit;

  RETURN NULL;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER evaluation_trigger AFTER UPDATE OR INSERT ON audit_criterion
FOR EACH ROW EXECUTE PROCEDURE evaluation_procedure();

CREATE FUNCTION check_audit_criterion_belongs_to_audit_subgroup_procedure() RETURNS TRIGGER AS $$
BEGIN
  PERFORM id_audit
  FROM audit_subgroup
  JOIN criterion ON criterion.id_subgroup = audit_subgroup.id_subgroup
  WHERE audit_subgroup.id_audit = NEW.id_audit AND criterion.id = NEW.id_criterion;
  IF NOT FOUND THEN
    RAISE EXCEPTION 'criterion not in audit subgroups';
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_audit_criterion_belongs_to_audit_subgroup_trigger BEFORE UPDATE OR INSERT ON audit_criterion
FOR EACH ROW EXECUTE PROCEDURE check_audit_criterion_belongs_to_audit_subgroup_procedure();

CREATE FUNCTION initialize_audit_criterion_on_audit_subgroup_procedure() RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO audit_criterion (id_audit, id_criterion, value)
  SELECT NEW.id_audit id, crit, 0 val
  FROM unnest(ARRAY(
    SELECT criterion.id
    FROM criterion
    WHERE criterion.id_subgroup = NEW.id_subgroup
  )) crit;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER initialize_audit_criterion_on_audit_subgroup_trigger AFTER INSERT ON audit_subgroup
FOR EACH ROW EXECUTE PROCEDURE initialize_audit_criterion_on_audit_subgroup_procedure();

CREATE FUNCTION audit_subgroup_audit_criterion_consistency_procedure() RETURNS TRIGGER AS $$
BEGIN
  PERFORM criterion.id
  FROM audit_criterion
  JOIN criterion ON criterion.id = audit_criterion.id_criterion AND criterion.id_subgroup = OLD.id_subgroup
  WHERE audit_criterion.id_audit = OLD.id_audit;
  IF FOUND THEN
    RAISE EXCEPTION 'audit subgroup delete/update not permitted related audit_criterion exists';
  END IF;
  IF TG_OP = 'DELETE' THEN
    return OLD;
  ELSIF TG_OP = 'UPDATE' THEN
    return NEW;
  END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER audit_subgroup_audit_criterion_consistency_trigger BEFORE DELETE OR UPDATE ON audit_subgroup
FOR EACH ROW EXECUTE PROCEDURE audit_subgroup_audit_criterion_consistency_procedure();
