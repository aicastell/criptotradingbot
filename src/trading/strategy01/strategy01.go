package strategy01

import(
	"fmt"
    "time"
	"finantial/ema"
    "markets/generic"
    "markets/bitstamp"
)

var UNDEF = int32(-1)
var TRUE  = int32(1)
var FALSE = int32(0)

func Start(pair string, period time.Duration, training_iters int, win_len_min int, win_len_max int) {

    var ema_fast ema.TFinantial_EMA
    var ema_slow ema.TFinantial_EMA
    var ema_vol  ema.TFinantial_EMA
    var market generic.TMarket

    ema_fast.Reset(win_len_min)
    ema_slow.Reset(win_len_max)
    ema_vol.Reset(10)
    market.Reset("bitcoin", "euro", 1000)

    var fast_gt_slow = int32(UNDEF)

	iter := 0

    for {
        fmt.Println("Iter")
        price, err := bitstamp.DoGet(pair)
        if (err != nil) {
            fmt.Println("Error en el doget de bitstamp")
            time.Sleep(period * time.Second)
            continue
        }

        ema_fast.New_price(price)
        ema_slow.New_price(price)
        fmt.Println("price: ", price, "ema_fast: ", ema_fast.Ema(), "ema_slow: ", ema_slow.Ema(), "time: ", time.Now())

        time.Sleep(period * time.Second)

		if (iter < training_iters) {
		    iter++;
			continue
		}

		// End of training, start trading

        // Initialize fast_gt_slow only once after training
        if (fast_gt_slow  == UNDEF) {
            if (ema_fast.Ema() > ema_slow.Ema()) {
                fast_gt_slow = TRUE
            } else {
                fast_gt_slow = FALSE
            }
		    fmt.Println("Training ready. Starting trade now...")
            continue
        }

        // fast_gt_slow already defined

        /*
        if (market.InsideMarket()) {
            if (price < market.LastBuyPrice()) {
                market.DoSell(price)
                fmt.Println("********************************** Activated: CONTROL1")
                fmt.Println("********************************** VENDE a: ", market.LastSellPrice())
                fmt.Println("********************************** FIAT: ", market.Fiat())
            } else {
                fmt.Println("===> He comprado y esta subiendo, GOOD SIGNAL")
            }
        } */

        if (fast_gt_slow == FALSE) {
            if (ema_fast.Ema() < ema_slow.Ema()) {
                fmt.Println("ema_fast < ema_slow... Se mantiene la tendencia de bajada")
                // tendency is maintained (BAJADA)
                continue
            } else {
                if (market.InsideMarket() == false) {
                    market.DoBuy(price)
                    fmt.Println("********************************** COMPRA a: ", market.LastBuyPrice())
                    fmt.Println("********************************** CRYPTO: ", market.Crypto())
                } else {
                    fmt.Println("===> Tocaba comprar pero ya estoy dentro")
                }
                fast_gt_slow = TRUE
            }
        } else {
            if (ema_fast.Ema() > ema_slow.Ema()) {
                fmt.Println("ema_fast > ema_slow... Se mantiene la tendencia de subida")
                // tendency is maintained (SUBIDA)
                continue
            } else {
                if (market.InsideMarket() == true) {
                    market.DoSell(price)
                    fmt.Println("********************************** VENDE a: ", market.LastSellPrice())
                    fmt.Println("********************************** FIAT: ", market.Fiat())
                } else {
                    fmt.Println("===> Tocaba vender pero estoy fuera")

                }
                fast_gt_slow = FALSE
            }
        }
    }
}