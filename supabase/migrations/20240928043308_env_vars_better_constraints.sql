set check_function_bodies = off;

CREATE OR REPLACE FUNCTION public.environment_variables_after_actions()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$BEGIN
  IF TG_OP = 'DELETE' THEN

    -- cannot delete if environment and variable exists
    IF ((SELECT COUNT(*) FROM environments r WHERE r.id = OLD.environment_id) IS DISTINCT FROM 0) AND ((SELECT COUNT(*) FROM variables r WHERE r.id = OLD.variable_id) IS DISTINCT FROM 0) THEN
      RAISE EXCEPTION '[environment_variables_after_actions] cannot delete environment_variable for an existing environment or variable';
    END IF;

  END IF;

  RETURN COALESCE (NEW, OLD);
END;$function$
;

CREATE OR REPLACE FUNCTION public.environments_after_actions()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$DECLARE
  ROW RECORD := NULL;
BEGIN

  IF TG_OP = 'INSERT' THEN
    FOR ROW IN SELECT * FROM variables r WHERE r.project_id = NEW.project_id 
    LOOP
      INSERT INTO environment_variables (
        environment_id,
        variable_id
      ) VALUES (
        NEW.id,
        ROW.id
      );
    END LOOP;
  END IF;

  RETURN COALESCE (NEW, OLD);

END;$function$
;

CREATE OR REPLACE FUNCTION public.variables_after_actions()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$DECLARE
  ROW RECORD := NULL;
BEGIN

  IF TG_OP = 'INSERT' THEN
    FOR ROW IN SELECT * FROM environments r WHERE r.project_id = NEW.project_id 
    LOOP
      INSERT INTO environment_variables (
        environment_id,
        variable_id
      ) VALUES (
        ROW.id,
        NEW.id
      );
    END LOOP;
  END IF;

  RETURN COALESCE (NEW, OLD);

END;$function$
;

CREATE TRIGGER environment_variables_after_actions AFTER INSERT OR DELETE OR UPDATE ON public.environment_variables FOR EACH ROW EXECUTE FUNCTION environment_variables_after_actions();

CREATE TRIGGER environments_after_actions AFTER INSERT OR DELETE OR UPDATE ON public.environments FOR EACH ROW EXECUTE FUNCTION environments_after_actions();

CREATE TRIGGER variables_after_actions AFTER INSERT OR DELETE OR UPDATE ON public.variables FOR EACH ROW EXECUTE FUNCTION variables_after_actions();


