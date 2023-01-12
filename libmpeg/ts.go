package libmpeg

const (
	PSIPAT = 0x0000
	PSICAT = 0x0001
	//PSINIT = 0x0000
	PSIPMT  = 0x0002
	PSITSDT = 0x0002
)

type TSHeader struct {
	syncByte                   byte //1byte fixed 0x47
	transportErrorIndicator    byte //1bit
	payloadUnitStartIndicator  byte //1bit
	transportPriority          byte //1bit
	pid                        int  //13bits //0x0000-PAT/0x001-CAT/0x002-TSDT/0x0004-0x000F reserved/0x1FFF null packet
	transportScramblingControl byte //2bits
	adaptationFieldControl     byte //2bits 10/11/01/11
	continuityCounter          byte //4bits
}

func readTSHeader(data []byte) (TSHeader, int) {
	h := TSHeader{}
	h.syncByte = data[0]
	h.transportErrorIndicator = data[1] >> 7 & 0x1
	h.payloadUnitStartIndicator = data[1] >> 6 & 0x1
	h.transportPriority = data[1] >> 5 & 0x1
	h.pid = int(data[1]&0x1F) << 8
	h.pid = h.pid | int(data[2])
	h.transportScramblingControl = data[3] >> 6 & 0x3
	h.adaptationFieldControl = data[3] >> 4 & 0x3
	h.continuityCounter = data[3] & 0xF
	index := 4

	switch h.adaptationFieldControl {
	case 0x00:
		//discard
		break
	case 0x01:
		break
	case 0x02:
		//adaptation field only,no payload
		break
	case 0x03:
		//2.4.3.4 adaptation_field
		length := data[4]
		index++
		index += int(length)
		break
	}

	return h, index
}

func readTableHeader(data []byte) int {
	//pointerField := data[0]
	//tableId := data[1]
	//sectionSyntaxIndicator := data[2] >> 7 & 0x01
	////'0'
	////2bits reserved
	sectionLength := int(data[2]&0xF) << 8
	sectionLength |= int(data[3])
	//transportStreamId := (int(data[4]) << 8) | int(data[5])
	////2bits reserved
	//versionNumber := data[6] >> 1 & 0x1F
	////1bit current_next_indicator
	//sectionNumber := data[7]
	//lastSectionNumber := data[8]
	//println(pointerField)
	//println(tableId)
	//println(sectionSyntaxIndicator)
	//println(transportStreamId)
	//println(versionNumber)
	//println(sectionNumber)
	//println(lastSectionNumber)

	return sectionLength
}

func readPAT(data []byte) []int {
	sectionLength := readTableHeader(data)
	index := 9
	sectionLength -= 5
	var pmt []int
	for sectionLength >= 8 {
		programNumber := (int(data[index]) << 8) | int(data[index+1])
		//reserved 3bits
		index += 2
		pid := (int(data[index]&0x1F) << 8) | int(data[index+1])
		if programNumber == 0 {
			//network pid
		} else {
			//pat
			pmt = append(pmt, pid)
		}
		index += 2
		sectionLength -= 4
	}

	return pmt
}

func readPMT(data []byte) []int {
	sectionLength := readTableHeader(data)
	index := 9
	sectionLength -= 5
	pcrPid := int(data[index]&0x1F) << 8
	index++
	pcrPid |= int(data[index])
	index++
	programInfoLength := int(data[index]&0xF) << 8
	index++
	programInfoLength |= int(data[index])
	index++

	var pid []int
	for sectionLength >= 9 {
		//streamType
		_ = data[index]
		index++
		elementaryPid := int(data[index]&0x1F) << 8
		index++
		elementaryPid |= int(data[index])
		index++
		esInfoLength := int(data[index]&0x1F) << 8
		pid = append(pid, elementaryPid)
		index++
		esInfoLength |= int(data[index])
		index++
		sectionLength -= 5
	}

	return pid
}

func readCASection(data []byte) int {

	return 0
}
