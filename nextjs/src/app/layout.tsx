import "./globals.css";
import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';
import {ReactNode} from "react";


type Props = Readonly<{
  children: ReactNode;
}>;

export default function RootLayout({children}: Props) {
  return (
    <html lang="en">
    <head>
      <title>Supavault</title>
      <meta name="description" content="A Supabase-backed Keystore"/>
      <meta name="viewport" content="width=device-width, initial-scale=1"/>
    </head>
    <body>
    <main>{children}</main>
    </body>
    </html>
  );
}
