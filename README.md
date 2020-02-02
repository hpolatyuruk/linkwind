# turkdev
turkdev

# MakeFile

We use gcflags and asmflags. Because we do not want
error's absolute path.

Also we assign environment variable in build time. We read ".env" file, take variables and assign in build time.
Because of this, we added LDFLAGS.
For use makefile command. You can run this command for example development : 

```sh
$ make build file=".env.dev"
```

For staging or production build, you should create prod and staging environment files and give parameter that file names.
But for use this command, you can install mingw. Check this site and install : 

https://www.ics.uci.edu/~pattis/common/handouts/mingweclipse/mingw.html

The go to "C:\MinGW\bin" directory and change "mingw-32 make" to "make" . That's it. 

For more information about Makefile : 

https://tutorialedge.net/golang/makefiles-for-go-developers/