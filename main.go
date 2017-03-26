package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/AstromechZA/gaze-web/storage"
	boltstorage "github.com/AstromechZA/gaze-web/storage/bolt"
	"github.com/AstromechZA/gaze-web/utils"
	"github.com/AstromechZA/gaze-web/webserver"
)

const usageString = `
Running gaze-web will launch a small webserver on this host with a sqlite database
in the current directory. The sqlite database is used to store records received
over its http API.
`

const logoImage = `
 _______  _______  _______  _______               _______  ______
(  ____ \(  ___  )/ ___   )(  ____ \    |\     /|(  ____ \(  ___ \
| (    \/| (   ) |\/   )  || (    \/    | )   ( || (    \/| (   ) )
| |      | (___) |    /   )| (__  _____ | | _ | || (__    | (__/ /
| | ____ |  ___  |   /   / |  __)(_____)| |( )| ||  __)   |  __ (
| | \_  )| (   ) |  /   /  | (          | || || || (      | (  \ \
| (___) || )   ( | /   (_/\| (____/\    | () () || (____/\| )___) )
(_______)|/     \|(_______/(_______/    (_______)(_______/|/ \___/
`

// These variables are filled by the `govvv` tool at compile time.
// There are a few more granular variables available if necessary.
var Version = "<unofficial build>"
var GitSummary = "<changes unknown>"
var BuildDate = "<no date>"

func mainInner() error {

	// first set up config flag options
	versionFlag := flag.Bool("version", false, "Print the version string")
	portFlag := flag.Int("port", 8080, "the port to listen on")

	// set a more verbose usage message.
	flag.Usage = func() {
		os.Stderr.WriteString(strings.TrimSpace(usageString) + "\n\n")
		flag.PrintDefaults()
	}
	// parse them
	flag.Parse()

	// do arg checking
	if *versionFlag {
		fmt.Printf("Version: %s (%s) on %s [%s]\n", Version, GitSummary, BuildDate, runtime.Version())
		fmt.Println(logoImage)
		fmt.Println("Project: github.com/AstromechZA/gaze-web")
		return nil
	}
	if *portFlag <= 0 {
		return fmt.Errorf("Port must be > 0")
	}

	utils.EmbeddedVersionString = fmt.Sprintf("Version: %s (%s) | %s | %s", Version, GitSummary, BuildDate, runtime.Version())

	// construct server directory
	srvDir, _ := filepath.Abs(os.Args[0])
	srvDir = filepath.Dir(srvDir)

	// set up database and models
	store, err := boltstorage.SetupBoltDBReportStore("gaze-web.db")
	if err != nil {
		return err
	}
	defer store.Close()
	storage.ActiveStore = store

	app := webserver.Setup(srvDir)
	app.Listen(fmt.Sprintf(":%v", *portFlag))
	return nil
}

func main() {
	if err := mainInner(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
