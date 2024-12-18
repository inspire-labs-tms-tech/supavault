import {Database} from "./types.gen";
import {SupabaseClient} from "@supabase/supabase-js";
import PostgrestFilterBuilder from "@supabase/postgrest-js/dist/cjs/PostgrestFilterBuilder";

type SupavaultSchema = Database["public"];
export type SupavaultTableName = string & keyof SupavaultSchema["Tables"];
export type SupavaultTable<T extends SupavaultTableName> = SupavaultSchema["Tables"][T];
export type SupavaultTableRow<T extends SupavaultTableName> = SupavaultTable<T>["Row"];
export type SupavaultTableColumn<T extends SupavaultTableName> = string & keyof SupavaultTableRow<T>;
export type SupavaultClient = SupabaseClient<Database, "public", SupavaultSchema>
export type SupavaultFilterChain<Table extends SupavaultTableName> = PostgrestFilterBuilder<Database["public"], Database["public"]["Tables"][Table]["Row"], SupavaultTableRow<Table>[]>;
export type SupavaultQuery<Table extends SupavaultTableName> = (query: SupavaultFilterChain<Table>) => SupavaultFilterChain<Table>;
