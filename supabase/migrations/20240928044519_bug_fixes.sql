set check_function_bodies = off;

CREATE OR REPLACE FUNCTION public.users_after_actions()
 RETURNS trigger
 LANGUAGE plpgsql
 SECURITY DEFINER
AS $function$BEGIN
  IF TG_OP = 'INSERT' THEN
    INSERT INTO public.clients (id) VALUES (NEW.id);
  END IF;

  RETURN COALESCE (NEW, OLD);
END;$function$
;


