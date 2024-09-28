revoke delete on table "public"."secrets" from "anon";

revoke insert on table "public"."secrets" from "anon";

revoke references on table "public"."secrets" from "anon";

revoke select on table "public"."secrets" from "anon";

revoke trigger on table "public"."secrets" from "anon";

revoke truncate on table "public"."secrets" from "anon";

revoke update on table "public"."secrets" from "anon";

revoke delete on table "public"."secrets" from "authenticated";

revoke insert on table "public"."secrets" from "authenticated";

revoke references on table "public"."secrets" from "authenticated";

revoke select on table "public"."secrets" from "authenticated";

revoke trigger on table "public"."secrets" from "authenticated";

revoke truncate on table "public"."secrets" from "authenticated";

revoke update on table "public"."secrets" from "authenticated";

revoke delete on table "public"."secrets" from "service_role";

revoke insert on table "public"."secrets" from "service_role";

revoke references on table "public"."secrets" from "service_role";

revoke select on table "public"."secrets" from "service_role";

revoke trigger on table "public"."secrets" from "service_role";

revoke truncate on table "public"."secrets" from "service_role";

revoke update on table "public"."secrets" from "service_role";

alter table "public"."secrets" drop constraint "secrets_id_check";

alter table "public"."secrets" drop constraint "secrets_project_id_fkey";

alter table "public"."secrets" drop constraint "secrets_pkey";

drop index if exists "public"."secrets_pkey";

drop table "public"."secrets";


