import {Database} from "./types.gen";
export {Database as SupavaultDatabase} from "./types.gen";

type SupavaultSchema = Database["public"];
export type SupavaultTableName = string & keyof SupavaultSchema["Tables"];
export type SupavaultTable<T extends SupavaultTableName> = SupavaultSchema["Tables"][T];
export type SupavaultTableRow<T extends SupavaultTableName> = SupavaultTable<T>["Row"];

// Define a type for a comma-separated list of valid column names for the specified table
type ColumnString<T extends SupavaultTableName> =
    ValidColumnNames<T> extends infer C
        ? C extends string
            ? `${C}` | `${C}, ${ColumnString<T>}`
            : never
        : never;

// Extract valid column names for the specified table
type ValidColumnNames<T extends SupavaultTableName> = keyof SupavaultTableRow<T>;

// Function to validate and handle the query
function queryTable<Table extends SupavaultTableName, Columns extends ColumnString<Table>>(
    table: Table,
    columns: Columns
): Pick<SupavaultTableRow<Table>, Extract<Columns, ValidColumnNames<Table>>> {
    const columnArray = columns.split(",").map((col) => col.trim()) as Array<ValidColumnNames<Table>>;

    // Example runtime validation
    const validColumns = Object.keys(SupavaultSchema["Tables"][table]["Row"]) as Array<ValidColumnNames<Table>>;

    columnArray.forEach((column) => {
        if (!validColumns.includes(column)) {
            throw new Error(`Invalid column: ${column}`);
        }
    });

    // Dummy return for type inference demonstration
    return {} as any;
}

queryTable("clients", "id")