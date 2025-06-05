package pkg

type configutation struct {
	PORT      uint
	JWT_TOKEN string
	DB_URL    string

	GIT_BRANCH string
	GIT_REMOTE string
	GIT_RESET  bool
}

// Config returns the configuration of whole app
func Config() *configutation {
	return &configutation{
		PORT:      uint(GetEnvInt("PORT", 8080)),
		JWT_TOKEN: GetEnv("JWT_SECRET", "some_secret_token"),
		DB_URL:    GetEnv("DB_URL", "file:app.db?_fk=1"),

		GIT_BRANCH: "master",
		GIT_REMOTE: "origin",
		GIT_RESET:  false,
	}
}
