package configs

type AppConfiguration struct {
	Mode       string `env:"GIN_MODE"`
	Port       int    `env:"PORT"`
	AppEnv     string `env:"APP_ENV"`
	Version    string `env:"VERSION"`
	Database   DBConfig
	AwsConf    AwsConfiguration
	GoogleAuth GoogleAuthConfig
	Cache      CacheConfig
	Swag       SwagConf
}

type DBConfig struct {
	Type           string `env:"DB_TYPE,default=postgres"`
	EndPoint       string `env:"DB_ENDPOINT"`
	ReadEndPoint   string `env:"READ_ENDPOINT"`
	Name           string `env:"DB_NAME,default=postgres"`
	User           string `env:"DB_USER,default=postgres"`
	Password       string `env:"DB_PASSWORD"`
	Migrate        bool   `env:"DB_MIGRATE,default=true"`
	MigrationsPath string `env:"DB_MIGRATIONS_PATH,default=db/migrations"`
}

type SwagConf struct {
	Host string `env:"SWAG_HOST,default=localhost"`
}

type CacheConfig struct {
	EndPoint string `env:"CACHE_ENDPOINT"`
	Port     string `env:"CACHE_PORT"`
	Password string `env:"CACHE_PASS"`
}

type AwsConfiguration struct {
	AwsProfile   string `env:"AWS_PROFILE"`
	AwsRegion    string `env:"AWS_REGION"`
	ClientId     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	UserPoolId   string `env:"USER_POOL_ID"`
}

type GoogleAuthConfig struct {
	ClientID     string `env:"CLIENT_ID_GOOGLE"`
	ClientSecret string `env:"CLIENT_SECRET_GOOGLE"`
	RedirectURL  string `env:"REDIRECT_URL_GOOGLE"`
}
