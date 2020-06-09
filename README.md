# Tweetcheck

Requirements: Chrome web browser installed.

## Build binary 
Locally: 
`go build tweetcheck.go`

For specific system architecture (see https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04 for list of architectures):

`env GOOS=target-OS GOARCH=target-architecture go build tweetcheck.go`



## Running

Run by specifying `TWITTER_USERNAME` and `TWITTER_PASSWORD` environment variables, e.g.:

```TWITTER_USERNAME=your_username TWITTER_PASSWORD=your_password ./tweetcheck```

Running the uncompiled tweetcheck.go file is very similar:

```TWITTER_USERNAME=your_username TWITTER_PASSWORD=your_password go run tweetcheck.go```

This will output results like the following:

```Navigating to Twitter..
Logging in to Twitter..
Retrieving message and notification data..

You have 20+ Notifications and 0 Messages.
```

If you have issues with timeouts, try setting TWEETCHECK_TIMEOUT to a value higher than 30.



