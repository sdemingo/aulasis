# Aulasis


Aulasis is a little CMS to manage a students group workflow in the
classroom.  It is written in Go and this code runs over Linux, Windows
and Mac after compiliing.  You can use Aulasis to publish the diaries
activities using [Markdown](http://es.wikipedia.org/wiki/Markdown)
syntax and they send you the assignments throught a web form.  Aulasis
uses your filesystem to store the information, and it has been
designed to run on a USB memory.

## For users

Go the the [releases
page](https://github.com/sdemingo/aulasis/releases) and download the
latest version package. After uncompress it, run the executable file
for your operating system. Now, Aulasis is running on 9090 port of
your machine. Type  http://localhost:9090 in your browser and enjoy
it.

You can run aulasis with the folowings flags:

```
  -p	Service port
  -d	Data directory (resources and course information)
```

## For developers

Aulasis has been compiled and testing with Go 1.2. If you want to
compile aulasis you must be know that the code has the followings
dependencies:

* [Blackfriday](http://github.com/russross/blackfriday): A great markdown parser

Get them with `go get` and clone it.
