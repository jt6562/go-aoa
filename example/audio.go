package main

import (
	aoa "github.com/jt6562/go-aoa"
)

func main() {
	acc := aoa.NewAccessory()
	defer acc.Close()
	acc.SwitchToAccessoryMode(aoa.MODE_AUDIO)

}
