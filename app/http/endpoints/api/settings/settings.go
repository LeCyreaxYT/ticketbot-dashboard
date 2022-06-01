package api

import (
	"context"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type Settings struct {
	database.Settings
	Prefix            string                `json:"prefix"`
	WelcomeMessage    string                `json:"welcome_message"`
	TicketLimit       uint8                 `json:"ticket_limit"`
	Category          uint64                `json:"category,string"`
	ArchiveChannel    *uint64               `json:"archive_channel,string"`
	NamingScheme      database.NamingScheme `json:"naming_scheme"`
	PingEveryone      bool                  `json:"ping_everyone"`
	UsersCanClose     bool                  `json:"users_can_close"`
	CloseConfirmation bool                  `json:"close_confirmation"`
	FeedbackEnabled   bool                  `json:"feedback_enabled"`
}

func GetSettingsHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	var settings Settings

	group, _ := errgroup.WithContext(context.Background())

	group.Go(func() (err error) {
		settings.Settings, err = dbclient.Client.Settings.Get(guildId)
		return
	})

	// prefix
	group.Go(func() (err error) {
		settings.Prefix, err = dbclient.Client.Prefix.Get(guildId)
		if err == nil && settings.Prefix == "" {
			settings.Prefix = "t!"
		}

		return
	})

	// welcome message
	group.Go(func() (err error) {
		settings.WelcomeMessage, err = dbclient.Client.WelcomeMessages.Get(guildId)
		if err == nil && settings.WelcomeMessage == "" {
			settings.WelcomeMessage = "Thank you for contacting support.\nPlease describe your issue and await a response."
		}

		return
	})

	// ticket limit
	group.Go(func() (err error) {
		settings.TicketLimit, err = dbclient.Client.TicketLimit.Get(guildId)
		if err == nil && settings.TicketLimit == 0 {
			settings.TicketLimit = 5 // Set default
		}

		return
	})

	// category
	group.Go(func() (err error) {
		settings.Category, err = dbclient.Client.ChannelCategory.Get(guildId)
		return
	})

	// archive channel
	group.Go(func() (err error) {
		settings.ArchiveChannel, err = dbclient.Client.ArchiveChannel.Get(guildId)
		return
	})

	// allow users to close
	group.Go(func() (err error) {
		settings.UsersCanClose, err = dbclient.Client.UsersCanClose.Get(guildId)
		return
	})

	// ping everyone
	group.Go(func() (err error) {
		settings.PingEveryone, err = dbclient.Client.PingEveryone.Get(guildId)
		return
	})

	// naming scheme
	group.Go(func() (err error) {
		settings.NamingScheme, err = dbclient.Client.NamingScheme.Get(guildId)
		return
	})

	// close confirmation
	group.Go(func() (err error) {
		settings.CloseConfirmation, err = dbclient.Client.CloseConfirmation.Get(guildId)
		return
	})

	// close confirmation
	group.Go(func() (err error) {
		settings.FeedbackEnabled, err = dbclient.Client.FeedbackEnabled.Get(guildId)
		return
	})

	if err := group.Wait(); err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(200, settings)
}
