// Copyright (c) 2013-2014, Jeffrey Wilcke. All rights reserved.
//
// This library is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public
// License as published by the Free Software Foundation; either
// version 2.1 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this library; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston,
// MA 02110-1301  USA

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"

	"bitbucket.org/kardianos/osext"
	"github.com/georzaza/go-ethereum-v0.7.10_official/logger"
	"github.com/georzaza/go-ethereum-v0.7.10_official/vm"
)

var (
	Identifier      string
	KeyRing         string
	KeyStore        string
	StartRpc        bool
	StartWebSockets bool
	RpcPort         int
	UseUPnP         bool
	OutboundPort    string
	ShowGenesis     bool
	AddPeer         string
	MaxPeer         int
	GenAddr         bool
	UseSeed         bool
	SecretFile      string
	ExportDir       string
	NonInteractive  bool
	Datadir         string
	LogFile         string
	ConfigFile      string
	DebugFile       string
	LogLevel        int
	VmType          int
)

// flags specific to gui client
var AssetPath string

//TODO: If we re-use the one defined in cmd.go the binary osx image crashes. If somebody finds out why we can dry this up.
func defaultAssetPath() string {
	var assetPath string
	// If the current working directory is the go-ethereum dir
	// assume a debug build and use the source directory as
	// asset directory.
	pwd, _ := os.Getwd()
	if pwd == path.Join(os.Getenv("GOPATH"), "src", "github.com", "ethereum", "go-ethereum", "cmd", "mist") {
		assetPath = path.Join(pwd, "assets")
	} else {
		switch runtime.GOOS {
		case "darwin":
			// Get Binary Directory
			exedir, _ := osext.ExecutableFolder()
			assetPath = filepath.Join(exedir, "../Resources")
		case "linux":
			assetPath = "/usr/share/mist"
		case "windows":
			assetPath = "./assets"
		default:
			assetPath = "."
		}
	}
	return assetPath
}
func defaultDataDir() string {
	usr, _ := user.Current()
	return path.Join(usr.HomeDir, ".mist")
}

var defaultConfigFile = path.Join(defaultDataDir(), "conf.ini")

func Init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s [options] [filename]:\noptions precedence: default < config file < environment variables < command line\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.IntVar(&VmType, "vm", 0, "Virtual Machine type: 0-1: standard, debug")
	flag.StringVar(&Identifier, "id", "", "Custom client identifier")
	flag.StringVar(&KeyRing, "keyring", "", "identifier for keyring to use")
	flag.StringVar(&KeyStore, "keystore", "db", "system to store keyrings: db|file (db)")
	flag.StringVar(&OutboundPort, "port", "30303", "listening port")
	flag.BoolVar(&UseUPnP, "upnp", true, "enable UPnP support")
	flag.IntVar(&MaxPeer, "maxpeer", 30, "maximum desired peers")
	flag.IntVar(&RpcPort, "rpcport", 8080, "port to start json-rpc server on")
	flag.BoolVar(&StartRpc, "rpc", false, "start rpc server")
	flag.BoolVar(&StartWebSockets, "ws", false, "start websocket server")
	flag.BoolVar(&NonInteractive, "y", false, "non-interactive mode (say yes to confirmations)")
	flag.BoolVar(&UseSeed, "seed", true, "seed peers")
	flag.BoolVar(&GenAddr, "genaddr", false, "create a new priv/pub key")
	flag.StringVar(&SecretFile, "import", "", "imports the file given (hex or mnemonic formats)")
	flag.StringVar(&ExportDir, "export", "", "exports the session keyring to files in the directory given")
	flag.StringVar(&LogFile, "logfile", "", "log file (defaults to standard output)")
	flag.StringVar(&Datadir, "datadir", defaultDataDir(), "specifies the datadir to use")
	flag.StringVar(&ConfigFile, "conf", defaultConfigFile, "config file")
	flag.StringVar(&DebugFile, "debug", "", "debug file (no debugging if not set)")
	flag.IntVar(&LogLevel, "loglevel", int(logger.InfoLevel), "loglevel: 0-5: silent,error,warn,info,debug,debug detail)")

	flag.StringVar(&AssetPath, "asset_path", defaultAssetPath(), "absolute path to GUI assets directory")

	flag.Parse()

	if VmType >= int(vm.MaxVmTy) {
		log.Fatal("Invalid VM type ", VmType)
	}
}
