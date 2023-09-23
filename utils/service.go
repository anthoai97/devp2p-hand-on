package utils

type PingService struct {
}

type PingResult struct {
	Ping string
	Pong string
}

func (p *PingService) Ping(ping string) PingResult {
	return PingResult{
		Ping: ping,
		Pong: "Pong",
	}
}
