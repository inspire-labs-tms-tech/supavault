"use client";

import {
  SupavaultClient,
  SupavaultFilterChain,
  SupavaultQuery,
  SupavaultTableName,
  SupavaultTableRow
} from "@supavault/typescript/types/supabase";
import {uuid} from "@supabase/supabase-js/dist/main/lib/helpers";
import {useEffect, useMemo, useState} from "react";
import {PostgrestMaybeSingleResponse, PostgrestResponse, PostgrestSingleResponse} from "@supabase/supabase-js";

type Modes = "single" | "maybe" | "many";
type Response<Table extends SupavaultTableName, Mode extends Modes> =
  Mode extends "single" ? PostgrestSingleResponse<SupavaultTableRow<Table>>["data"] :
    Mode extends "maybe" ? PostgrestMaybeSingleResponse<SupavaultTableRow<Table>>["data"] :
      Mode extends "many" ? PostgrestResponse<SupavaultTableRow<Table>>["data"] :
        never;

type LoadingState = {
  loading: true;
  value: undefined;
  error: undefined;
};

type ErrorState = {
  loading: false;
  value: null;
  error: true;
}

type DataState<Table extends SupavaultTableName, Mode extends Modes> = {
  loading: false;
  value: Response<Table, Mode>;
  error: false;
}

type RealtimeData<Table extends SupavaultTableName, Mode extends Modes> =
  LoadingState
  | ErrorState
  | DataState<Table, Mode>;

type Props<Table extends SupavaultTableName, Mode extends Modes> = {
  table: Table;
  filter: SupavaultQuery<Table>;
  type: Mode;
  client: SupavaultClient;
}

const refresh = <Table extends SupavaultTableName, Mode extends Modes>(config: Props<Table, Mode>, callback: (result: DataState<Table, Mode>) => any) => {

  let query = config.filter(config.client.from(config.table).select());
  switch (config.type) {
    case "single":
      query = query.single() as unknown as SupavaultFilterChain<Table>;
      break;
    case "maybe":
      query = query.maybeSingle() as unknown as SupavaultFilterChain<Table>;
      break;
    case "many":
      break;
    default:
      throw new Error(`Unsupported query type: ${config.type}`);
  }

  query.then(({data, error}) => error != null || data == null ? ({
    loading: false,
    value: null,
    error: true,
  }) : ({
    loading: false,
    value: data as unknown as Response<Table, Mode>,
    error: false
  }));
}

export const useRealtimeData = <Table extends SupavaultTableName, Mode extends Modes>(props: Props<Table, Mode>): RealtimeData<Table, Mode> => {

  const channel = useMemo(() => uuid(), []);
  const [state, setState] = useState<RealtimeData<Table, Mode>>({
    loading: true,
    value: undefined,
    error: undefined,
  });

  useEffect(() => {
    const subscription = props.client.channel(channel)
      .on("postgres_changes", {
        event: "*",
        schema: "public",
        table: props.table
      }, () => refresh(props, setState))
      .subscribe((status, err) => {
        if (err) console.error({status, err});
        refresh(props, setState); // set the initial/current state
      });

    return () => {
      subscription.unsubscribe();
    }
  }, []);

  return state;

}
