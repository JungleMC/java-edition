package startup

import (
	"context"
	"github.com/JungleMC/java-edition/internal/config"
	"github.com/JungleMC/java-edition/internal/net"
	"github.com/caarlos0/env"
	"github.com/go-redis/redis/v8"
)

func Start(rdb *redis.Client) {
	PopulateDummyData(rdb)

	config.Get = &config.Config{}
	if err := env.Parse(config.Get); err != nil {
		panic(err)
	}

	_, err := net.Bootstrap(rdb, config.Get.ListenAddress, config.Get.ListenPort, config.Get.OnlineMode)
	if err != nil {
		panic(err)
	}
}

func PopulateDummyData(rdb *redis.Client) {
	favicon := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAEAAAABACAMAAACdt4HsAAABHVBMVEVnvUfS57h0SSHI4615xFrQ5rZsv0xovkjM5bHF4qrK5LCKy2xqvkqh1IS53Z2e04GY0Xp0wlTB4KXA4KWn14uj1YeWz3nN5bOp142Ax2B2w1fR5baw2pSNdEnO5bR9xl9vwE+KbUOBXTXG46y73qCz25eOzXCej2NywVJwwVFnuEVpskSHaD+FZjx2TiTJ2Km2t4uzsoWilmp7xVyReVBsrkZ9WTDN37C+36Ks2I+SznWFyWaWgliDYjmBYDd6Uyt1TCPF0qPBy56+xZey2pa6vpGqpHimnnKMzG6CyGSaiFx4wlmMk1WUf1RwskqNckhprEFppT5siDVzai633Ju0yI+lzoSvrH+npnactnKZsG17uliGiExtuUlvmz8RFckjAAADU0lEQVRYw+2X12LaQBBF9wrUEEXEdDAGG2NjB1zj3mviHju9/f9nhC2KJVaLZD/nvMADOuzM7syOyH/UmFW3krUtXbfsbMWtmuRFpJ2zMgKUG06axGU+YyEEKzNP4mAUdSjQMwaJ5H0CY0g4UcEXEcH22EXkJxDJRJ4oSc4gBok15fMWYpFKKtZvg2GXkq1aknH9bsg1/15r5btbYMyERmGI+DvGmBzPijyE/aYCzhQZQw2cSsj+g2ONP/U2OI4UQMLLMRnLlvez0SCKiCd4A0FmpH70gMDsFndHgi/1goJmj/jJICDoAqgSyuIif74M7AYEyAW2xwoKPgEo0C+PmvZIP3MAskFByt8fHAQFWU8w0DRtQaxw4p9A3oizEcG2WOKCNoQGUQIwOyJo+Pqfr3/NkCENL0ZNrMD1BB1flzOJRxXP6HSD6aHYGn7eUcGdiKlM/4olSzqzLnwU0qTLTB+IcUMFv1pkiu1yiZiBflOSyoBjJcCxm3tUcK/bXn5S8FMgHlkouKWCNhRkpQqRWKeCPhTYxCMFBfdUsK/sTMRDh4I2FUxDgR4t6FPBYb1eVwgiQziggoEyhOgkbmiUuiqJ8jYmuvPVanWtVsv3XLqsQybYjNzGSliNFwAcMcGSQlCQj7JL8kneQFgLqGuMY3quTVZ1u5b/KMvFVCQdmzw8LRh6gbaATS44oZUo6MjF5K8x1+wkzD/aYgtFKljiglMWrknMnTTvB/IF0IB3q+lwb7QfFaRKTeCYC86pYCdlrOk5RUMhDnzMaSvgnHDBMhPYv3+iqGppaX9ulrWv4JxywdzKysr3b/2nzyjPQm6qoq0/sznYb4JxzgXTb4e050Z2M0f89PzlsKpNLom1MC4QQrNFiGoJ9T1t0J4c0ueCS4SQUVyunOX1/ekhG1ywCpkZaUJwIHMhkggZJ2zAkLjkgitFGShGnEA2GXvyqJcOHZFthWBdGvTyijFvtDPNccGkasyTDSNruOKCtjRoxh11b7mgL426cYftdS44iDdsCxzfiZrkgg3f+XFIJEam+XyvMA4hKGcMEodeTlT3Wy44EsnPtV7w0tWwxL0iLgZLvHTFx5xyCxtC8KU0ZZLXMBCCRfJKjoTg4bUCTfCRvI6FOIK/RRpB9J2EXRgAAAAASUVORK5CYII="

	rdb.Set(context.Background(), "config:description", "A JungleTree Server", 0)
	rdb.Set(context.Background(), "config:favicon", favicon, 0)
	rdb.Set(context.Background(), "config:max_players", 20, 0)
}
