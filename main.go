package main

import(
	"fmt"
	"os"
	"time"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/eiannone/keyboard"
)
var ctrl *beep.Ctrl
var streamer beep.StreamSeekCloser
var format beep.Format
//func metadata(path string)string{
//	f,err := os.Open(path)
//}
func player(path string , ready chan struct{}){
	f,err := os.Open(path)
	if err != nil{
		panic(err)
	}
	defer f.Close()
	var err2 error
	streamer,format,err2 = mp3.Decode(f)
	if err2 != nil {
		panic(err)
	}
	defer streamer.Close()
	speaker.Init(format.SampleRate,format.SampleRate.N(time.Second/10))
	ctrl = &beep.Ctrl{Streamer: streamer,Paused:false}
	speaker.Play(ctrl)
	close(ready)
	select{}
}
//func readPlaylist(path string){
//	f,err := os.Open(path)
//}
func printTime() {
	for {
		pos := streamer.Position()
		seconds := time.Duration(float64(pos) / float64(format.SampleRate)) * time.Second
		fmt.Printf("\r%s", seconds.Truncate(time.Second))
		time.Sleep(1 * time.Second)
	}
}
func main(){
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()
	ready := make(chan struct{})
	go player("test.mp3",ready)
	<-ready
	go printTime()
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		if key == keyboard.KeySpace{
			speaker.Lock()
			ctrl.Paused = !ctrl.Paused
			speaker.Unlock()
		}else if char == 'q'{
			break
		}
	}
}
