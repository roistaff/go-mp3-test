package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)
func main() {
	f, err := os.Open("test.mp3")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		panic(err)
	}
	defer streamer.Close()

	// スピーカー初期化
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}

	// 再生開始
	speaker.Play(ctrl)

	fmt.Println("再生開始：p = 再生/一時停止、s = 10秒スキップ、q = 終了")

	// 入力処理（非同期）
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			cmd := scanner.Text()
			switch cmd {
			case "p":
				speaker.Lock()
				ctrl.Paused = !ctrl.Paused
				state := "再生"
				if ctrl.Paused {
					state = "一時停止"
				}
				fmt.Println(state)
				speaker.Unlock()
			case "s":
				speaker.Lock()
				pos := streamer.Position()
				skip := format.SampleRate.N(10 * time.Second)
				newPos := pos + skip
				if newPos < streamer.Len() {
					streamer.Seek(newPos)
					fmt.Println("10秒スキップ")
				} else {
					fmt.Println("スキップ範囲外")
				}
				speaker.Unlock()
			case "q":
				fmt.Println("終了します。")
				os.Exit(0)
			default:
				fmt.Println("無効なコマンドです（p = 再生/一時停止、s = スキップ、q = 終了）")
			}
		}
	}()

	// メインスレッドを終了させない
	select {}
}

