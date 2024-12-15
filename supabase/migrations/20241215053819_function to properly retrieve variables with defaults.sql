set check_function_bodies = off;

CREATE OR REPLACE FUNCTION public.get_variable_ids()
 RETURNS TABLE(key text, value text, env_var_id uuid)
 LANGUAGE plpgsql
AS $function$
DECLARE
    env_id uuid;
BEGIN
    -- Check if the current user has a matching client and retrieve the environment_id
    SELECT c.environment_id INTO env_id
    FROM public.clients c
    WHERE c.id = auth.uid();

    IF env_id IS NULL THEN
        RAISE EXCEPTION 'No matching client found for the current user.';
    END IF;

    -- Return the variable IDs along with the corresponding trimmed values and environment_variables.id
    RETURN QUERY
    SELECT 
        v.id AS key, 
        TRIM(COALESCE(NULLIF(ev.value, ''), v.default)) AS value,
        ev.id AS env_var_id
    FROM public.variables v
    LEFT JOIN public.environment_variables ev ON ev.environment_id = env_id AND ev.variable_id = v.id;
END;
$function$
;


create schema if not exists "supavault";

create table "supavault"."version_history" (
    "version" text not null,
    "at" timestamp with time zone not null default CURRENT_TIMESTAMP
);


CREATE UNIQUE INDEX version_history_pkey ON supavault.version_history USING btree (version);

alter table "supavault"."version_history" add constraint "version_history_pkey" PRIMARY KEY using index "version_history_pkey";

alter table "supavault"."version_history" add constraint "valid_semver" CHECK ((version ~ '^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-(0|[1-9A-Za-z-][0-9A-Za-z-]*)(\.[0-9A-Za-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$'::text)) not valid;

alter table "supavault"."version_history" validate constraint "valid_semver";


