package config

type Configuration struct {
	WebServer struct {
		Address            string   `toml:"address" default:":6830" validate:"required"` // Listen address for the web application server e.g. ":6830" or "0.0.0.0:6830"
		TLSDir             string   `toml:"tls_dir" default:""`                          // Directory containing a crt and a key file for TLS. Leave empty to use HTTP instead of HTTPS.
		EnableCSRF         bool     `toml:"enable_csrf" default:"true"`                  // Enable CSRF protection middleware
		CORSAllowedOrigins []string `toml:"cors_allowed_origins" default:"[\"*\"]"`      // CORS allowed origins for the web server
	} `toml:"web_server"` // Web server configuration

	Database struct {
		File string `toml:"file" default:"florr-legacy.db" validate:"required"` // Path to the MySQL database file
	} `toml:"database"` // Database configuration

	CommandServer struct {
		Address string `toml:"address" default:":6831" validate:"required"` // Listen address for the command server e.g. ":6831" or "0.0.0.0":6831"
	} `toml:"command_server"` // Command server configuration

	Discord struct {
		ClientID       string `toml:"client_id" default:"" validate:"required"`                                                         // Discord OAuth2 Client ID
		ClientSecret   string `toml:"client_secret" default:"" validate:"required"`                                                     // Discord OAuth2 Client Secret
		ClientRedirect string `toml:"client_redirect" default:"http://localhost:6800/api/v1/auth/discord/callback" validate:"required"` // Discord OAuth2 Redirect URI
	} `toml:"discord"` // Discord OAuth2 configuration

	Auth struct {
		SessionTTLHours           int    `toml:"session_ttl_hours" default:"720"`                                                          // Session token lifetime in hours (default 30 days)
		ChallengeTTLSeconds       int    `toml:"challenge_ttl_seconds" default:"300"`                                                      // Seconds before a challenge expires
		PowDifficulty             int    `toml:"pow_difficulty" default:"3"`                                                               // Leading hex zeros required for PoW challenges
		PowDynamic                bool   `toml:"pow_dynamic" default:"true"`                                                               // Enable dynamic PoW difficulty
		PowBurstWindowSeconds     int    `toml:"pow_burst_window_seconds" default:"300"`                                                   // Seconds of PoW request history to consider
		PowBurstStep              int    `toml:"pow_burst_step" default:"3"`                                                               // Requests per window step to increase PoW difficulty
		PowMaxPenalty             int    `toml:"pow_max_penalty" default:"4"`                                                              // Max extra difficulty added to PoW
		PasswordMinLength         int    `toml:"password_min_length" default:"8"`                                                          // Minimum password length
		UsernameMinLength         int    `toml:"username_min_length" default:"3"`                                                          // Minimum username length
		UsernameMaxLength         int    `toml:"username_max_length" default:"16"`                                                         // Maximum username length
		PasswordSalt              string `toml:"password_salt" default:"" validate:"required"`                                             // Salt prefix for password hashing
		TurnstileEnabled          bool   `toml:"turnstile_enabled" default:"false"`                                                        // Enable Cloudflare Turnstile challenges
		TurnstileSiteKey          string `toml:"turnstile_site_key" default:""`                                                            // Cloudflare Turnstile site key (client)
		TurnstileSecretKey        string `toml:"turnstile_secret_key" default:""`                                                          // Cloudflare Turnstile secret key (server)
		TurnstileAllowedHostnames string `toml:"turnstile_allowed_hostnames" default:""`                                                   // Comma-separated hostnames allowed in successful Turnstile siteverify responses
		TurnstileVerifyURL        string `toml:"turnstile_verify_url" default:"https://challenges.cloudflare.com/turnstile/v0/siteverify"` // Turnstile verification endpoint
	} `toml:"auth"` // Authentication configuration
}

var Config Configuration
