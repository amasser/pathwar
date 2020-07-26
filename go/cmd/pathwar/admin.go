package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/peterbourgon/ff"
	"github.com/peterbourgon/ff/ffcli"
	"moul.io/godev"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

var adminJSONFormat bool

func adminCommand() *ffcli.Command {
	adminFlags := flag.NewFlagSet("admin", flag.ExitOnError)
	adminFlags.StringVar(&httpAPIAddr, "http-api-addr", defaultHTTPApiAddr, "HTTP API address")
	adminFlags.StringVar(&ssoOpts.TokenFile, "sso-token-file", ssoOpts.TokenFile, "Token file")
	adminFlags.BoolVar(&adminJSONFormat, "json", false, "Print JSON and exit")
	return &ffcli.Command{
		Name:  "admin",
		Usage: "pathwar [global flags] admin [admin flags] <subcommand> [flags] [args...]",
		Subcommands: []*ffcli.Command{
			adminChallengeCommand(),
			adminUsersCommand(),
			adminAgentsCommand(),
			adminActivitiesCommand(),
			adminOrganizationsCommand(),
			adminTeamsCommand(),
			adminRedumpCommand(),
			adminChallengeAddCommand(),
			adminChallengeFlavorAddCommand(),
			adminChallengeInstanceAddCommand(),
		},
		ShortHelp: "admin commands",
		FlagSet:   adminFlags,
		Options:   []ff.Option{ff.WithEnvVarNoPrefix()},
		Exec:      func([]string) error { return flag.ErrHelp },
	}
}

func adminChallengeCommand() *ffcli.Command {
	adminChallengesFlags := flag.NewFlagSet("admin challenges", flag.ExitOnError)
	return &ffcli.Command{
		Name:    "challenges",
		Usage:   "pathwar [global flags] admin [admin flags] challenges [flags]",
		FlagSet: adminChallengesFlags,
		Exec: func(args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			ctx := context.Background()
			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminListChallenges(ctx, &pwapi.AdminListChallenges_Input{})
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			// tree
			{
				fmt.Println("TREE")
				for _, challenge := range ret.Challenges {
					fmt.Printf("- %s (%d)\n", challenge.Slug, challenge.ID)
					for _, flavor := range challenge.Flavors {
						fmt.Printf("  - %s (%d)\n", flavor.Slug, flavor.ID)
						for _, instance := range flavor.Instances {
							fmt.Printf("    - %d\n", instance.ID)
						}
					}
				}
				fmt.Println("")
			}

			// challenges table
			{
				fmt.Println("CHALLENGES")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"CHALLENGE", "NAME", "AUTHOR", "CREATED", "UPDATED", "FLAVORS", "ID"})
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				table.SetBorder(false)

				for _, challenge := range ret.Challenges {
					slug := challenge.Slug
					name := challenge.Name
					author := challenge.Author
					createdAgo := humanize.Time(*challenge.CreatedAt)
					updatedAgo := humanize.Time(*challenge.UpdatedAt)
					flavors := fmt.Sprintf("%d", len(challenge.Flavors))
					id := fmt.Sprintf("%d", challenge.ID)
					table.Append([]string{slug, name, author, createdAgo, updatedAgo, flavors, id})
				}
				table.Render()
				fmt.Println("")
			}

			// flavors table
			{
				fmt.Println("FLAVORS")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"FLAVOR", "CHALLENGE", "CREATED", "UPDATED", "INSTANCES", "SEASON CHALLENGES", "ID"})
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				table.SetBorder(false)

				for _, challenge := range ret.Challenges {
					for _, flavor := range challenge.Flavors {
						slug := flavor.Slug
						challengeSlug := challenge.Slug
						createdAgo := humanize.Time(*flavor.CreatedAt)
						updatedAgo := humanize.Time(*flavor.UpdatedAt)
						instances := asciiInstancesStats(flavor.Instances)
						seasonChallengeParts := []string{}
						for _, seasonChallenge := range flavor.SeasonChallenges {
							seasonChallengeParts = append(seasonChallengeParts, seasonChallenge.Season.Slug)
						}
						seasonChallenges := strings.Join(seasonChallengeParts, ", ")
						id := fmt.Sprintf("%d", flavor.ID)
						table.Append([]string{slug, challengeSlug, createdAgo, updatedAgo, instances, seasonChallenges, id})
					}
				}
				table.Render()
				fmt.Println("")
			}

			// instances
			{
				fmt.Println("INSTANCES")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"INSTANCE", "FLAVOR", "STATUS", "CREATED", "UPDATED", "CONFIG", "SEASON CHALLENGES", "PRICE/REWARD"})
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				table.SetBorder(false)

				for _, challenge := range ret.Challenges {
					for _, flavor := range challenge.Flavors {
						for _, instance := range flavor.Instances {
							//fmt.Println(godev.PrettyJSONPB(instance))
							id := fmt.Sprintf("%d", instance.ID)
							status := asciiStatus(instance.Status.String())
							flavorSlug := fmt.Sprintf("%s@%s", challenge.Slug, flavor.Slug)
							createdAgo := humanize.Time(*instance.CreatedAt)
							updatedAgo := humanize.Time(*instance.UpdatedAt)
							configStruct, _ := instance.ParseInstanceConfig()
							config := godev.JSONPB(configStruct)
							seasonChallenges := fmt.Sprintf("%d", len(flavor.SeasonChallenges))
							price := "free"
							if flavor.PurchasePrice > 0 {
								price = fmt.Sprintf("$%d", flavor.PurchasePrice)
							}
							priceReward := fmt.Sprintf("%s / $%d", price, flavor.ValidationReward)
							table.Append([]string{id, flavorSlug, status, createdAgo, updatedAgo, config, seasonChallenges, priceReward})
						}
					}
				}
				table.Render()
			}

			return nil
		},
	}
}

