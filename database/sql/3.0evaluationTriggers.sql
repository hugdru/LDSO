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
  _id_audit INTEGER := 0;
BEGIN
  SET search_path TO places4all, public;

  IF TG_OP = 'DELETE' THEN
    _id_audit := OLD.id_audit;
  ELSIF TG_OP = 'UPDATE' OR TG_OP = 'INSERT' THEN
    _id_audit := NEW.id_audit;
  ELSE
    RAISE EXCEPTION 'EXPECTING ONLY (DELETE|UPDATE|INSERT)';
  END IF;

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
      WHERE audit_criterion.id_audit = _id_audit
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
  WHERE audit.id = _id_audit;

  RETURN NULL;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER evaluation_trigger AFTER UPDATE OR INSERT OR DELETE ON audit_criterion
FOR EACH ROW EXECUTE PROCEDURE evaluation_procedure();

CREATE FUNCTION check_audit_criterion_belongs_to_audit_subgroup_procedure() RETURNS TRIGGER AS $$
BEGIN
  SET search_path TO places4all, public;

  PERFORM id_audit
  FROM audit_subgroup
  JOIN criterion ON criterion.id_subgroup = audit_subgroup.id_subgroup
  WHERE audit_subgroup.id_audit = NEW.id_audit AND criterion.id = NEW.id_criterion;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'CRITERION NOT IN AUDIT SUBGROUPS';
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER check_audit_criterion_belongs_to_audit_subgroup_trigger BEFORE UPDATE OR INSERT ON audit_criterion
FOR EACH ROW EXECUTE PROCEDURE check_audit_criterion_belongs_to_audit_subgroup_procedure();

-- CREATE FUNCTION initialize_audit_criterion_on_audit_subgroup_procedure() RETURNS TRIGGER AS $$
-- BEGIN
--   SET search_path TO places4all, public;
--
--   INSERT INTO audit_criterion (id_audit, id_criterion, value)
--   SELECT NEW.id_audit id, crit, 0 val
--   FROM unnest(ARRAY(
--     SELECT criterion.id
--     FROM criterion
--     WHERE criterion.id_subgroup = NEW.id_subgroup
--   )) crit;
--
--   RETURN NULL;
-- END;
-- $$ LANGUAGE plpgsql;
-- CREATE TRIGGER initialize_audit_criterion_on_audit_subgroup_trigger AFTER INSERT ON audit_subgroup
-- FOR EACH ROW EXECUTE PROCEDURE initialize_audit_criterion_on_audit_subgroup_procedure();

CREATE FUNCTION audit_subgroup_audit_criterion_consistency_procedure() RETURNS TRIGGER AS $$
BEGIN
  SET search_path TO places4all, public;

  PERFORM criterion.id
  FROM audit_criterion
  JOIN criterion ON criterion.id = audit_criterion.id_criterion AND criterion.id_subgroup = OLD.id_subgroup
  WHERE audit_criterion.id_audit = OLD.id_audit;

  IF FOUND THEN
    RAISE EXCEPTION 'AUDIT SUBGROUP DELETE/UPDATE NOT PERMITTED RELATED AUDIT_CRITERION EXISTS';
  END IF;

  IF TG_OP = 'DELETE' THEN
    return OLD;
  ELSIF TG_OP = 'UPDATE' THEN
    return NEW;
  ELSE
    RAISE 'EXPECTING ONLY (DELETE | UPDATE)';
  END IF;

  RETURN NULL;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER audit_subgroup_audit_criterion_consistency_trigger BEFORE DELETE OR UPDATE ON audit_subgroup
FOR EACH ROW EXECUTE PROCEDURE audit_subgroup_audit_criterion_consistency_procedure();

CREATE FUNCTION audit_subgroup_belongs_to_audit_template_procedure() RETURNS TRIGGER AS $$
BEGIN
  SET search_path TO places4all, public;

  PERFORM subgroup.id
  FROM subgroup
  JOIN maingroup ON maingroup.id = subgroup.id_maingroup
  JOIN template ON template.id = maingroup.id_template
  WHERE subgroup.id = NEW.id_subgroup AND template.id = (SELECT id_template FROM audit WHERE audit.id = NEW.id_audit);

  IF NOT FOUND THEN
    RAISE EXCEPTION 'AUDIT SUBGROUP DOES NOT BELONG TO AUDIT TEMPLATE';
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER audit_subgroup_belongs_to_audit_template_trigger BEFORE UPDATE OR INSERT ON audit_subgroup
FOR EACH ROW EXECUTE PROCEDURE audit_subgroup_belongs_to_audit_template_procedure();

