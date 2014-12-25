package main

/*
#include <linux/spi/spidev.h>
#include <sys/ioctl.h>
typedef struct spi_ioc_transfer SPI_IOC_TRANSFER;
const unsigned long SPI_IOC_MESSAGE_1=SPI_IOC_MESSAGE(1);

*/
import "C"

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"
	"unsafe"
)

func xfer(file *os.File, txb *[]byte) (*[]byte, error) {
	len := len(*txb)
	rxb := make([]byte, len)
	var tr C.SPI_IOC_TRANSFER
	tr.tx_buf = C.__u64(uintptr(unsafe.Pointer(&(*txb)[0])))
	tr.rx_buf = C.__u64(uintptr(unsafe.Pointer(&rxb[0])))
	tr.len = C.__u32(len)
	err := ioctl(file, uintptr(C.SPI_IOC_MESSAGE_1), uintptr(unsafe.Pointer(&tr)))
	return &rxb, err
}

func ioctl(file *os.File, request, argp uintptr) error {
	_, _, errorp := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), request, argp)
	return os.NewSyscallError("ioctl", errorp)
}

type Radio struct {
	file *os.File
}

func NewRadio(filename string) *Radio {
	r := new(Radio)
	file, err := os.OpenFile(filename, os.O_RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}
	r.file = file
	return r
}

func (r *Radio) setReg(reg byte, val byte) {
	txb := []byte{reg, val}
	xfer(r.file, &txb)
}

func (r *Radio) sendStrobe(reg byte) {
	txb := []byte{reg}
	xfer(r.file, &txb)
}

func (r *Radio) setupRX() {

	conf := []byte{
		0x29, //0x0000 IOCFG2 - GDO2 Output Pin Configuration
		0x2e, //0x0001 IOCFG1 - GDO1 Output Pin Configuration
		0x06, //0x0002 IOCFG0 - GDO0 Output Pin Configuration
		0x47, //0x0003 FIFOTHR - RX FIFO and TX FIFO Thresholds
		0xd3, //0x0004 SYNC1 - Sync Word, High Byte
		0x91, //0x0005 SYNC0 - Sync Word, Low Byte
		0xff, //0x0006 PKTLEN - Packet Length
		0x04, //0x0007 PKTCTRL1 - Packet Automation Control
		0x06, //0x0008 PKTCTRL0 - Packet Automation Control
		0x00, //0x0009 ADDR - Device Address
		0x00, //0x000a CHANNR - Channel Number
		0x06, //0x000b FSCTRL1 - Frequency Synthesizer Control
		0x00, //0x000c FSCTRL0 - Frequency Synthesizer Control
		0x10, //0x000d FREQ2 - Frequency Control Word, High Byte
		0xb0, //0x000e FREQ1 - Frequency Control Word, Middle Byte
		0x71, //0x000f FREQ0 - Frequency Control Word, Low Byte
		0xc8, //0x0010 MDMCFG4 - Modem Configuration
		0x93, //0x0011 MDMCFG3 - Modem Configuration
		0x30, //0x0012 MDMCFG2 - Modem Configuration
		0x22, //0x0013 MDMCFG1 - Modem Configuration
		0xf8, //0x0014 MDMCFG0 - Modem Configuration
		0x24, //0x0015 DEVIATN - Modem Deviation Setting
		0x07, //0x0016 MCSM2 - Main Radio Control State Machine Configuration
		0x30, //0x0017 MCSM1 - Main Radio Control State Machine Configuration
		0x18, //0x0018 MCSM0 - Main Radio Control State Machine Configuration
		0x16, //0x0019 FOCCFG - Frequency Offset Compensation Configuration
		0x6c, //0x001a BSCFG - Bit Synchronization Configuration
		0x43, //0x001b AGCCTRL2 - AGC Control
		0x40, //0x001c AGCCTRL1 - AGC Control
		0x91, //0x001d AGCCTRL0 - AGC Control
		0x87, //0x001e WOREVT1 - High Byte Event0 Timeout
		0x6b, //0x001f WOREVT0 - Low Byte Event0 Timeout
		0xfb, //0x0020 WORCTRL - Wake On Radio Control
		0x56, //0x0021 FREND1 - Front End RX Configuration
		0x11, //0x0022 FREND0 - Front End TX Configuration
		0xe9, //0x0023 FSCAL3 - Frequency Synthesizer Calibration
		0x2a, //0x0024 FSCAL2 - Frequency Synthesizer Calibration
		0x00, //0x0025 FSCAL1 - Frequency Synthesizer Calibration
		0x1f, //0x0026 FSCAL0 - Frequency Synthesizer Calibration
		0x41, //0x0027 RCCTRL1 - RC Oscillator Configuration
		0x00, //0x0028 RCCTRL0 - RC Oscillator Configuration
	}

	conf[0x001b] = 0x07 //AGCCTRL2 - AGC Control
	conf[0x001c] = 0x00 //AGCCTRL1 - AGC Control
	conf[0x001d] = 0x90 //AGCCTRL0 - AGC Control

	buf := append([]byte{0x40}, conf...)

	fmt.Printf("Conf: % x\n", buf)

	xfer(r.file, &buf)

}

func main() {
	fmt.Println("Start...")

	radio := NewRadio("/dev/spidev0.0")
	file := radio.file

	radio.sendStrobe(0x30) //Reset

	radio.setupRX()
	rxb,_ := xfer(file, &[]byte{0x7e, 0, 0xc0, 0, 0, 0, 0, 0, 0})

	fmt.Printf("OutBytes: % x\n", *rxb)
	radio.sendStrobe(0x34) //RX-Mode

	var counter rlCounter
	rawChan := make(chan rune, 100)
	go processRaw(rawChan)

	for {
		rxb, _ := xfer(file, &[]byte{0x3d})
		if (*rxb)[0] == 31 {
			rxb, _ := xfer(file, &[]byte{0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
			//fmt.Printf("OutBytes: % x\n", *rxb)
			a := (*rxb)[1:]
			for i := range a {
				for b := 0; b < 8; b++ {
					counter.count((a[i]>>uint(7-b))&1, rawChan)
				}
			}
		} else {
			time.Sleep(1 * time.Millisecond)
		}
	}

}

type rlCounter struct {
	current byte
	counter uint
}

func (this *rlCounter) count(val byte, ch chan rune) {
	if val == this.current {
		this.counter++
	} else {
		//fmt.Printf("%b %d\n", this.current, this.counter)
		switch {
		case this.counter >= 2 && this.counter <= 6:
			ch <- '1'
		case this.counter >= 9 && this.counter <= 13:
			ch <- '3'
		case this.counter > 50 && this.current == 0:
			ch <- 'S'
		default:
			ch <- 'X'
		}
		this.current = val
		this.counter = 1
	}
}

func processRaw(ch chan rune) {
	var symbols bytes.Buffer
	ok := false
	for {
		c := <-ch
		switch c {
		case '1', '3':
			if ok {
				symbols.WriteByte(byte(c))
			}
		case 'X':
			ok = false
			symbols.Reset()
		case 'S':
			fmt.Println(symbols.String())

			symbols.Reset()
			ok = true
		}
		//fmt.Printf("#%s\n", c)

	}

}