func adminUsersCommand() *ffcli.Command {
	adminUsersFlags := flag.NewFlagSet("admin users", flag.ExitOnError)
	return &ffcli.Command{
		Name:    "users",
		Usage:   "pathwar [global flags] admin [admin flags] users [flags]",
		FlagSet: adminUsersFlags,
		Exec: func(args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			ctx := context.Background()
			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminListUsers(ctx, &pwapi.AdminListUsers_Input{})
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			// users table
			{
				fmt.Println("USERS")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"USER", "STATUS", "USERNAME", "EMAIL", "CREATED", "UPDATED", "TEAMS", "ORGS", "ID"})
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				table.SetBorder(false)

				for _, user := range ret.Users {
					//fmt.Println(godev.PrettyJSONPB(user))
					slug := user.Slug
					email := user.Email
					username := user.Username
					status := asciiStatus(user.DeletionStatus.String())
					createdAgo := humanize.Time(*user.CreatedAt)
					updatedAgo := humanize.Time(*user.UpdatedAt)
					teams := fmt.Sprintf("%d", len(user.TeamMemberships))
					orgs := fmt.Sprintf("%d", len(user.OrganizationMemberships))
					id := fmt.Sprintf("%d", user.ID)
					table.Append([]string{slug, status, username, email, createdAgo, updatedAgo, teams, orgs, id})
				}
				table.Render()
				fmt.Println("")
			}

			return nil
		},
	}
}

func adminAgentsCommand() *ffcli.Command {
	adminAgentsFlags := flag.NewFlagSet("admin agents", flag.ExitOnError)
	return &ffcli.Command{
		Name:    "agents",
		Usage:   "pathwar [global flags] admin [admin flags] agents [flags]",
		FlagSet: adminAgentsFlags,
		Exec: func(args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			ctx := context.Background()
			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminListAgents(ctx, &pwapi.AdminListAgents_Input{})
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			// agents table
			{
				fmt.Println("AGENTS")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"AGENT", "HOSTNAME", "SUFFIX", "STATUS", "CREATED", "UPDATED", "SEEN", "STATS", "INSTANCES", "DEFAULT", "ID"})
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				table.SetBorder(false)

				for _, agent := range ret.Agents {
					//fmt.Println(godev.PrettyJSONPB(agent))
					slug := agent.Slug
					createdAgo := humanize.Time(*agent.CreatedAt)
					updatedAgo := humanize.Time(*agent.UpdatedAt)
					seenAgo := humanize.Time(*agent.LastSeenAt)
					id := fmt.Sprintf("%d", agent.ID)
					instances := asciiInstancesStats(agent.ChallengeInstances)
					stats := fmt.Sprintf("%d seen / %d reg.", agent.TimesSeen, agent.TimesRegistered)
					status := asciiStatus(agent.Status.String())
					isDefault := asciiBool(agent.DefaultAgent)
					suffix := agent.DomainSuffix
					hostname := agent.Hostname
					table.Append([]string{slug, hostname, suffix, status, createdAgo, updatedAgo, seenAgo, stats, instances, isDefault, id})
				}
				table.Render()
				fmt.Println("")
			}

			return nil
		},
	}
}

