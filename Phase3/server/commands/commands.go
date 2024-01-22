package commands

import (
	"net"

	"github.com/reeflective/console"
	"github.com/spf13/cobra"

	con "phase3/server/console"
)

func Commands() console.Commands {
	return func() *cobra.Command {
		root := &cobra.Command{}
		root.Short = "short help"

		// listen
		listen := &cobra.Command{
			Use:   "listen",
			Short: "Start a tcp listener",
			Run: func(cmd *cobra.Command, args []string) {
				h, _ := cmd.Flags().GetIP("lhost")
				p, _ := cmd.Flags().GetUint32("lport")

				listen(h, p)
			},
		}

		listen.Flags().IPP("lhost", "l", net.IPv4(0, 0, 0, 0), "local ip (default 0.0.0.0)")
		listen.Flags().Uint32P("lport", "p", 1334, "tcp listening port")
		root.AddCommand(listen)

		// job
		job := &cobra.Command{
			Use:   "job",
			Short: "Show all background jobs",
			Run: func(cmd *cobra.Command, args []string) {
				k, _ := cmd.Flags().GetInt("kill")
				job(k)
			},
		}
		job.Flags().IntP("kill", "k", 0, "kill a background job by id")
		root.AddCommand(job)

		// session
		session := &cobra.Command{
			Use:   "session",
			Short: "Show all sessions and interact with them",
			Run: func(cmd *cobra.Command, args []string) {
				k, _ := cmd.Flags().GetInt("kill")
				i, _ := cmd.Flags().GetInt("interact")

				session(k, i)
			},
		}
		session.Flags().IntP("kill", "k", 0, "kill a session by id")
		session.Flags().IntP("interact", "i", 0, "interact with a session by id")
		root.AddCommand(session)

		root.CompletionOptions.DisableDefaultCmd = true
		root.DisableFlagsInUseLine = true

		return root
	}
}

func SessionCmd() console.Commands {
	return func() *cobra.Command {
		rootCmd := &cobra.Command{}

		bgCmd := &cobra.Command{
			Use:   "bg",
			Short: "back to main menu (Ctrl-D)",
			Run: func(cmd *cobra.Command, args []string) {
				PrintInfo("Switch back to main menu\n")
				con.Con.SessionAct = 0
				con.Con.App.SwitchMenu("")
			},
		}
		rootCmd.AddCommand(bgCmd)

		lsCmd := &cobra.Command{
			Use:   "ls",
			Short: "list directory content",
			Run: func(cmd *cobra.Command, args []string) {
				ls()
			},
		}
		rootCmd.AddCommand(lsCmd)

		rootCmd.CompletionOptions.DisableDefaultCmd = true
		rootCmd.DisableFlagsInUseLine = true

		return rootCmd
	}
}
