package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"path/filepath"

	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/view"
)

const usageString = `
TODO
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

var RandomSource = rand.New(rand.NewSource(time.Now().UnixNano()))

var Database DatabaseRef

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
		fmt.Printf("Version: %s (%s) on %s \n", Version, GitSummary, BuildDate)
		fmt.Println(logoImage)
		fmt.Println("Project: <project url here>")
		return nil
	}
	if *portFlag <= 0 {
		return fmt.Errorf("Port must be > 0")
	}

	// construct server directory
	srvDir, _ := filepath.Abs(os.Args[0])
	srvDir = filepath.Dir(srvDir)

	// set up database and models
	db, err := initDatabase("test.db")
	if err != nil {
		return err
	}
	defer db.Close()
	Database.Active = db

	app := iris.New()

	app.Adapt(iris.DevLogger())
	app.Adapt(httprouter.New())

	engine := view.HTML(filepath.Join(srvDir, "templates"), ".html")
	engine.Layout("root/layout.html")
	engine.Funcs(buildTemplateFuncsMap())
	app.Adapt(engine)

	app.Use(loggerMiddleware{})

	app.StaticWeb("/static", filepath.Join(srvDir, "static"))

	app.Get("/", indexHandler)
	app.Post("/report", newReportHandler)
	app.Put("/report", newReportHandler)
	app.Get("/reports", listReportsHandler)
	app.Get("/reports/:ulid", getReportHandler)

	app.OnError(iris.StatusInternalServerError, error500Handler)

	app.Listen(fmt.Sprintf(":%v", *portFlag))
	return nil
}

func main() {
	if err := mainInner(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
