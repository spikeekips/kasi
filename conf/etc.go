package kasi_conf

import "github.com/spikeekips/kasi/util"

// SSLSetting
type SSLSetting struct {
	Cert string
	Key  string
	Pem  string
}

func (setting *SSLSetting) String() string {
	return kasi_util.ToJson(setting)
}
