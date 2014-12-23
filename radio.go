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

func xfer(file *os.File, txb []byte) ([]byte, error) {
	len := len(txb)
	rxb := make([]byte, len)
	var tr C.SPI_IOC_TRANSFER
	tr.tx_buf = C.__u64(uintptr(unsafe.Pointer(&txb[0])))
	tr.rx_buf = C.__u64(uintptr(unsafe.Pointer(&rxb[0])))
	tr.len = C.__u32(len)
	err := ioctl(file, uintptr(C.SPI_IOC_MESSAGE_1), uintptr(unsafe.Pointer(&tr)))
	return rxb, err
}

func ioctl(file *os.File, request, argp uintptr) error {
	_, _, errorp := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), request, argp)
	return os.NewSyscallError("ioctl", errorp)
}

func setReg(file *os.File, reg, val byte) {
	txb := []byte{reg, val}
	xfer(file, txb)
}

func sendStrobe(file *os.File, reg byte) {
	txb := []byte{reg}
	xfer(file, txb)
}

func main() {
	fmt.Println("Start...")
	file, err := os.OpenFile("/dev/spidev0.0", os.O_RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}

	sendStrobe(file, 0x30) //Reset




//
// Rf settings for CC1101
//
setReg(file, 0x0000, 0x29) //IOCFG2 - GDO2 Output Pin Configuration
setReg(file, 0x0001, 0x2e) //IOCFG1 - GDO1 Output Pin Configuration
setReg(file, 0x0002, 0x06) //IOCFG0 - GDO0 Output Pin Configuration
setReg(file, 0x0003, 0x47) //FIFOTHR - RX FIFO and TX FIFO Thresholds
setReg(file, 0x0004, 0xd3) //SYNC1 - Sync Word, High Byte
setReg(file, 0x0005, 0x91) //SYNC0 - Sync Word, Low Byte
setReg(file, 0x0006, 0xff) //PKTLEN - Packet Length
setReg(file, 0x0007, 0x04) //PKTCTRL1 - Packet Automation Control
setReg(file, 0x0008, 0x06) //PKTCTRL0 - Packet Automation Control
setReg(file, 0x0009, 0x00) //ADDR - Device Address
setReg(file, 0x000a, 0x00) //CHANNR - Channel Number
setReg(file, 0x000b, 0x06) //FSCTRL1 - Frequency Synthesizer Control
setReg(file, 0x000c, 0x00) //FSCTRL0 - Frequency Synthesizer Control
setReg(file, 0x000d, 0x10) //FREQ2 - Frequency Control Word, High Byte
setReg(file, 0x000e, 0xb0) //FREQ1 - Frequency Control Word, Middle Byte
setReg(file, 0x000f, 0x71) //FREQ0 - Frequency Control Word, Low Byte
setReg(file, 0x0010, 0xc8) //MDMCFG4 - Modem Configuration
setReg(file, 0x0011, 0x93) //MDMCFG3 - Modem Configuration
setReg(file, 0x0012, 0x30) //MDMCFG2 - Modem Configuration
setReg(file, 0x0013, 0x22) //MDMCFG1 - Modem Configuration
setReg(file, 0x0014, 0xf8) //MDMCFG0 - Modem Configuration
setReg(file, 0x0015, 0x24) //DEVIATN - Modem Deviation Setting
setReg(file, 0x0016, 0x07) //MCSM2 - Main Radio Control State Machine Configuration
setReg(file, 0x0017, 0x30) //MCSM1 - Main Radio Control State Machine Configuration
setReg(file, 0x0018, 0x18) //MCSM0 - Main Radio Control State Machine Configuration
setReg(file, 0x0019, 0x16) //FOCCFG - Frequency Offset Compensation Configuration
setReg(file, 0x001a, 0x6c) //BSCFG - Bit Synchronization Configuration
setReg(file, 0x001b, 0x43) //AGCCTRL2 - AGC Control
setReg(file, 0x001c, 0x40) //AGCCTRL1 - AGC Control
setReg(file, 0x001d, 0x91) //AGCCTRL0 - AGC Control
setReg(file, 0x001e, 0x87) //WOREVT1 - High Byte Event0 Timeout
setReg(file, 0x001f, 0x6b) //WOREVT0 - Low Byte Event0 Timeout
setReg(file, 0x0020, 0xfb) //WORCTRL - Wake On Radio Control
setReg(file, 0x0021, 0x56) //FREND1 - Front End RX Configuration
setReg(file, 0x0022, 0x11) //FREND0 - Front End TX Configuration
setReg(file, 0x0023, 0xe9) //FSCAL3 - Frequency Synthesizer Calibration
setReg(file, 0x0024, 0x2a) //FSCAL2 - Frequency Synthesizer Calibration
setReg(file, 0x0025, 0x00) //FSCAL1 - Frequency Synthesizer Calibration
setReg(file, 0x0026, 0x1f) //FSCAL0 - Frequency Synthesizer Calibration
setReg(file, 0x0027, 0x41) //RCCTRL1 - RC Oscillator Configuration
setReg(file, 0x0028, 0x00) //RCCTRL0 - RC Oscillator Configuration




	setReg(file, 0x001b, 0x07) //AGCCTRL2 - AGC Control
	setReg(file, 0x001c, 0x00) //AGCCTRL1 - AGC Control
	setReg(file, 0x001d, 0x90) //AGCCTRL0 - AGC Control

	xfer(file, []byte{0x7e, 0, 0xc0, 0, 0, 0, 0, 0, 0})

	sendStrobe(file, 0x34) //RX-Mode

	var counter rlCounter
	rawChan := make(chan rune, 100)
	go processRaw(rawChan)

	for {
		rxb, _ := xfer(file, []byte{0xc0})
		if rxb[0] == 31 {
			rxb, _ := xfer(file, []byte{0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				//fmt.Printf("OutBytes: % x\n", rxb)
			a := rxb[1:]
			for i := range a {
				for b := 0; b < 8; b++ {
					counter.count((a[i]>>uint(7-b))&1, rawChan)
				}
			}
		} else {
			time.Sleep(10 * time.Millisecond)
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

