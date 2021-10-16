package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"os"
	"os/user"
)

type Archive struct {
	Body       string
	Categories string
	Date       string
	Filename   string
	Title      string
	Updated    string
}

type Post struct {
	ID        int    `json:"id"`
	Author    string `json:"author"`
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	Markdown  string `json:"markdown"`
	HTML      string `json:"html"`
	Published string `json:"published"`
	Updated   string `json:"updated"`
    Highlight *int    `json:"highlight"`
}

type Database struct {
	Type             string
	ConnectionString string
}

func (d *Database) getPosts() ([]Post, error) {
	db, err := sql.Open(d.Type, d.ConnectionString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	results, err := db.Query("SELECT * FROM entries")
	if err != nil {
		return nil, err
	}
	var posts []Post
	for results.Next() {
		var post Post
		err = results.Scan(
			&post.ID,
			&post.Author,
			&post.Slug,
			&post.Title,
			&post.Markdown,
			&post.HTML,
			&post.Published,
			&post.Updated,
            &post.Highlight,
		)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (a *Archive) write() error {
	data := []byte(
		"---\n" +
			a.Title + "\n" +
			a.Date + "\n" +
			a.Updated + "\n" +
			"categories: [\"Archive\"]\n" +
			"---\n" +
			"\n" +
			a.Body)
	fmt.Printf("Writing %v\n", a.Filename)
	err := ioutil.WriteFile(a.Filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	fmt.Println("blogshovel: mysql -> archive(markdown)\n")

	// Setup default output path
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	outPath := usr.HomeDir + "/archive/"

	dryrun := flag.Bool("dryrun", false, "skip writing to file")
	outdir := flag.String("outdir", outPath, "output files to this directory")
	// Set a database connection string, this needs improvement.
	defaultDbString := "user:password@tcp(host.domain:3306)/database"
	dbstring := flag.String("dbconnstring", defaultDbString, "database connection string")
	flag.Parse()
	if *dbstring == defaultDbString {
		fmt.Printf("You must provide a -dbconstring\n")
		flag.Usage()
		os.Exit(1)
	}

	var db Database
	db.Type = "mysql"
	db.ConnectionString = *dbstring
	posts, err := db.getPosts()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Found %d posts to export\n", len(posts))
	for _, p := range posts {
		fmt.Printf("Archiving: %v\n", p.Slug)
		out := Archive{
			Filename:   *outdir + p.Slug + ".md",
			Title:      "title: " + p.Title,
			Date:       "date: " + p.Published,
			Categories: "[\"Archive\"]",
			Body:       p.Markdown,
			Updated:    "updated: " + p.Updated,
		}
		if *dryrun != true {
			out.write()
		}
	}
}
