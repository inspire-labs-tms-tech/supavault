package commands

import (
	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/fatih/color"
	commands "github.com/inspire-labs-tms-tech/supavault/pkg/helpers"
	"github.com/inspire-labs-tms-tech/supavault/pkg/helpers/auth/secrets"
	"github.com/inspire-labs-tms-tech/supavault/pkg/helpers/supabase"
	"github.com/inspire-labs-tms-tech/supavault/pkg/models/auth"
	"github.com/urfave/cli/v2"
	"strings"
)

const DEFAULT_SERVER_URL = "http://127.0.0.1:54321"
const DEFAULT_ANON_KEY = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZS1kZW1vIiwicm9sZSI6ImFub24iLCJleHAiOjE5ODM4MTI5OTZ9.CRXP1A7WOeoJeXxjNni43kdQwgnWNReilDMblYTn_I0"

var AuthCommand = &cli.Command{
	Name:        "auth",
	Description: "manage authentication status",
	Subcommands: []*cli.Command{
		{
			Name:        "show",
			Description: "show existing saved credentials",
			Usage:       "show existing saved credentials",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "verify",
					Usage:       "verify the existing login credentials by attempting to authenticate",
					DefaultText: "false",
					Value:       false,
				},
			},
			Action: func(c *cli.Context) error {

				secret, err := secrets.GetSecret()
				if err != nil {
					color.Red(err.Error())
					return err
				}

				verify := c.Bool("verify")
				var verified bool
				if verify {
					client, clientErr := supabase.GetClient()
					if clientErr != nil {
						color.Red("verification failed - an error occurred while creating a client: %s", clientErr.Error())
						return nil
					} else {
						v, err := client.Authenticate()
						verified = v
						if err != nil {
							color.Red("verification failed - an error occurred during authentication: %s", err.Error())
							return nil
						}
					}
				}

				color.Blue("   email: %s", secret.Email)
				if secret.Password != "" {
					color.Blue("password: %s", strings.Repeat("*", len(secret.Password)))
				} else {
					color.RGB(255, 128, 0).Println("password: (missing)")
				}

				color.Blue("     url: %s", secret.Url)
				color.Blue("    anon: %s", secret.AnonKey)

				if verify {
					print(color.BlueString("verified: "))
					if verified {
						color.Green("✓")
					} else {
						color.Red("failed")
					}
				}

				return nil
			},
		},
		{
			Name:        "logout",
			Description: "remove authentication credentials",
			Usage:       "remove authentication credentials",
			Action: func(c *cli.Context) error {
				_ = secrets.RemoveSecret()
				color.Blue("logout successful")
				return nil
			},
		},
		{
			Name:        "login",
			Description: "save authentication credentials",
			Usage:       "save authentication credentials",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "email",
					Aliases: []string{"u"},
					Usage:   "the email address of the client to authenticate with",
				},
				&cli.StringFlag{
					Name:        "url",
					Usage:       "the URL of the Supabase instance to authenticate with",
					DefaultText: DEFAULT_SERVER_URL,
				},
				&cli.StringFlag{
					Name:        "anon",
					Usage:       "the Anon(ymous) Key of the Supabase instance to authenticate with",
					DefaultText: DEFAULT_ANON_KEY,
				},
				&cli.StringFlag{
					Name:        "password",
					Usage:       "the password address of the client to authenticate with",
					DefaultText: DEFAULT_ANON_KEY,
				},
			},
			Action: func(c *cli.Context) error {

				verifier := emailverifier.NewVerifier()
				email, _ := commands.Prompt(c, "email")
				if res, e := verifier.Verify(email); e != nil {
					color.Red(e.Error())
				} else if !res.Syntax.Valid {
					color.Red("invalid email address: %s", email)
					return nil
				}

				pass, _ := commands.Prompt(c, "password")
				url, _ := commands.PromptWithDefault(c, "url", DEFAULT_SERVER_URL)
				anon, _ := commands.PromptWithDefault(c, "anon", DEFAULT_ANON_KEY)

				err := secrets.SetSecret(auth.ClientCredentials{
					Password: pass,
					AnonKey:  anon,
					Url:      url,
					Email:    email,
				})

				if err != nil {
					color.Red(err.Error())
				} else {
					color.Green("login successful")
				}

				return nil
			},
		},
	},
}
