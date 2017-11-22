# go-chat-application

> This is a silly little demo app to learn Golang.

A simple web-based chat application that allows multiple users to have a real-time conversation right in their web browser.
I will start by building a simple web server using the net/http package, which will serve the HTML files. We will then go on to add support for web sockets through which our messages will flow.


This project has three different implementations of profile images: Gravatar, Local files and Google or Github profile image.

It also has authentication with Google and Github providers.

## Run server
```bash
$ ./chat -addr=":3000"

// 2017/10/31 19:46:11 Starting web server on :3000

```