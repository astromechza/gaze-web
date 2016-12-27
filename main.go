package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"path/filepath"

	"github.com/kataras/go-template/html"
	"github.com/kataras/iris"
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

// VersionString is the version string inserted by whatever build script
// format should be 'X.YZ'
// Set this at build time using the -ldflags="-X main.VersionString=X.YZ"
var VersionString = "<unofficial build>"

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
		fmt.Println("Version: " + VersionString)
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

	iris.Static("/static", "./static", 1)
	iris.UseTemplate(html.New(html.Config{
		Layout: "root/layout.html",
	})).Directory(filepath.Join(srvDir, "templates"), ".html")

	iris.Get("/", indexHandler)
	iris.Post("/report", newReportHandler)
	iris.Put("/report", newReportHandler)
	iris.Get("/reports", listReportsHandler)
	iris.Get("/reports/:ulid", getReportHandler)

	iris.OnError(iris.StatusInternalServerError, error500Handler)

	iris.Listen(fmt.Sprintf(":%v", *portFlag))
	return nil
}

func main() {
	if err := mainInner(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