func adminActivitiesCommand() *ffcli.Command {
	adminActivitiesFlags := flag.NewFlagSet("admin activities", flag.ExitOnError)
	return &ffcli.Command{
		Name:    "activities",
		Usage:   "pathwar [global flags] admin [admin flags] activities [flags]",
		FlagSet: adminActivitiesFlags,
		Exec: func(args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			ctx := context.Background()
			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminListActivities(ctx, &pwapi.AdminListActivities_Input{})
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			// activities table
			{
				fmt.Println("ACTIVITIES")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"ID", "KIND", "HAPPENED", "AUTHOR", "TEAM", "USER", "ORG", "SEASON", "CHALLENGE", "COUPON", "SEASON CHAL.", "TEAM MEMBER", "CHALLENGE SUBS"})
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				table.SetBorder(false)

				for _, activity := range ret.Activities {
					//fmt.Println(godev.PrettyJSONPB(activity))
					author := activity.Author.Slug
					id := fmt.Sprintf("%d", activity.ID)
					kind := activity.Kind.String()
					createdAgo := humanize.Time(*activity.CreatedAt)
					team := activity.Team.ASCIIID()
					user := activity.User.ASCIIID()
					organization := activity.Organization.ASCIIID()
					season := activity.Season.ASCIIID()
					challenge := activity.Challenge.ASCIIID()
					coupon := activity.Coupon.ASCIIID()
					seasonChallenge := activity.SeasonChallenge.ASCIIID()
					teamMember := activity.TeamMember.ASCIIID()
					challengeSubscription := activity.ChallengeSubscription.ASCIIID()
					table.Append([]string{id, kind, createdAgo, author, team, user, organization, season, challenge, coupon, seasonChallenge, teamMember, challengeSubscription})
				}
				table.Render()
				fmt.Println("")
			}

			return nil
		},
	}
}

func adminOrganizationsCommand() *ffcli.Command {
	adminOrganizationsFlags := flag.NewFlagSet("admin organizations", flag.ExitOnError)
	return &ffcli.Command{
		Name:    "organizations",
		Usage:   "pathwar [global flags] admin [admin flags] organizations [flags]",
		FlagSet: adminOrganizationsFlags,
		Exec: func(args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			ctx := context.Background()
			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminListOrganizations(ctx, &pwapi.AdminListOrganizations_Input{})
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			// organizations table
			{
				fmt.Println("ORGANIZATIONS")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"SLUG", "NAME", "STATUS", "CREATED", "UPDATED", "TEAMS", "MEMBERS", "ID"})
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				table.SetBorder(false)

				for _, organization := range ret.Organizations {
					//fmt.Println(godev.PrettyJSONPB(organization))
					id := fmt.Sprintf("%d", organization.ID)
					createdAgo := humanize.Time(*organization.CreatedAt)
					updatedAgo := humanize.Time(*organization.UpdatedAt)
					slug := organization.Slug
					name := organization.Name
					status := asciiStatus(organization.DeletionStatus.String())
					teamParts := []string{}
					for _, team := range organization.Teams {
						teamParts = append(teamParts, team.Season.ASCIIID())
					}
					teams := strings.Join(teamParts, ",")
					memberParts := []string{}
					for _, member := range organization.Members {
						memberParts = append(memberParts, member.User.ASCIIID())
					}
					members := strings.Join(memberParts, ",")
					table.Append([]string{slug, name, status, createdAgo, updatedAgo, teams, members, id})
				}
				table.Render()
				fmt.Println("")
			}

			return nil
		},
	}
}

