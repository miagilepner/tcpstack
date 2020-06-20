package main

import (
	"flag"
	"os/user"

	"github.com/miagilepner/tcpstack/socket"
	"github.com/miagilepner/tcpstack/tcp"
	"github.com/rs/zerolog/log"
)

func main() {
	ip := flag.String("dest_ip", "1.2.3.4", "destination IP")
	port := flag.Int("dest_port", 0, "destination port")
	iface := flag.String("iface", "en0", "interface name")
	flag.Parse()
	us, err := user.Current()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get user")
	}
	log.Info().Msgf("running as %s", us.Username)
	sock := socket.NewSocket()
	attachErr := sock.Attach(*iface)
	if attachErr != nil {
		log.Fatal().Err(attachErr).Msg("failed to attach")
	}
	if err := sock.Close(); err != nil {
		log.Fatal().Err(err).Msg("closing")
	}
	tcp.Handshake(sock, *ip, *port)
}
