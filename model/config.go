package model


type Config struct {
	BeforeCommand []string
	AfterCommand []string
	IpAddress Ips
	NginxConfigPath string
	NginxUpdateGap int
	NginxDomain string//
	NginxRestartCommand string
	DataSavePath string
	OpenScanHost bool
	RetryTimesWhenFail int // 1 < times < 10
	FailAction string //goon rollback stop
	HookMail []string // send email when fail over maxtime

}

