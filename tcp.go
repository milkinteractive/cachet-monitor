package cachet

import (
	"fmt"
	"net"
	"time"

	"github.com/Sirupsen/logrus"
)

// Investigating template
var defaultTCPInvestigatingTpl = MessageTemplate{
	Subject: `{{ .Monitor.Name }} - {{ .SystemName }}`,
	Message: `{{ .Monitor.Name }} check **failed** (server time: {{ .now }})

{{ .FailReason }}`,
}

// Fixed template
var defaultTCPFixedTpl = MessageTemplate{
	Subject: `{{ .Monitor.Name }} - {{ .SystemName }}`,
	Message: `**Resolved** - {{ .now }}

- - -

{{ .incident.Message }}`,
}

// TCPMonitor struct
type TCPMonitor struct {
	AbstractMonitor `mapstructure:",squash"`
	Port            string
}

// CheckTCPPortAlive func
func CheckTCPPortAlive(ip, port string, timeout int64) (bool, error) {

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), time.Duration(timeout)*time.Second)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		return false, err
	} else {
		return true, nil
	}

}

// test if it available
func (m *TCPMonitor) test() bool {
	if alive, e := CheckTCPPortAlive(m.Target, m.Port, int64(m.Timeout)); alive {
		return true
	} else {
		msg := fmt.Sprintf("TCP check failed: %v", e)
		logrus.Error(msg)
		m.lastFailReason = msg
		return false
	}
}

// Validate configuration
func (m *TCPMonitor) Validate() []string {

	// set incident temp
	m.Template.Investigating.SetDefault(defaultTCPInvestigatingTpl)
	m.Template.Fixed.SetDefault(defaultTCPFixedTpl)

	// super.Validate()
	errs := m.AbstractMonitor.Validate()

	if m.Target == "" {
		errs = append(errs, "Target is required")
	}

	if m.Port == "" {
		errs = append(errs, "Port is required")
	}

	return errs
}