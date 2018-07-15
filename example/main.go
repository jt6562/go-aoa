package main

import (
	aoa "github.com/jt6562/go-aoa"
)

func main() {
	acc := aoa.NewAccessory()
	defer acc.Close()
	acc.SwitchToAccessoryMode(aoa.MODE_AUDIO)

	// err := acc.SwitchToAccessoryMode(aoa.MODE_ACCESSORY) //aoa.MODE_ACCESSORY | aoa.MODE_AUDIO
	// rw, err := acc.OpenAcessoryInterface()
	// fmt.Println(rw, err)

}
