package hackrf

import (
	"fmt"
	"testing"
	"time"
)

func TestHackRF(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	defer Exit()
	dev, err := Open()
	if err != nil {
		t.Fatal(err)
	}
	defer dev.Close()
	if ver, err := dev.Version(); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Version: %s", ver)
	}
	total := 0
	if err := dev.StartRX(func(buf []byte) error {
		total += len(buf)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	if err := dev.StopRX(); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%d bytes\n", total)
}
