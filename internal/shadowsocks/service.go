package shadowsocks

import (
	"encoding/base64"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/quarkcms/quark-go/internal/models"
	"github.com/shadowsocks/go-shadowsocks2/core"
)

var verbose = true
var TCPCork = true
var UDPTimeout = 5 * time.Minute

func Start() {

	serverInfo := (&models.Server{}).Info(1)

	addr := "ss://" + serverInfo.EncryptType + ":" + serverInfo.Password + "@:" + serverInfo.Port
	cipher := ""
	password := ""
	plugin := serverInfo.Plugin
	pluginOpts := serverInfo.PluginOpts
	flagsKey := serverInfo.Key
	flagsUDP := false
	flagsTcp := true
	var err error

	var key []byte
	if flagsKey != "" {
		k, err := base64.URLEncoding.DecodeString(flagsKey)
		if err != nil {
			log.Fatal(err)
		}
		key = k
	}

	if strings.HasPrefix(addr, "ss://") {
		addr, cipher, password, err = parseURL(addr)
		if err != nil {
			log.Fatal(err)
		}
	}

	udpAddr := addr

	if plugin != "" {
		addr, err = startPlugin(plugin, pluginOpts, addr, true)
		if err != nil {
			log.Fatal(err)
		}
	}

	ciph, err := core.PickCipher(cipher, key, password)
	if err != nil {
		log.Fatal(err)
	}

	if flagsUDP {
		go udpRemote(udpAddr, ciph.PacketConn)
	}
	if flagsTcp {
		go tcpRemote(addr, ciph.StreamConn)
	}

	killPlugin()
}

func parseURL(s string) (addr, cipher, password string, err error) {
	u, err := url.Parse(s)
	if err != nil {
		return
	}

	addr = u.Host
	if u.User != nil {
		cipher = u.User.Username()
		password, _ = u.User.Password()
	}
	return
}