CREATE FUNCTION block_maingroup_if_template_closed_procedure() RETURNS TRIGGER AS $$
BEGIN
  SET search_path TO places4all, public;

  PERFORM maingroup.id
  FROM maingroup
  JOIN template ON template.id = maingroup.id_template AND template.closed_date IS NULL AND template.closed IS FALSE
  WHERE maingroup.id = OLD.id;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'TEMPLATE IS CLOSED OPERATION NOT PERMITTED';
  END IF;

  IF TG_OP = 'DELETE' THEN
    return OLD;
  ELSIF TG_OP = 'UPDATE' THEN
    return NEW;
  END IF;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER block_maingroup_if_template_closed_trigger BEFORE UPDATE OR DELETE ON maingroup FOR EACH ROW EXECUTE PROCEDURE block_maingroup_if_template_closed_procedure();

CREATE FUNCTION block_subgroup_if_template_closed_procedure() RETURNS TRIGGER AS $$
BEGIN
  SET search_path TO places4all, public;

  PERFORM subgroup.id
  FROM subgroup
  JOIN maingroup ON maingroup.id = subgroup.id_maingroup
  JOIN template ON template.id = maingroup.id_template AND template.closed_date IS NULL AND template.closed IS FALSE
  WHERE subgroup.id = OLD.id;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'template is closed operation not permitted';
  END IF;

  IF TG_OP = 'DELETE' THEN
    return OLD;
  ELSIF TG_OP = 'UPDATE' THEN
    return NEW;
  END IF;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER block_subgroup_if_template_closed_trigger BEFORE UPDATE OR DELETE ON subgroup FOR EACH ROW EXECUTE PROCEDURE block_subgroup_if_template_closed_procedure();

CREATE FUNCTION block_criterion_if_template_closed_procedure() RETURNS TRIGGER AS $$
BEGIN
  SET search_path TO places4all, public;

  PERFORM criterion.id
  FROM criterion
  JOIN subgroup ON subgroup.id = criterion.id_subgroup
  JOIN maingroup ON maingroup.id = subgroup.id_maingroup
  JOIN template ON template.id = maingroup.id_template AND template.closed_date IS NULL AND template.closed IS FALSE
  WHERE criterion.id = OLD.id;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'TEMPLATE IS CLOSED OPERATION NOT PERMITTED';
  END IF;

  IF TG_OP = 'DELETE' THEN
    return OLD;
  ELSIF TG_OP = 'UPDATE' THEN
    return NEW;
  END IF;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER block_criterion_if_template_closed_trigger BEFORE UPDATE OR DELETE ON criterion FOR EACH ROW EXECUTE PROCEDURE block_criterion_if_template_closed_procedure();

CREATE FUNCTION block_criterion_accessibility_if_template_closed_procedure() RETURNS TRIGGER AS $$
BEGIN
  SET search_path TO places4all, public;

  PERFORM criterion.id
  FROM criterion
  JOIN subgroup ON subgroup.id = criterion.id_subgroup
  JOIN maingroup ON maingroup.id = subgroup.id_maingroup
  JOIN template ON template.id = maingroup.id_template AND template.closed_date IS NULL AND template.closed IS FALSE
  WHERE criterion.id = OLD.id_criterion;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'template is closed operation not permitted';
  END IF;

  IF TG_OP = 'DELETE' THEN
    return OLD;
  ELSIF TG_OP = 'UPDATE' THEN
    return NEW;
  END IF;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER block_criterion_accessibility_if_template_closed_trigger BEFORE UPDATE OR DELETE ON criterion_accessibility FOR EACH ROW EXECUTE PROCEDURE block_criterion_accessibility_if_template_closed_procedure();

