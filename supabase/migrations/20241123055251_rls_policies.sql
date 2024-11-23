create policy "select based on client"
on "public"."projects"
as permissive
for select
to authenticated
using ((id = ( SELECT env.project_id
   FROM environments env
 LIMIT 1)));



