package config

type Configure struct {
	RDB     *RDBConfig     `yaml:"RDB"`
	Redis   *RedisConfig   `yaml:"Redis"`
	Srv     *SrvConfig     `yaml:"Server"`
	Email   *EmailConfig   `yaml:"Email"`
	Contest *ContestConfig `yaml:"Contest"`
}

type RDBConfig struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
}

type RedisConfig struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	Password string `yaml:"Password"`
}

type SrvConfig struct {
	Domain string `yaml:"Domain"`
	Port   int    `yaml:"Port"`
}

type EmailConfig struct {
	SMTPAddr      string `yaml:"SMTPAddr"`
	SMTPPort      int    `yaml:"SMTPPort"`
	EmailAddr     string `yaml:"EmailAddr"`
	EmailPassword string `yaml:"EmailPassword"`
	EmailSign     string `yaml:"EmailSign"`
	EmailAlias    string `yaml:"EmailAlias"`
}

type ContestConfig struct {
	Name            string   `yaml:"Name"`
	StartTime       string   `yaml:"StartTime"`
	EndTime         string   `yaml:"EndTime"`
	Note            string   `yaml:"Note"`
	ValidTshirtSize []string `yaml:"ValidTshirtSize"`
	SchoolName      []string `yaml:"SchoolName"`
}
