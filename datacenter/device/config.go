package device

var config struct {
	Name    string `required:"true"`
	Version string `required:"true"`

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