func adminTeamsCommand() *ffcli.Command {
	adminTeamsFlags := flag.NewFlagSet("admin teams", flag.ExitOnError)
	return &ffcli.Command{
		Name:    "teams",
		Usage:   "pathwar [global flags] admin [admin flags] teams [flags]",
		FlagSet: adminTeamsFlags,
		Exec: func(args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			ctx := context.Background()
			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminListTeams(ctx, &pwapi.AdminListTeams_Input{})
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			// teams table
			{
				fmt.Println("TEAMS")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"SLUG", "CREATED", "UPDATED", "SEASON", "MEMBERS", "ORGANIZATION", "CHALLENGES", "ACHIEVEMENTS", "STATUS", "ID"})
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				table.SetBorder(false)

				for _, team := range ret.Teams {
					//fmt.Println(godev.PrettyJSONPB(team))
					id := fmt.Sprintf("%d", team.ID)
					createdAgo := humanize.Time(*team.CreatedAt)
					updatedAgo := humanize.Time(*team.UpdatedAt)
					season := team.Season.ASCIIID()
					memberParts := []string{}
					for _, member := range team.Members {
						memberParts = append(memberParts, member.User.ASCIIID())
					}
					members := strings.Join(memberParts, ",")
					organization := team.Organization.ASCIIID()
					status := asciiStatus(team.DeletionStatus.String())
					challenges := asciiSubscriptionsStats(team.ChallengeSubscriptions)
					achievements := fmt.Sprintf("%d", len(team.Achievements))
					slug := team.Slug
					table.Append([]string{slug, createdAgo, updatedAgo, season, members, organization, challenges, achievements, status, id})
				}
				table.Render()
				fmt.Println("")
			}

			return nil
		},
	}
}

func adminRedumpCommand() *ffcli.Command {
	adminRedumpFlags := flag.NewFlagSet("admin redump", flag.ExitOnError)
	return &ffcli.Command{
		Name:    "redump",
		Usage:   "pathwar [global flags] admin [admin flags] redump [flags] ID...",
		FlagSet: adminRedumpFlags,
		Exec: func(args []string) error {
			if len(args) < 1 {
				return flag.ErrHelp
			}

			if err := globalPreRun(); err != nil {
				return err
			}

			ctx := context.Background()
			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminRedump(ctx, &pwapi.AdminRedump_Input{
				Identifiers: args,
			})
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			fmt.Println("OK")

			return nil
		},
	}
}

func adminChallengeAddCommand() *ffcli.Command {
	adminChallengeAddFlags := flag.NewFlagSet("admin challenge add", flag.ExitOnError)
	adminChallengeAddFlags.StringVar(&adminChallengeAddInput.Challenge.Name, "name", "", "Challenge name")
	adminChallengeAddFlags.StringVar(&adminChallengeAddInput.Challenge.Description, "description", "", "Challenge description")
	adminChallengeAddFlags.StringVar(&adminChallengeAddInput.Challenge.Author, "author", "", "Challenge author")
	adminChallengeAddFlags.StringVar(&adminChallengeAddInput.Challenge.Locale, "locale", "", "Challenge Locale")
	adminChallengeAddFlags.BoolVar(&adminChallengeAddInput.Challenge.IsDraft, "is-draft", true, "Is challenge production ready ?")
	adminChallengeAddFlags.StringVar(&adminChallengeAddInput.Challenge.PreviewUrl, "preview-url", "", "Challenge preview URL")
	adminChallengeAddFlags.StringVar(&adminChallengeAddInput.Challenge.Homepage, "homepage", "", "Challenge homepage URL")

	return &ffcli.Command{
		Name:      "challenge-add",
		Usage:     "pathwar [global flags] admin [admin flags] challenge-add [flags] [args...]",
		ShortHelp: "add a challenge",
		FlagSet:   adminChallengeAddFlags,
		Exec: func(args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			ctx := context.Background()
			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminAddChallenge(ctx, &adminChallengeAddInput)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}
			if globalDebug {
				fmt.Fprintln(os.Stderr, godev.PrettyJSONPB(&ret))
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			fmt.Println(ret.Challenge.ID)
			return nil
		},
	}
}

