package main

import (
	"image/color"
	"io"
	"machine"
	"net"
	"runtime"
	"time"

	"tinygo.org/x/drivers/examples/ili9341/initdisplay"
	"tinygo.org/x/drivers/ili9341"
	"tinygo.org/x/drivers/netlink"
	"tinygo.org/x/drivers/netlink/probe"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

var (
	ssid string
	pass string
)

const NTP_PACKET_SIZE = 48

const seventyYears = 2208988800

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

var (
	black = color.RGBA{0, 0, 0, 255}
	white = color.RGBA{255, 255, 255, 255}
)

func main() {
	display := initdisplay.InitDisplay()

	width, height := display.Size()
	if width < 320 || height < 240 {
		display.SetRotation(ili9341.Rotation270)
	}

	display.FillScreen(black)

	for !machine.Serial.DTR() {
		time.Sleep(100 * time.Millisecond)
	}

	link, _ := probe.Probe()

	if err := link.NetConnect(&netlink.ConnectParams{
		Ssid:       ssid,
		Passphrase: pass,
	}); err != nil {
		panic(err)
	}

	conn, err := net.Dial("udp", "ntp.nict.jp:123")
	if err != nil {
		panic(err)
	}

	var req = [48]byte{
		0xe3,
	}

	if _, err := conn.Write(req[:]); err != nil {
		panic(err)
	}

	res := make([]byte, NTP_PACKET_SIZE)

	n, err := conn.Read(res)
	if err != nil && err != io.EOF {
		panic(err)
	}

	if n != NTP_PACKET_SIZE {
		panic("short read")
	}

	t := uint32(res[40])<<24 | uint32(res[41])<<16 | uint32(res[42])<<8 | uint32(res[43])

	tm := time.Unix(int64(t-seventyYears), 0)

	conn.Close()

	link.NetDisconnect()

	runtime.AdjustTimeOffset(-1 * int64(time.Since(tm)))

	for {
		now := time.Now().In(jst).Format("03:04:05")

		_, o := tinyfont.LineWidth(&freemono.Regular24pt7b, now)

		x := (320 - o) / 2

		tinyfont.WriteLine(display, &freemono.Regular24pt7b, int16(x), 130, now, white)

		time.Sleep(time.Second)

		display.FillScreen(black)
	}
}
