package vivo

type Config struct {
	AppID  string
	AppKey string
}

type Vivo struct {
	appID  string
	appKey string
}

func NewVivoAIGC(config Config) *Vivo {
	return &Vivo{
		appID:  config.AppID,
		appKey: config.AppKey,
	}
}
