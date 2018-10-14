package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/cheggaaa/pb"
	"github.com/urfave/cli"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/neurosky"
)

func main() {
	app := cli.NewApp()
	app.Name = "meisou"
	app.Version = "1.0.0"
	app.Usage = "瞑想力(0-100)が一定値を超えたらコマンドが実行されます。"

	var border float64

	app.Flags = []cli.Flag{
		cli.Float64Flag{
			Name:        "border, b",
			Value:       60,
			Usage:       "設定された瞑想力を超えると指定されたコマンドが実行される",
			Destination: &border,
		},
	}

	app.Action = func(c *cli.Context) error {
		command := c.Args().Get(0)
		fmt.Println("目標瞑想力:" + strconv.Itoa(int(border)) + "%")
		if command != "" {
			fmt.Println("実行コマンド:" + command)
		}
		meditation := 0.0
		basePart := 0.9
		currentPart := 1.0 - basePart

		bar := pb.New(int(border))
		bar.ShowTimeLeft = false
		bar.ShowSpeed = false

		adaptor := neurosky.NewAdaptor("/dev/tty.MindWaveMobile-SerialPo")
		neuro := neurosky.NewDriver(adaptor)
		var robot *gobot.Robot

		work := func() {
			neuro.On(neuro.Event("meditation"), func(data interface{}) {
				if currentMediation, ok := data.(uint8); ok {
					meditation = (meditation * basePart) + (float64(currentMediation) * currentPart)
					if int(meditation) != 0 {
						bar.Start()
					}
					bar.Set(int(meditation))
					if meditation > border {
						bar.Finish()
						fmt.Println("目標達成めっちゃ集中してる！！！")
						if err := robot.Stop(); err == nil {
							if command != "" {
								exec.Command(command)
							}
							os.Exit(1)
						}
					}
				}
			})
		}

		robot = gobot.NewRobot("brainBot",
			[]gobot.Connection{adaptor},
			[]gobot.Device{neuro},
			work,
		)

		robot.Start()

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
