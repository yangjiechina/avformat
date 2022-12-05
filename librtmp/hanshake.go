package librtmp

import "avformat/utils"

const (
	VERSION             = 3
	HandshakePacketSize = 1536
	DefaultPort         = 1935
	RTMPSDefaultPort    = 443
)

/**
0 1 2 3 4 5 6 7
+-+-+-+-+-+-+-+-+
| version |
+-+-+-+-+-+-+-+-+
C0 and S0 bits
*/

/**
 The C1 and S1 packets are 1536 octets long, consisting of the
 following fields:
0 1 2 3
0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
| time (4 bytes) |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
| zero (4 bytes) |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
| random bytes |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
| random bytes |
| (cont) |
| .... |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
C1 and S1 bits
*/

/**
0 1 2 3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 | time (4 bytes) |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 | time2 (4 bytes) |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 | random echo |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 | random echo |
 | (cont) |
 | .... |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 C2 and S2 bits
*/

//C0+C1一起发，收到S1发C2。

func writeHandshakeC0(dst []byte) {
	dst[0] = VERSION
}

func writeHandshakeC1(dst []byte, time int) {
	utils.WriteDWORD(dst, uint32(time))
	utils.WriteDWORD(dst, 0)
}

func writeHandshakeRandom(dst []byte) {

}
