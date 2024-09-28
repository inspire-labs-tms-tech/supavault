create table "public"."secrets" (
    "id" text not null,
    "created_at" timestamp with time zone not null default (now() AT TIME ZONE 'utc'::text),
    "description" text not null default ''::text
);


CREATE UNIQUE INDEX secrets_pkey ON public.secrets USING btree (id);

alter table "public"."secrets" add constraint "secrets_pkey" PRIMARY KEY using index "secrets_pkey";

alter table "public"."secrets" add constraint "secrets_id_check" CHECK ((id ~ '^[A-Z_]{1,}[A-Z0-9_]{0,}$'::text)) not valid;

alter table "public"."secrets" validate constraint "secrets_id_check";

grant delete on table "public"."secrets" to "anon";

grant insert on table "public"."secrets" to "anon";

grant references on table "public"."secrets" to "anon";

grant select on table "public"."secrets" to "anon";

grant trigger on table "public"."secrets" to "anon";

grant truncate on table "public"."secrets" to "anon";

grant update on table "public"."secrets" to "anon";

grant delete on table "public"."secrets" to "authenticated";

grant insert on table "public"."secrets" to "authenticated";

grant references on table "public"."secrets" to "authenticated";

grant select on table "public"."secrets" to "authenticated";

grant trigger on table "public"."secrets" to "authenticated";

grant truncate on table "public"."secrets" to "authenticated";

grant update on table "public"."secrets" to "authenticated";

grant delete on table "public"."secrets" to "service_role";

grant insert on table "public"."secrets" to "service_role";

grant references on table "public"."secrets" to "service_role";

grant select on table "public"."secrets" to "service_role";

grant trigger on table "public"."secrets" to "service_role";

grant truncate on table "public"."secrets" to "service_role";

grant update on table "public"."secrets" to "service_role";


