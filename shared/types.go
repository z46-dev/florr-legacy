package shared

import "time"

type (
	AuthType            uint8
	DiscordHistoryEvent uint8

	Account struct {
		ID                string    `json:"id" gosqlite:"id,primary,unique"`
		Username          string    `json:"username" gosqlite:"username,unique"`
		CreationDate      time.Time `json:"creation_date" gosqlite:"creation_date"`
		LastSeen          time.Time `json:"last_seen" gosqlite:"last_seen"`
		UsernameChangedAt time.Time `json:"username_changed_at" gosqlite:"username_changed_at"`
		Disabled          bool      `json:"disabled" gosqlite:"disabled"`
		Discord           *Discord  `json:"discord,omitempty"`
	}

	AuthMethod struct {
		ID              int       `json:"id" gosqlite:"id,primary,unique,increment"`
		AccountUsername string    `json:"accountUsername" gosqlite:"accountUsername"`
		Type            AuthType  `json:"type" gosqlite:"type"`
		Identifier      string    `json:"identifier" gosqlite:"identifier"`
		SecretHash      string    `json:"-" gosqlite:"secretHash"`
		CreatedAt       time.Time `json:"createdAt" gosqlite:"createdAt"`
	}

	Discord struct {
		ID       string `json:"id" gosqlite:"id,primary,unique"`
		Username string `json:"username" gosqlite:"username"`
		Avatar   string `json:"avatar" gosqlite:"avatar"`
	}

	DiscordHistory struct {
		ID              int                 `json:"id" gosqlite:"id,primary,unique,increment"`
		DiscordID       string              `json:"discordId" gosqlite:"discordId"`
		AccountUsername string              `json:"accountUsername" gosqlite:"accountUsername"`
		Event           DiscordHistoryEvent `json:"event" gosqlite:"event"`
		CreatedAt       time.Time           `json:"createdAt" gosqlite:"createdAt"`
	}
)
