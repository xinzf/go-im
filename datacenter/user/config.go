package user

var config struct {
	Name    string `required:"true"`
	Version string `required:"true"`

	Db struct {
		Host         string `required:"true"`
		User         string `required:"true"`
		Pswd         string `required:"true"`
		Name         string `required:"true"`
		MaxOpenConns int    `default:100`
		MaxIdleConns int    `default:100`
	}

	Redis struct {
		Host         string `required:"true"`
		Password     string `required:"true"`
		Db           int    `default:0`
		MaxOpenConns int    `default:10`
		MaxIdleConns int    `default:10`
	}

	Discovery struct {
		Addrs    string
		Ttl      int
		Interval int
	}
}