func adminChallengeFlavorAddCommand() *ffcli.Command {
	adminChallengeFlavorAddFlags := flag.NewFlagSet("admin challenge flavor add", flag.ExitOnError)
	adminChallengeFlavorAddFlags.StringVar(&adminChallengeFlavorAddInput.ChallengeFlavor.Version, "version", "1.0.0", "Challenge flavor version")
	adminChallengeFlavorAddFlags.StringVar(&adminChallengeFlavorAddInput.ChallengeFlavor.ComposeBundle, "compose-bundle", "", "Challenge flavor compose bundle")
	adminChallengeFlavorAddFlags.Int64Var(&adminChallengeFlavorAddInput.ChallengeFlavor.ChallengeID, "challenge-id", 0, "Challenge id")
	return &ffcli.Command{
		Name:      "challenge-flavor-add",
		Usage:     "pathwar [global flags] admin [admin flags] challenge-flavor-add [flags] [args...]",
		ShortHelp: "add a challenge flavor",
		FlagSet:   adminChallengeFlavorAddFlags,
		Exec: func(args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			ctx := context.Background()
			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminAddChallengeFlavor(ctx, &adminChallengeFlavorAddInput)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}
			if globalDebug {
				fmt.Fprintln(os.Stderr, godev.PrettyJSONPB(&ret))
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			fmt.Println(ret.ChallengeFlavor.ID)
			return nil
		},
	}
}
func adminChallengeInstanceAddCommand() *ffcli.Command {
	adminChallengeInstanceAddFlags := flag.NewFlagSet("admin challenge instance add", flag.ExitOnError)
	adminChallengeInstanceAddFlags.Int64Var(&adminChallengeInstanceAddInput.ChallengeInstance.AgentID, "agent-id", 0, "Id of the agent that will host the instance")
	adminChallengeInstanceAddFlags.Int64Var(&adminChallengeInstanceAddInput.ChallengeInstance.FlavorID, "flavor-id", 0, "Challenge flavor id")
	return &ffcli.Command{
		Name:      "challenge-instance-add",
		Usage:     "pathwar [global flags] admin [admin flags] challenge-instance-add [flags] [args...]",
		ShortHelp: "add a challenge instance",
		FlagSet:   adminChallengeInstanceAddFlags,
		Exec: func(args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			ctx := context.Background()
			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminAddChallengeInstance(ctx, &adminChallengeInstanceAddInput)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}
			if globalDebug {
				fmt.Fprintln(os.Stderr, godev.PrettyJSONPB(&ret))
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			fmt.Println(ret.ChallengeInstance.ID)
			return nil
		},
	}
}

func asciiInstancesStats(instances []*pwdb.ChallengeInstance) string {
	if len(instances) == 0 {
		return "🚫"
	}

	instanceGreen := 0
	instanceRed := 0
	for _, instance := range instances {
		if instance.Status == pwdb.ChallengeInstance_Available {
			instanceGreen++
		} else {
			instanceRed++
		}
	}
	instanceParts := []string{}
	if instanceGreen > 0 {
		instanceParts = append(instanceParts, fmt.Sprintf("%dx🟢", instanceGreen))
	}
	if instanceRed > 0 {
		instanceParts = append(instanceParts, fmt.Sprintf("%dx🔴", instanceRed))
	}
	stats := strings.Join(instanceParts, " + ")
	return stats
}

func asciiSubscriptionsStats(subscriptions []*pwdb.ChallengeSubscription) string {
	if len(subscriptions) == 0 {
		return "🚫"
	}

	subscriptionGreen := 0
	subscriptionRed := 0
	for _, subscription := range subscriptions {
		if subscription.Status == pwdb.ChallengeSubscription_Active {
			subscriptionGreen++
		} else {
			subscriptionRed++
		}
	}
	subscriptionParts := []string{}
	if subscriptionGreen > 0 {
		subscriptionParts = append(subscriptionParts, fmt.Sprintf("%dx🟢", subscriptionGreen))
	}
	if subscriptionRed > 0 {
		subscriptionParts = append(subscriptionParts, fmt.Sprintf("%dx🔴", subscriptionRed))
	}
	stats := strings.Join(subscriptionParts, " + ")
	return stats
}

func asciiBool(input bool) string {
	if !input {
		return "❌"
	}
	return "✅"
}

func asciiStatus(status string) string {
	switch status {
	case "Active":
		status += " 🟢"
	default:
		status += " 🔴"
	}
	return status
}
