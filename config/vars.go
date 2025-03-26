package config

type GlobalConfig struct {
	AppName   string       `yaml:"AppName"`
	MODE      string       `yaml:"Mode"` // dev或prod
	VERSION   string       `yaml:"Version"`
	Host      string       `yaml:"Host"`
	Port      string       `yaml:"Port"`
	Databases []Datasource `yaml:"Databases"`
	Caches    []Cache      `yaml:"Caches"`
	Minio     struct {
		Endpoint        string `yaml:"endpoint"`
		AccessKeyID     string `yaml:"accessKeyID"`
		SecretAccessKey string `yaml:"secretAccessKey"`
		UseSSL          bool   `yaml:"useSSL"`
	} `yaml:"minio"`
	Elasticsearch struct {
		Addresses []string `yaml:"addresses"`
	} `yaml:"elasticsearch"`
	Wechat struct {
		AppId     string `yaml:"appid"`
		AppSecret string `yaml:"appsecret"`
	} `yaml:"wechat"`
	MQ struct {
		Broker        string `yaml:"broker"`         // MQ 的类型（例如 "rabbitmq", "kafka"）
		Address       string `yaml:"address"`        // MQ 地址（例如 RabbitMQ 的主机地址）
		Port          string `yaml:"port"`           // MQ 端口
		Username      string `yaml:"username"`       // MQ 用户名
		Password      string `yaml:"password"`       // MQ 密码
		VirtualHost   string `yaml:"virtual_host"`   // RabbitMQ 的虚拟主机
		Exchange      string `yaml:"exchange"`       // 默认交换机
		ExchangeType  string `yaml:"exchange_type"`  // 交换机类型（如 "direct", "fanout", "topic"）
		QueueName     string `yaml:"queue_name"`     // 默认队列名称
		RetryAttempts int    `yaml:"retry_attempts"` // 重试次数
		RetryDelay    int    `yaml:"retry_delay"`    // 重试延迟（秒）
	} `yaml:"mq"`
	Jwt struct {
		//关键点：不要留secret，甚至是_secret也不行
		Mundo    string `yaml:"mundo"`
		Offercat string `yaml:"offercat"`
	} `yaml:"jwt"`
	Apmq struct {
		Url string `yaml:"url"`
	} `yaml:"apmq"`
}

type Datasource struct {
	Key      string `yaml:"Key"`
	Type     string `yaml:"Type"`
	IP       string `yaml:"Ip"`
	PORT     string `yaml:"Port"`
	USER     string `yaml:"User"`
	PASSWORD string `yaml:"Password"`
	DATABASE string `yaml:"Database"`
}

type Cache struct {
	Key      string `yaml:"Key"`
	Type     string `yaml:"Type"`
	IP       string `yaml:"Ip"`
	PORT     string `yaml:"Port"`
	PASSWORD string `yaml:"Password"`
	DB       int    `yaml:"Db"`
}

type Oss struct {
	AccessKeySecret string `yaml:"AccessKeySecret"`
	AccessKeyId     string `yaml:"AccessKeyId"`
	EndPoint        string `yaml:"EndPoint"`
	BucketName      string `yaml:"BucketName"`
	BaseURL         string `yaml:"BaseURL"`
	Path            string `yaml:"Path"`
	CallbackUrl     string `yaml:"CallbackUrl"`
	ExpireTime      int64  `yaml:"ExpireTime"`
}

type Mail struct {
	SMTP     string `yaml:"Smtp"`
	PORT     int    `yaml:"Port"`
	ACCOUNT  string `yaml:"Account"`
	PASSWORD string `yaml:"Password"`
}

type Cms struct {
	SecretId   string `yaml:"SecretId"`
	SecretKey  string `yaml:"SecretKey"`
	AppId      string `yaml:"AppId"`
	TemplateId string `yaml:"TemplateId"`
	Sign       string `yaml:"Sign"`
}
