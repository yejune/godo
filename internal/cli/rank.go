package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/yejune/godo/internal/rank"
	"github.com/spf13/cobra"
)

var rankCmd = &cobra.Command{
	Use:   "rank",
	Short: "Interact with the Rank leaderboard API",
	Long: `Rank provides authentication and submission commands for the
Do Rank leaderboard system.`,
}

var rankLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with the Rank API via OAuth",
	RunE: func(cmd *cobra.Command, args []string) error {
		baseURL := os.Getenv("RANK_API_URL")
		if baseURL == "" {
			baseURL = "https://rank.do.dev"
		}
		creds, err := rank.StartOAuthFlow(baseURL)
		if err != nil {
			return fmt.Errorf("oauth flow: %w", err)
		}
		if err := rank.SaveCredentials(creds); err != nil {
			return fmt.Errorf("save credentials: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), "logged in successfully")
		return nil
	},
}

var rankSubmitCmd = &cobra.Command{
	Use:   "submit [session-id]",
	Short: "Submit a session transcript to the leaderboard",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID := args[0]

		creds, err := rank.LoadCredentials()
		if err != nil {
			return fmt.Errorf("load credentials (run 'godo rank login' first): %w", err)
		}

		path := rank.FindTranscriptForSession(sessionID)
		if path == "" {
			return fmt.Errorf("transcript not found for session %s", sessionID)
		}

		usage, err := rank.ParseTranscript(path)
		if err != nil {
			return fmt.Errorf("parse transcript: %w", err)
		}

		client := rank.NewClient(creds.APIKey)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		hash, err := rank.ComputeSessionHash(
			usage.EndedAt,
			rank.ClampTokens(usage.InputTokens),
			rank.ClampTokens(usage.OutputTokens),
		)
		if err != nil {
			return fmt.Errorf("compute hash: %w", err)
		}

		submission := &rank.SessionSubmission{
			SessionHash:  hash,
			EndedAt:      usage.EndedAt,
			InputTokens:  rank.ClampTokens(usage.InputTokens),
			OutputTokens: rank.ClampTokens(usage.OutputTokens),
		}

		if err := client.SubmitSession(ctx, submission); err != nil {
			return fmt.Errorf("submit session: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "session %s submitted\n", sessionID)
		return nil
	},
}

var rankStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current user rank",
	RunE: func(cmd *cobra.Command, args []string) error {
		creds, err := rank.LoadCredentials()
		if err != nil {
			return fmt.Errorf("load credentials (run 'godo rank login' first): %w", err)
		}

		client := rank.NewClient(creds.APIKey)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		info, err := client.GetUserRank(ctx)
		if err != nil {
			return fmt.Errorf("get rank: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "User: %s\nSessions: %d\nTotal Tokens: %d\n",
			info.Username, info.TotalSessions, info.TotalTokens)
		if info.AllTime != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "All-time Position: #%d\n", info.AllTime.Position)
		}
		return nil
	},
}

var rankLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored Rank API credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := rank.DeleteCredentials(); err != nil {
			return fmt.Errorf("delete credentials: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), "logged out")
		return nil
	},
}

func init() {
	rankCmd.AddCommand(rankLoginCmd)
	rankCmd.AddCommand(rankSubmitCmd)
	rankCmd.AddCommand(rankStatusCmd)
	rankCmd.AddCommand(rankLogoutCmd)
	rootCmd.AddCommand(rankCmd)
}
