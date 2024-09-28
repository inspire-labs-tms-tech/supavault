set check_function_bodies = off;

CREATE OR REPLACE FUNCTION public.users_after_actions()
 RETURNS trigger
 LANGUAGE plpgsql
 SECURITY DEFINER
AS $function$BEGIN
  IF TG_OP = 'INSERT' THEN
    INSERT INTO public.clients (id) VALUES (NEW.id);
  END IF;
END;$function$
;


CREATE TRIGGER users_after_actions AFTER INSERT OR DELETE OR UPDATE ON auth.users FOR EACH ROW EXECUTE FUNCTION users_after_actions();


