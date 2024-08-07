package main

import "fmt"

func main() {
	sl := []int{-5, -10}
	for i := 0; i < len(sl); i++ {
		sl = append(sl, i*10)
		fmt.Println(sl)
	}

	//cfg := config.Must()
	//zap.ReplaceGlobals(logger.Must(cfg.Logger))
	//if err := app.Run(cfg); err != nil {
	//	zap.S().Fatal(err)
	//}
}
