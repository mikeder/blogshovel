# blogshovel
tornado -> hugo post shovel written in go

## Install & Build

```bash
$ go get github.com/mikeder/blogshovel
$ cd $GOPATH/src/github.com/mikeder/blogshovel
$ go install
```

### Help
```bash
meder@debian:~$ ./blogshovel
blogshovel: mysql -> archive(markdown)

You must provide a -dbconstring
Usage of ./blogshovel:
  -dbconnstring string
        database connection string (default "user:password@tcp(host.domain:3306)/database")
  -dryrun
        skip writing to file
  -outdir string
        output files to this directory (default "/home/meder/archive/")
```
