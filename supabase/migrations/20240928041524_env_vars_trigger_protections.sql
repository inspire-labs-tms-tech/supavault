set check_function_bodies = off;

CREATE OR REPLACE FUNCTION public.environment_variables_before_actions()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$DECLARE
  _env_project_id uuid := NULL;
  _var_project_id uuid := NULL;
BEGIN

  -- check if NEW.value is NULL
  IF TG_OP != 'DELETE' THEN

    -- project must match
    SELECT r.project_id INTO _env_project_id FROM environments r WHERE r.id = NEW.environment_id;
    SELECT r.project_id INTO _var_project_id FROM variables r WHERE r.id = NEW.variable_id; 
    IF _env_project_id IS NULL OR _var_project_id IS NULL THEN
      RAISE EXCEPTION '[environment_variables_before_actions] invalid environment or variable configuration';
    END IF;
    IF _env_project_id <> _var_project_id THEN
      RAISE EXCEPTION '[environment_variables_before_actions]  variable is not from this environment''s project';
    END IF;

    -- set NEW.value to the default from the variables table
    IF NEW.value IS NULL THEN
      SELECT v.default INTO NEW.value FROM variables v WHERE v.id = NEW.variable_id;
    END IF;

  END IF;

  RETURN COALESCE(NEW, OLD);

END;$function$
;


