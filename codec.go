package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"proxy/encrypt"
	"strings"
)

var buffSize = 1024

func EncryptPack(data []byte) []byte {
	transformed := encrypt.AesEncryptData(data)
	dataLen := uint16(len(transformed))
	if dataLen > 0 {
		var pkg = new(bytes.Buffer)
		err := binary.Write(pkg, binary.BigEndian, dataLen)
		errPrint(err)
		err = binary.Write(pkg, binary.BigEndian, transformed)
		errPrint(err)
		return pkg.Bytes()
	}
	return make([]byte, 0)
}

func transformIoEncrypt(dst io.Writer, src io.Reader) {
	for {
		buf := make([]byte, buffSize-encrypt.BlockSize)
		n, err := src.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}
		if n == 0 {
			continue
		}
		pack := EncryptPack(buf[:n])
		if len(pack) > 0 {
			_, ew := dst.Write(pack)
			if errPrint(ew) {
				break
			}
		}
	}
}

func DecryptUnpackOne(src io.Reader) ([]byte, error) {
	cnt := 0
	head := make([]byte, 2)
	var data []byte
	var packLen int
	for {
		//read length head
		if cnt < 2 {
			tmp := make([]byte, 2-cnt)
			n, err := src.Read(tmp)
			if err != nil {
				if err != io.EOF && !strings.Contains(err.Error(), "use of closed network connection") {
					errPrint(err)
				}
				return make([]byte, 0), err
			}
			for i := 0; i < n; i++ {
				head[cnt+i] = tmp[i]
			}
			cnt += n
			if cnt < 2 {
				continue
			}
		}
		var dataLen uint16
		err := binary.Read(bytes.NewBuffer(head), binary.BigEndian, &dataLen)
		if errPrint(err) || dataLen == 0 {
			return make([]byte, 0), err
		}
		packLen = int(dataLen) + 2
		tmp := make([]byte, packLen-cnt)
		n, err := src.Read(tmp)
		if errPrint(err) {
			return make([]byte, 0), err
		}
		cnt += n
		data = append(data, tmp[:n]...)
		if cnt < packLen {
			continue
		}
		transformed, err := encrypt.AesDecryptData(data)
		errPrint(err)
		return transformed, err
	}
}

func transformIoDecrypt(dst io.Writer, src io.Reader) {
	for {
		onePack, err := DecryptUnpackOne(src)
		if len(onePack) > 0 {
			_, ew := dst.Write(onePack)
			if errPrint(ew) {
				break
			}
		}
		if err != nil {
			break
		}
	}
}
