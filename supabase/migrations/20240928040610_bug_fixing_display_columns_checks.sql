alter table "public"."environment_variables" alter column "value" drop default;

alter table "public"."environments" add constraint "environments_display_check" CHECK ((display ~ '^[a-zA-Z0-9_-](?:[a-zA-Z0-9 _-]*[a-zA-Z0-9_-])?$'::text)) not valid;

alter table "public"."environments" validate constraint "environments_display_check";

alter table "public"."projects" add constraint "projects_display_check" CHECK ((display ~ '^[a-zA-Z0-9_-](?:[a-zA-Z0-9 _-]*[a-zA-Z0-9_-])?$'::text)) not valid;

alter table "public"."projects" validate constraint "projects_display_check";

set check_function_bodies = off;

CREATE OR REPLACE FUNCTION public.environment_variables_before_actions()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$BEGIN

  -- check if NEW.value is NULL
  IF TG_OP != 'DELETE' AND NEW.value IS NULL THEN
    -- set NEW.value to the default from the variables table
    SELECT v.default INTO NEW.value FROM variables v WHERE v.id = NEW.variable_id;
  END IF;

  RETURN COALESCE(NEW, OLD);

END;$function$
;

CREATE TRIGGER environment_variables_before_actions BEFORE INSERT OR DELETE OR UPDATE ON public.environment_variables FOR EACH ROW EXECUTE FUNCTION environment_variables_before_actions();


