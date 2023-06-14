package elk

import "testing"

func TestXxx(t *testing.T) {
	InitEs("192.168.4.41", "9200")
	InitLog("192.168.4.41", "index")

	LogrusObj.Debugln("11")
	LogrusObj.Infoln("11")
	LogrusObj.Warningln("11")
	LogrusObj.Errorln("11")
}