CREATE FUNCTION only_closed_templates_on_audit_procedure() RETURNS TRIGGER AS $$
BEGIN
  SET search_path TO places4all, public;

  IF TG_OP = 'INSERT' AND NEW.id_template IS NULL THEN
    RAISE EXCEPTION 'AUDIT INSERT MUST ALWAYS COME WITH ID_TEMPLATE';
  END IF;

  IF TG_OP = 'UPDATE' AND NEW.id_template != OLD.id_template THEN
    RAISE EXCEPTION 'ONCE ID_TEMPLATE SET IN AUDIT THERE IS NO WAY BACK';
  END IF;

  PERFORM template.id
  FROM template
  WHERE template.id = NEW.id_template AND template.closed_date IS NOT NULL AND template.closed IS TRUE;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'AUDIT ONLY ACCEPTS CLOSED TEMPLATES';
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER only_closed_templates_on_audit_trigger BEFORE UPDATE OR INSERT ON audit
FOR EACH ROW EXECUTE PROCEDURE only_closed_templates_on_audit_procedure();

CREATE FUNCTION template_closed_guard_procedure() RETURNS TRIGGER AS $$
BEGIN
  SET search_path TO places4all, public;

  IF NEW.closed THEN
    IF NEW.closed_date IS NULL THEN
      RAISE EXCEPTION 'NEW ROW SHOULD ASSIGN A DATE TO CLOSED_DATE';
    END IF;

    IF TG_OP = 'UPDATE' THEN
      IF OLD.closed = NEW.closed THEN
        RAISE EXCEPTION 'IS CLOSED';
      END IF;
    END IF;

    UPDATE maingroup
    SET closed = NEW.closed
    WHERE maingroup.id_template = NEW.id;

    UPDATE subgroup
    SET closed = NEW.closed
    FROM maingroup
    WHERE maingroup.id_template = NEW.id AND subgroup.id_maingroup = maingroup.id;

    UPDATE criterion
    SET closed = NEW.closed
    FROM subgroup
    JOIN maingroup ON maingroup.id_template = NEW.id AND subgroup.id_maingroup = maingroup.id
    WHERE criterion.id_subgroup = subgroup.id;

    UPDATE criterion_accessibility
    SET closed = NEW.closed
    FROM criterion
    JOIN subgroup ON subgroup.id = criterion.id_subgroup
    JOIN maingroup ON maingroup.id_template = NEW.id AND subgroup.id_maingroup = maingroup.id
    WHERE criterion.id_subgroup = subgroup.id AND criterion_accessibility.id_criterion = criterion.id;
  ELSE
    IF NEW.closed_date IS NOT NULL THEN
      RAISE EXCEPTION 'NEW ROW SHOULD ASSIGN NULL TO CLOSED_DATE';
    END IF;

    DECLARE
      _template_id INTEGER := 0;
    BEGIN
      IF TG_OP = 'INSERT' THEN
        _template_id := NEW.id;
      ELSIF TG_OP = 'UPDATE' THEN
        _template_id := OLD.id;
      END IF;

      PERFORM audit.id
      FROM audit
      WHERE audit.id_template = _template_id;

      IF FOUND THEN
        RAISE EXCEPTION 'CANNOT SET CLOSED TO FALSE BECAUSE THE TEMPLATE IS USED IN AUDITS';
      END IF;
    END;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER template_closed_guard_trigger BEFORE UPDATE OR INSERT ON template
FOR EACH ROW EXECUTE PROCEDURE template_closed_guard_procedure();

CREATE FUNCTION cascade_closed_on_change_template_procedure() RETURNS TRIGGER AS $$
BEGIN
  IF OLD.closed = TRUE AND NEW.closed = FALSE THEN
    UPDATE maingroup
    SET closed = NEW.closed
    WHERE maingroup.id_template = NEW.id;

    UPDATE subgroup
    SET closed = NEW.closed
    FROM maingroup
    WHERE maingroup.id_template = NEW.id AND subgroup.id_maingroup = maingroup.id;

    UPDATE criterion
    SET closed = NEW.closed
    FROM subgroup
    JOIN maingroup ON maingroup.id_template = NEW.id AND subgroup.id_maingroup = maingroup.id
    WHERE criterion.id_subgroup = subgroup.id;

    UPDATE criterion_accessibility
    SET closed = NEW.closed
    FROM criterion
    JOIN subgroup ON subgroup.id = criterion.id_subgroup
    JOIN maingroup ON maingroup.id_template = NEW.id AND subgroup.id_maingroup = maingroup.id
    WHERE criterion.id_subgroup = subgroup.id AND criterion_accessibility.id_criterion = criterion.id;
  END IF;

  RETURN NULL;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER cascade_closed_on_change_template_trigger AFTER UPDATE ON template
FOR EACH ROW EXECUTE PROCEDURE cascade_closed_on_change_template_procedure();
