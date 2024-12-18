"use client";
import {useSearchParams} from "next/navigation";

export default function Home() {

  const params = useSearchParams();
  const SUPABASE_URL = params.get("supabase-url");
  const SUPABASE_ANON_KEY = params.get("supabase-anon-key");

  console.log(SUPABASE_URL, SUPABASE_ANON_KEY);

  return (
    <div>
    </div>
  );
}
