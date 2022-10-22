package utils

type AVCodecID int

const (
	//lint:ignore U1000 Ignore unused constant temporarily for debugging
	AVCodecIdNONE            = AVCodecID(0)
	AVCodecIdMPEG1VIDEO      = AVCodecID(1)
	AVCodecIdMPEG2VIDEO      = AVCodecID(2) ///< preferred ID for MPEG-1/2 video decodn
	AVCodecIdH261            = AVCodecID(3)
	AVCodecIdH263            = AVCodecID(4)
	AVCodecIdRV10            = AVCodecID(5)
	AVCodecIdRV20            = AVCodecID(6)
	AVCodecIdMJPEG           = AVCodecID(7)
	AVCodecIdMJPEGB          = AVCodecID(8)
	AVCodecIdLJPEG           = AVCodecID(9)
	AVCodecIdSP5X            = AVCodecID(0xa)
	AVCodecIdJPEGLS          = AVCodecID(0xb)
	AVCodecIdMPEG4           = AVCodecID(0xc)
	AVCodecIdRAWVIDEO        = AVCodecID(0xd)
	AVCodecIdMSMPEG4V1       = AVCodecID(0xe)
	AVCodecIdMSMPEG4V2       = AVCodecID(0xf)
	AVCodecIdMSMPEG4V3       = AVCodecID(0x10)
	AVCodecIdWMV1            = AVCodecID(0x11)
	AVCodecIdWMV2            = AVCodecID(0x12)
	AVCodecIdH263P           = AVCodecID(0x13)
	AVCodecIdH263I           = AVCodecID(0x14)
	AVCodecIdFLV1            = AVCodecID(0x15)
	AVCodecIdSVQ1            = AVCodecID(0x16)
	AVCodecIdSVQ3            = AVCodecID(0x17)
	AVCodecIdDVVIDEO         = AVCodecID(0x18)
	AVCodecIdHUFFYUV         = AVCodecID(0x19)
	AVCodecIdCYUV            = AVCodecID(0x1a)
	AVCodecIdH264            = AVCodecID(0x1b)
	AVCodecIdINDEO3          = AVCodecID(0x1c)
	AVCodecIdVP3             = AVCodecID(0x1d)
	AVCodecIdTHEORA          = AVCodecID(0x1e)
	AVCodecIdASV1            = AVCodecID(0x1f)
	AVCodecIdASV2            = AVCodecID(0x20)
	AVCodecIdFFV1            = AVCodecID(0x21)
	AVCodecId4XM             = AVCodecID(0x22)
	AVCodecIdVCR1            = AVCodecID(0x23)
	AVCodecIdCLJR            = AVCodecID(0x24)
	AVCodecIdMDEC            = AVCodecID(0x25)
	AVCodecIdROQ             = AVCodecID(0x26)
	AVCodecIdINTERPLAYVIDEO  = AVCodecID(0x27)
	AVCodecIdXANWC3          = AVCodecID(0x28)
	AVCodecIdXANWC4          = AVCodecID(0x29)
	AVCodecIdRPZA            = AVCodecID(0x2a)
	AVCodecIdCINEPAK         = AVCodecID(0x2b)
	AVCodecIdWSVQA           = AVCodecID(0x2c)
	AVCodecIdMSRLE           = AVCodecID(0x2d)
	AVCodecIdMSVIDEO1        = AVCodecID(0x2e)
	AVCodecIdIDCIN           = AVCodecID(0x2f)
	AVCodecId8BPS            = AVCodecID(0x30)
	AVCodecIdSMC             = AVCodecID(0x31)
	AVCodecIdFLIC            = AVCodecID(0x32)
	AVCodecIdTRUEMOTION1     = AVCodecID(0x33)
	AVCodecIdVMDVIDEO        = AVCodecID(0x34)
	AVCodecIdMSZH            = AVCodecID(0x35)
	AVCodecIdZLIB            = AVCodecID(0x36)
	AVCodecIdQTRLE           = AVCodecID(0x37)
	AVCodecIdTSCC            = AVCodecID(0x38)
	AVCodecIdULTI            = AVCodecID(0x39)
	AVCodecIdQDRAW           = AVCodecID(0x3a)
	AVCodecIdVIXL            = AVCodecID(0x3b)
	AVCodecIdQPEG            = AVCodecID(0x3c)
	AVCodecIdPNG             = AVCodecID(0x3d)
	AVCodecIdPPM             = AVCodecID(0x3e)
	AVCodecIdPBM             = AVCodecID(0x3f)
	AVCodecIdPGM             = AVCodecID(0x40)
	AVCodecIdPGMYUV          = AVCodecID(0x41)
	AVCodecIdPAM             = AVCodecID(0x42)
	AVCodecIdFFVHUFF         = AVCodecID(0x43)
	AVCodecIdRV30            = AVCodecID(0x44)
	AVCodecIdRV40            = AVCodecID(0x45)
	AVCodecIdVC1             = AVCodecID(0x46)
	AVCodecIdWMV3            = AVCodecID(0x47)
	AVCodecIdLOCO            = AVCodecID(0x48)
	AVCodecIdWNV1            = AVCodecID(0x49)
	AVCodecIdAASC            = AVCodecID(0x4a)
	AVCodecIdINDEO2          = AVCodecID(0x4b)
	AVCodecIdFRAPS           = AVCodecID(0x4c)
	AVCodecIdTRUEMOTION2     = AVCodecID(0x4d)
	AVCodecIdBMP             = AVCodecID(0x4e)
	AVCodecIdCSCD            = AVCodecID(0x4f)
	AVCodecIdMMVIDEO         = AVCodecID(0x50)
	AVCodecIdZMBV            = AVCodecID(0x51)
	AVCodecIdAVS             = AVCodecID(0x52)
	AVCodecIdSMACKVIDEO      = AVCodecID(0x53)
	AVCodecIdNUV             = AVCodecID(0x54)
	AVCodecIdKMVC            = AVCodecID(0x55)
	AVCodecIdFLASHSV         = AVCodecID(0x56)
	AVCodecIdCAVS            = AVCodecID(0x57)
	AVCodecIdJPEG2000        = AVCodecID(0x58)
	AVCodecIdVMNC            = AVCodecID(0x59)
	AVCodecIdVP5             = AVCodecID(0x5a)
	AVCodecIdVP6             = AVCodecID(0x5b)
	AVCodecIdVP6F            = AVCodecID(0x5c)
	AVCodecIdTARGA           = AVCodecID(0x5d)
	AVCodecIdDSICINVIDEO     = AVCodecID(0x5e)
	AVCodecIdTIERTEXSEQVIDEO = AVCodecID(0x5f)
	AVCodecIdTIFF            = AVCodecID(0x60)
	AVCodecIdGIF             = AVCodecID(0x61)
	AVCodecIdDXA             = AVCodecID(0x62)
	AVCodecIdDNXHD           = AVCodecID(0x63)
	AVCodecIdTHP             = AVCodecID(0x64)
	AVCodecIdSGI             = AVCodecID(0x65)
	AVCodecIdC93             = AVCodecID(0x66)
	AVCodecIdBETHSOFTVID     = AVCodecID(0x67)
	AVCodecIdPTX             = AVCodecID(0x68)
	AVCodecIdTXD             = AVCodecID(0x69)
	AVCodecIdVP6A            = AVCodecID(0x6a)
	AVCodecIdAMV             = AVCodecID(0x6b)
	AVCodecIdVB              = AVCodecID(0x6c)
	AVCodecIdPCX             = AVCodecID(0x6d)
	AVCodecIdSUNRAST         = AVCodecID(0x6e)
	AVCodecIdINDEO4          = AVCodecID(0x6f)
	AVCodecIdINDEO5          = AVCodecID(0x70)
	AVCodecIdMIMIC           = AVCodecID(0x71)
	AVCodecIdRL2             = AVCodecID(0x72)
	AVCodecIdESCAPE124       = AVCodecID(0x73)
	AVCodecIdDIRAC           = AVCodecID(0x74)
	AVCodecIdBFI             = AVCodecID(0x75)
	AVCodecIdCMV             = AVCodecID(0x76)
	AVCodecIdMOTIONPIXELS    = AVCodecID(0x77)
	AVCodecIdTGV             = AVCodecID(0x78)
	AVCodecIdTGQ             = AVCodecID(0x79)
	AVCodecIdTQI             = AVCodecID(0x7a)
	AVCodecIdAURA            = AVCodecID(0x7b)
	AVCodecIdAURA2           = AVCodecID(0x7c)
	AVCodecIdV210X           = AVCodecID(0x7d)
	AVCodecIdTMV             = AVCodecID(0x7e)
	AVCodecIdV210            = AVCodecID(0x7f)
	AVCodecIdDPX             = AVCodecID(0x80)
	AVCodecIdMAD             = AVCodecID(0x81)
	AVCodecIdFRWU            = AVCodecID(0x82)
	AVCodecIdFLASHSV2        = AVCodecID(0x83)
	AVCodecIdCDGRAPHICS      = AVCodecID(0x84)
	AVCodecIdR210            = AVCodecID(0x85)
	AVCodecIdANM             = AVCodecID(0x86)
	AVCodecIdBINKVIDEO       = AVCodecID(0x87)
	AVCodecIdIFFILBM         = AVCodecID(0x88)
	AVCodecIdIFFBYTERUN1     = AVCodecIdIFFILBM
	AVCodecIdKGV1            = AVCodecID(0x89)
	AVCodecIdYOP             = AVCodecID(0x8a)
	AVCodecIdVP8             = AVCodecID(0x8b)
	AVCodecIdPICTOR          = AVCodecID(0x8c)
	AVCodecIdANSI            = AVCodecID(0x8d)
	AVCodecIdA64MULTI        = AVCodecID(0x8e)
	AVCodecIdA64MULTI5       = AVCodecID(0x8f)
	AVCodecIdR10K            = AVCodecID(0x90)
	AVCodecIdMXPEG           = AVCodecID(0x91)
	AVCodecIdLAGARITH        = AVCodecID(0x92)
	AVCodecIdPRORES          = AVCodecID(0x93)
	AVCodecIdJV              = AVCodecID(0x94)
	AVCodecIdDFA             = AVCodecID(0x95)
	AVCodecIdWMV3IMAGE       = AVCodecID(0x96)
	AVCodecIdVC1IMAGE        = AVCodecID(0x97)
	AVCodecIdUTVIDEO         = AVCodecID(0x98)
	AVCodecIdBMVVIDEO        = AVCodecID(0x99)
	AVCodecIdVBLE            = AVCodecID(0x9a)
	AVCodecIdDXTORY          = AVCodecID(0x9b)
	AVCodecIdV410            = AVCodecID(0x9c)
	AVCodecIdXWD             = AVCodecID(0x9d)
	AVCodecIdCDXL            = AVCodecID(0x9e)
	AVCodecIdXBM             = AVCodecID(0x9f)
	AVCodecIdZEROCODEC       = AVCodecID(0xa0)
	AVCodecIdMSS1            = AVCodecID(0xa1)
	AVCodecIdMSA1            = AVCodecID(0xa2)
	AVCodecIdTSCC2           = AVCodecID(0xa3)
	AVCodecIdMTS2            = AVCodecID(0xa4)
	AVCodecIdCLLC            = AVCodecID(0xa5)
	AVCodecIdMSS2            = AVCodecID(0xa6)
	AVCodecIdVP9             = AVCodecID(0xa7)
	AVCodecIdAIC             = AVCodecID(0xa8)
	AVCodecIdESCAPE130       = AVCodecID(0xa9)
	AVCodecIdG2M             = AVCodecID(0xaa)
	AVCodecIdWEBP            = AVCodecID(0xab)
	AVCodecIdHNM4VIDEO       = AVCodecID(0xac)
	AVCodecIdHEVC            = AVCodecID(0xad)
	AVCodecIdH265            = AVCodecIdHEVC
	AVCodecIdFIC             = AVCodecID(0xae)
	AVCodecIdALIASPIX        = AVCodecID(0xaf)
	AVCodecIdBRENDERPIX      = AVCodecID(0xb0)
	AVCodecIdPAFVIDEO        = AVCodecID(0xb1)
	AVCodecIdEXR             = AVCodecID(0xb2)
	AVCodecIdVP7             = AVCodecID(0xb3)
	AVCodecIdSANM            = AVCodecID(0xb4)
	AVCodecIdSGIRLE          = AVCodecID(0xb5)
	AVCodecIdMVC1            = AVCodecID(0xb6)
	AVCodecIdMVC2            = AVCodecID(0xb7)
	AVCodecIdHQX             = AVCodecID(0xb8)
	AVCodecIdTDSC            = AVCodecID(0xb9)
	AVCodecIdHQHQA           = AVCodecID(0xba)
	AVCodecIdHAP             = AVCodecID(0xbb)
	AVCodecIdDDS             = AVCodecID(0xbc)
	AVCodecIdDXV             = AVCodecID(0xbd)
	AVCodecIdSCREENPRESSO    = AVCodecID(0xbe)
	AVCodecIdRSCC            = AVCodecID(0xbf)
	AVCodecIdAVS2            = AVCodecID(0xc0)
	AVCodecIdPGX             = AVCodecID(0xc1)
	AVCodecIdAVS3            = AVCodecID(0xc2)
	AVCodecIdMSP2            = AVCodecID(0xc3)
	AVCodecIdVVC             = AVCodecID(0xc4)
	AVCodecIdH266            = AVCodecIdVVC

	AVCodecIdY41P          = AVCodecID(0x8000)
	AVCodecIdAVRP          = AVCodecID(0x8001)
	AVCodecId012V          = AVCodecID(0x8002)
	AVCodecIdAVUI          = AVCodecID(0x8003)
	AVCodecIdAYUV          = AVCodecID(0x8004)
	AVCodecIdTARGAY216     = AVCodecID(0x8005)
	AVCodecIdV308          = AVCodecID(0x8006)
	AVCodecIdV408          = AVCodecID(0x8007)
	AVCodecIdYUV4          = AVCodecID(0x8008)
	AVCodecIdAVRN          = AVCodecID(0x8009)
	AVCodecIdCPIA          = AVCodecID(0x800a)
	AVCodecIdXFACE         = AVCodecID(0x800b)
	AVCodecIdSNOW          = AVCodecID(0x800c)
	AVCodecIdSMVJPEG       = AVCodecID(0x800d)
	AVCodecIdAPNG          = AVCodecID(0x800e)
	AVCodecIdDAALA         = AVCodecID(0x800f)
	AVCodecIdCFHD          = AVCodecID(0x8010)
	AVCodecIdTRUEMOTION2RT = AVCodecID(0x8011)
	AVCodecIdM101          = AVCodecID(0x8012)
	AVCodecIdMAGICYUV      = AVCodecID(0x8013)
	AVCodecIdSHEERVIDEO    = AVCodecID(0x8014)
	AVCodecIdYLC           = AVCodecID(0x8015)
	AVCodecIdPSD           = AVCodecID(0x8016)
	AVCodecIdPIXLET        = AVCodecID(0x8017)
	AVCodecIdSPEEDHQ       = AVCodecID(0x8018)
	AVCodecIdFMVC          = AVCodecID(0x8019)
	AVCodecIdSCPR          = AVCodecID(0x801a)
	AVCodecIdCLEARVIDEO    = AVCodecID(0x801b)
	AVCodecIdXPM           = AVCodecID(0x801c)
	AVCodecIdAV1           = AVCodecID(0x801d)
	AVCodecIdBITPACKED     = AVCodecID(0x801e)
	AVCodecIdMSCC          = AVCodecID(0x801f)
	AVCodecIdSRGC          = AVCodecID(0x8020)
	AVCodecIdSVG           = AVCodecID(0x8021)
	AVCodecIdGDV           = AVCodecID(0x8022)
	AVCodecIdFITS          = AVCodecID(0x8023)
	AVCodecIdIMM4          = AVCodecID(0x8024)
	AVCodecIdPROSUMER      = AVCodecID(0x8025)
	AVCodecIdMWSC          = AVCodecID(0x8026)
	AVCodecIdWCMV          = AVCodecID(0x8027)
	AVCodecIdRASC          = AVCodecID(0x8028)
	AVCodecIdHYMT          = AVCodecID(0x8029)
	AVCodecIdARBC          = AVCodecID(0x802a)
	AVCodecIdAGM           = AVCodecID(0x802b)
	AVCodecIdLSCR          = AVCodecID(0x802c)
	AVCodecIdVP4           = AVCodecID(0x802d)
	AVCodecIdIMM5          = AVCodecID(0x802e)
	AVCodecIdMVDV          = AVCodecID(0x802f)
	AVCodecIdMVHA          = AVCodecID(0x8030)
	AVCodecIdCDTOONS       = AVCodecID(0x8031)
	AVCodecIdMV30          = AVCodecID(0x8032)
	AVCodecIdNOTCHLC       = AVCodecID(0x8033)
	AVCodecIdPFM           = AVCodecID(0x8034)
	AVCodecIdMOBICLIP      = AVCodecID(0x8035)
	AVCodecIdPHOTOCD       = AVCodecID(0x8036)
	AVCodecIdIPU           = AVCodecID(0x8037)
	AVCodecIdARGO          = AVCodecID(0x8038)
	AVCodecIdCRI           = AVCodecID(0x8039)
	AVCodecIdSIMBIOSISIMX  = AVCodecID(0x803a)
	AVCodecIdSGAVIDEO      = AVCodecID(0x803b)

	/* various PCM "codecs" */
	AVCodecIdFIRSTAUDIO     = AVCodecID(0x10000) ///< A dummy id pointing at the start of audio codecs
	AVCodecIdPCMS16LE       = AVCodecID(0x10000)
	AVCodecIdPCMS16BE       = AVCodecID(0x10001)
	AVCodecIdPCMU16LE       = AVCodecID(0x10002)
	AVCodecIdPCMU16BE       = AVCodecID(0x10003)
	AVCodecIdPCMS8          = AVCodecID(0x10004)
	AVCodecIdPCMU8          = AVCodecID(0x10005)
	AVCodecIdPCMMULAW       = AVCodecID(0x10006)
	AVCodecIdPCMALAW        = AVCodecID(0x10007)
	AVCodecIdPCMS32LE       = AVCodecID(0x10008)
	AVCodecIdPCMS32BE       = AVCodecID(0x10009)
	AVCodecIdPCMU32LE       = AVCodecID(0x1000a)
	AVCodecIdPCMU32BE       = AVCodecID(0x1000b)
	AVCodecIdPCMS24LE       = AVCodecID(0x1000c)
	AVCodecIdPCMS24BE       = AVCodecID(0x1000d)
	AVCodecIdPCMU24LE       = AVCodecID(0x1000e)
	AVCodecIdPCMU24BE       = AVCodecID(0x1000f)
	AVCodecIdPCMS24DAUD     = AVCodecID(0x10010)
	AVCodecIdPCMZORK        = AVCodecID(0x10011)
	AVCodecIdPCMS16LEPLANAR = AVCodecID(0x10012)
	AVCodecIdPCMDVD         = AVCodecID(0x10013)
	AVCodecIdPCMF32BE       = AVCodecID(0x10014)
	AVCodecIdPCMF32LE       = AVCodecID(0x10015)
	AVCodecIdPCMF64BE       = AVCodecID(0x10016)
	AVCodecIdPCMF64LE       = AVCodecID(0x10017)
	AVCodecIdPCMBLURAY      = AVCodecID(0x10018)
	AVCodecIdPCMLXF         = AVCodecID(0x10019)
	AVCodecIdS302M          = AVCodecID(0x1001a)
	AVCodecIdPCMS8PLANAR    = AVCodecID(0x1001b)
	AVCodecIdPCMS24LEPLANAR = AVCodecID(0x1001c)
	AVCodecIdPCMS32LEPLANAR = AVCodecID(0x1001d)
	AVCodecIdPCMS16BEPLANAR = AVCodecID(0x1001e)

	AVCodecIdPCMS64LE = AVCodecID(0x10800)
	AVCodecIdPCMS64BE = AVCodecID(0x10801)
	AVCodecIdPCMF16LE = AVCodecID(0x10802)
	AVCodecIdPCMF24LE = AVCodecID(0x10803)
	AVCodecIdPCMVIDC  = AVCodecID(0x10804)
	AVCodecIdPCMSGA   = AVCodecID(0x10805)

	/* various ADPCM codecs */
	AVCodecIdADPCMIMAQT     = AVCodecID(0x11000)
	AVCodecIdADPCMIMAWAV    = AVCodecID(0x11001)
	AVCodecIdADPCMIMADK3    = AVCodecID(0x11002)
	AVCodecIdADPCMIMADK4    = AVCodecID(0x11003)
	AVCodecIdADPCMIMAWS     = AVCodecID(0x11004)
	AVCodecIdADPCMIMASMJPEG = AVCodecID(0x11005)
	AVCodecIdADPCMMS        = AVCodecID(0x11006)
	AVCodecIdADPCM4XM       = AVCodecID(0x11007)
	AVCodecIdADPCMXA        = AVCodecID(0x11008)
	AVCodecIdADPCMADX       = AVCodecID(0x11009)
	AVCodecIdADPCMEA        = AVCodecID(0x1100a)
	AVCodecIdADPCMG726      = AVCodecID(0x1100b)
	AVCodecIdADPCMCT        = AVCodecID(0x1100c)
	AVCodecIdADPCMSWF       = AVCodecID(0x1100d)
	AVCodecIdADPCMYAMAHA    = AVCodecID(0x1100e)
	AVCodecIdADPCMSBPRO4    = AVCodecID(0x1100f)
	AVCodecIdADPCMSBPRO3    = AVCodecID(0x11010)
	AVCodecIdADPCMSBPRO2    = AVCodecID(0x11011)
	AVCodecIdADPCMTHP       = AVCodecID(0x11012)
	AVCodecIdADPCMIMAAMV    = AVCodecID(0x11013)
	AVCodecIdADPCMEAR1      = AVCodecID(0x11014)
	AVCodecIdADPCMEAR3      = AVCodecID(0x11015)
	AVCodecIdADPCMEAR2      = AVCodecID(0x11016)
	AVCodecIdADPCMIMAEASEAD = AVCodecID(0x11017)
	AVCodecIdADPCMIMAEAEACS = AVCodecID(0x11018)
	AVCodecIdADPCMEAXAS     = AVCodecID(0x11019)
	AVCodecIdADPCMEAMAXISXA = AVCodecID(0x1101a)
	AVCodecIdADPCMIMAISS    = AVCodecID(0x1101b)
	AVCodecIdADPCMG722      = AVCodecID(0x1101c)
	AVCodecIdADPCMIMAAPC    = AVCodecID(0x1101d)
	AVCodecIdADPCMVIMA      = AVCodecID(0x1101e)

	AVCodecIdADPCMAFC        = AVCodecID(0x11800)
	AVCodecIdADPCMIMAOKI     = AVCodecID(0x11801)
	AVCodecIdADPCMDTK        = AVCodecID(0x11802)
	AVCodecIdADPCMIMARAD     = AVCodecID(0x11803)
	AVCodecIdADPCMG726LE     = AVCodecID(0x11804)
	AVCodecIdADPCMTHPLE      = AVCodecID(0x11805)
	AVCodecIdADPCMPSX        = AVCodecID(0x11806)
	AVCodecIdADPCMAICA       = AVCodecID(0x11807)
	AVCodecIdADPCMIMADAT4    = AVCodecID(0x11808)
	AVCodecIdADPCMMTAF       = AVCodecID(0x11809)
	AVCodecIdADPCMAGM        = AVCodecID(0x1180a)
	AVCodecIdADPCMARGO       = AVCodecID(0x1180b)
	AVCodecIdADPCMIMASSI     = AVCodecID(0x1180c)
	AVCodecIdADPCMZORK       = AVCodecID(0x1180d)
	AVCodecIdADPCMIMAAPM     = AVCodecID(0x1180e)
	AVCodecIdADPCMIMAALP     = AVCodecID(0x1180f)
	AVCodecIdADPCMIMAMTF     = AVCodecID(0x11810)
	AVCodecIdADPCMIMACUNNING = AVCodecID(0x11811)
	AVCodecIdADPCMIMAMOFLEX  = AVCodecID(0x11812)

	/* AMR */
	AVCodecIdAMRNB = AVCodecID(0x12000)
	AVCodecIdAMRWB = AVCodecID(0x12001)

	/* RealAudio codecs*/
	AVCodecIdRA144 = AVCodecID(0x13000)
	AVCodecIdRA288 = AVCodecID(0x13001)

	/* various DPCM codecs */
	AVCodecIdROQDPCM       = AVCodecID(0x14000)
	AVCodecIdINTERPLAYDPCM = AVCodecID(0x14001)
	AVCodecIdXANDPCM       = AVCodecID(0x14002)
	AVCodecIdSOLDPCM       = AVCodecID(0x14003)

	AVCodecIdSDX2DPCM    = AVCodecID(0x14800)
	AVCodecIdGREMLINDPCM = AVCodecID(0x14801)
	AVCodecIdDERFDPCM    = AVCodecID(0x14802)

	/* audio codecs */
	AVCodecIdMP2           = AVCodecID(0x15000)
	AVCodecIdMP3           = AVCodecID(0x15001) ///< preferred ID for decoding MPEG audio layer 1, 2 or 3
	AVCodecIdAAC           = AVCodecID(0x15002)
	AVCodecIdAC3           = AVCodecID(0x15003)
	AVCodecIdDTS           = AVCodecID(0x15004)
	AVCodecIdVORBIS        = AVCodecID(0x15005)
	AVCodecIdDVAUDIO       = AVCodecID(0x15006)
	AVCodecIdWMAV1         = AVCodecID(0x15007)
	AVCodecIdWMAV2         = AVCodecID(0x15008)
	AVCodecIdMACE3         = AVCodecID(0x15009)
	AVCodecIdMACE6         = AVCodecID(0x1500a)
	AVCodecIdVMDAUDIO      = AVCodecID(0x1500b)
	AVCodecIdFLAC          = AVCodecID(0x1500c)
	AVCodecIdMP3ADU        = AVCodecID(0x1500d)
	AVCodecIdMP3ON4        = AVCodecID(0x1500e)
	AVCodecIdSHORTEN       = AVCodecID(0x1500f)
	AVCodecIdALAC          = AVCodecID(0x15010)
	AVCodecIdWESTWOODSND1  = AVCodecID(0x15011)
	AVCodecIdGSM           = AVCodecID(0x15012) ///< as in Berlin toast format
	AVCodecIdQDM2          = AVCodecID(0x15013)
	AVCodecIdCOOK          = AVCodecID(0x15014)
	AVCodecIdTRUESPEECH    = AVCodecID(0x15015)
	AVCodecIdTTA           = AVCodecID(0x15016)
	AVCodecIdSMACKAUDIO    = AVCodecID(0x15017)
	AVCodecIdQCELP         = AVCodecID(0x15018)
	AVCodecIdWAVPACK       = AVCodecID(0x15019)
	AVCodecIdDSICINAUDIO   = AVCodecID(0x1501a)
	AVCodecIdIMC           = AVCodecID(0x1501b)
	AVCodecIdMUSEPACK7     = AVCodecID(0x1501c)
	AVCodecIdMLP           = AVCodecID(0x1501d)
	AVCodecIdGSMMS         = AVCodecID(0x1501e) /* as found in WAV */
	AVCodecIdATRAC3        = AVCodecID(0x1501f)
	AVCodecIdAPE           = AVCodecID(0x15020)
	AVCodecIdNELLYMOSER    = AVCodecID(0x15021)
	AVCodecIdMUSEPACK8     = AVCodecID(0x15022)
	AVCodecIdSPEEX         = AVCodecID(0x15023)
	AVCodecIdWMAVOICE      = AVCodecID(0x15024)
	AVCodecIdWMAPRO        = AVCodecID(0x15025)
	AVCodecIdWMALOSSLESS   = AVCodecID(0x15026)
	AVCodecIdATRAC3P       = AVCodecID(0x15027)
	AVCodecIdEAC3          = AVCodecID(0x15028)
	AVCodecIdSIPR          = AVCodecID(0x15029)
	AVCodecIdMP1           = AVCodecID(0x1502a)
	AVCodecIdTWINVQ        = AVCodecID(0x1502b)
	AVCodecIdTRUEHD        = AVCodecID(0x1502c)
	AVCodecIdMP4ALS        = AVCodecID(0x1502d)
	AVCodecIdATRAC1        = AVCodecID(0x1502e)
	AVCodecIdBINKAUDIORDFT = AVCodecID(0x1502f)
	AVCodecIdBINKAUDIODCT  = AVCodecID(0x15030)
	AVCodecIdAACLATM       = AVCodecID(0x15031)
	AVCodecIdQDMC          = AVCodecID(0x15032)
	AVCodecIdCELT          = AVCodecID(0x15033)
	AVCodecIdG7231         = AVCodecID(0x15034)
	AVCodecIdG729          = AVCodecID(0x15035)
	AVCodecId8SVXEXP       = AVCodecID(0x15036)
	AVCodecId8SVXFIB       = AVCodecID(0x15037)
	AVCodecIdBMVAUDIO      = AVCodecID(0x15038)
	AVCodecIdRALF          = AVCodecID(0x15039)
	AVCodecIdIAC           = AVCodecID(0x1503a)
	AVCodecIdILBC          = AVCodecID(0x1503b)
	AVCodecIdOPUS          = AVCodecID(0x1503c)
	AVCodecIdCOMFORTNOISE  = AVCodecID(0x1503d)
	AVCodecIdTAK           = AVCodecID(0x1503e)
	AVCodecIdMETASOUND     = AVCodecID(0x1503f)
	AVCodecIdPAFAUDIO      = AVCodecID(0x15040)
	AVCodecIdON2AVC        = AVCodecID(0x15041)
	AVCodecIdDSSSP         = AVCodecID(0x15042)
	AVCodecIdCODEC2        = AVCodecID(0x15043)

	AVCodecIdFFWAVESYNTH   = 0x15800
	AVCodecIdSONIC         = AVCodecID(0x15801)
	AVCodecIdSONICLS       = AVCodecID(0x15802)
	AVCodecIdEVRC          = AVCodecID(0x15803)
	AVCodecIdSMV           = AVCodecID(0x15804)
	AVCodecIdDSDLSBF       = AVCodecID(0x15805)
	AVCodecIdDSDMSBF       = AVCodecID(0x15806)
	AVCodecIdDSDLSBFPLANAR = AVCodecID(0x15807)
	AVCodecIdDSDMSBFPLANAR = AVCodecID(0x15808)
	AVCodecId4GV           = AVCodecID(0x15809)
	AVCodecIdINTERPLAYACM  = AVCodecID(0x1580a)
	AVCodecIdXMA1          = AVCodecID(0x1580b)
	AVCodecIdXMA2          = AVCodecID(0x1580c)
	AVCodecIdDST           = AVCodecID(0x1580d)
	AVCodecIdATRAC3AL      = AVCodecID(0x1580e)
	AVCodecIdATRAC3PAL     = AVCodecID(0x1580f)
	AVCodecIdDOLBYE        = AVCodecID(0x15810)
	AVCodecIdAPTX          = AVCodecID(0x15811)
	AVCodecIdAPTXHD        = AVCodecID(0x15812)
	AVCodecIdSBC           = AVCodecID(0x15813)
	AVCodecIdATRAC9        = AVCodecID(0x15814)
	AVCodecIdHCOM          = AVCodecID(0x15815)
	AVCodecIdACELPKELVIN   = AVCodecID(0x15816)
	AVCodecIdMPEGH3DAUDIO  = AVCodecID(0x15817)
	AVCodecIdSIREN         = AVCodecID(0x15818)
	AVCodecIdHCA           = AVCodecID(0x15819)
	AVCodecIdFASTAUDIO     = AVCodecID(0x1581a)

	/* subtitle codecs */
	AVCodecIdFIRSTSUBTITLE   = AVCodecID(0x17000) ///< A dummy ID pointing at the start of subtitle codecs.
	AVCodecIdDVDSUBTITLE     = AVCodecID(0x17000)
	AVCodecIdDVBSUBTITLE     = AVCodecID(0x17001)
	AVCodecIdTEXT            = AVCodecID(0x17002) ///< raw UTF-8 text
	AVCodecIdXSUB            = AVCodecID(0x17003)
	AVCodecIdSSA             = AVCodecID(0x17004)
	AVCodecIdMOVTEXT         = AVCodecID(0x17005)
	AVCodecIdHDMVPGSSUBTITLE = AVCodecID(0x17006)
	AVCodecIdDVBTELETEXT     = AVCodecID(0x17007)
	AVCodecIdSRT             = AVCodecID(0x17008)

	AVCodecIdMICRODVD         = AVCodecID(0x17800)
	AVCodecIdEIA608           = AVCodecID(0x17801)
	AVCodecIdJACOSUB          = AVCodecID(0x17802)
	AVCodecIdSAMI             = AVCodecID(0x17803)
	AVCodecIdREALTEXT         = AVCodecID(0x17804)
	AVCodecIdSTL              = AVCodecID(0x17805)
	AVCodecIdSUBVIEWER1       = AVCodecID(0x17806)
	AVCodecIdSUBVIEWER        = AVCodecID(0x17807)
	AVCodecIdSUBRIP           = AVCodecID(0x17808)
	AVCodecIdWEBVTT           = AVCodecID(0x17809)
	AVCodecIdMPL2             = AVCodecID(0x1780A)
	AVCodecIdVPLAYER          = AVCodecID(0x1780B)
	AVCodecIdPJS              = AVCodecID(0x1780C)
	AVCodecIdASS              = AVCodecID(0x1780D)
	AVCodecIdHDMVTEXTSUBTITLE = AVCodecID(0x1780E)
	AVCodecIdTTML             = AVCodecID(0x178F0)
	AVCodecIdARIBCAPTION      = AVCodecID(0x178F1)

	/* other specific kind of codecs (generally used for attachments) */
	AVCodecIdFIRSTUNKNOWN = AVCodecID(0x18000) ///< A dummy ID pointing at the start of various fake codecs.
	AVCodecIdTTF          = AVCodecID(0x18000)

	AVCodecIdSCTE35   = AVCodecID(0x18001) ///< Contain timestamp estimated through PCR of program stream.
	AVCodecIdEPG      = AVCodecID(0x18002)
	AVCodecIdBINTEXT  = AVCodecID(0x18800)
	AVCodecIdXBIN     = AVCodecID(0x18801)
	AVCodecIdIDF      = AVCodecID(0x18802)
	AVCodecIdOTF      = AVCodecID(0x18803)
	AVCodecIdSMPTEKLV = AVCodecID(0x18804)
	AVCodecIdDVDNAV   = AVCodecID(0x18805)
	AVCodecIdTIMEDID3 = AVCodecID(0x18806)
	AVCodecIdBINDATA  = AVCodecID(0x18807)

	AVCodecIdPROBE = AVCodecID(0x19000) ///< codecId is not known (like AVCodecIdNONE) but lavf should attempt to identify it

	AVCodecIdMPEG2TS = AVCodecID(0x20000) /**< FAKE codec to indicate a raw MPEG-2 TS
	 * stream (only used by libavformat) */
	AVCodecIdMPEG4SYSTEMS = AVCodecID(0x20001) /**< FAKE codec to indicate a MPEG-4 Systems
	 * stream (only used by libavformat) */
	AVCodecIdFFMETADATA     = AVCodecID(0x21000) ///< Dummy codec for streams containing only metadata information.
	AVCodecIdWRAPPEDAVFRAME = AVCodecID(0x21001) ///< Passthrough codec, AVFrames wrapped in AVPacket
)
