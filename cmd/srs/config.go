package main

type Config struct {
	Port        int    `envconfig:"PORT" required:"true"`
	OpenAIToken string `envconfig:"OPENAI_TOKEN" required:"true"`
	Auth        Auth
}

type Auth struct {
	ClientID    string `envconfig:"AUTH_CLIENT_ID" required:"true"`
	CallbackURL string `envconfig:"AUTH_CALLBACK_URL" required:"true"`
	SessionKey  string `envconfig:"AUTH_SESSION_KEY" required:"true"`
}
