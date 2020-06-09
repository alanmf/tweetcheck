package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	// Set opts so we can run in a viewable mode if we need to for debugging
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
	)

	// create allocaor
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create context
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	var notifications string
	var messages string

	// Set the default timeout for running the whole operation to 30 seconds
	// Override with TWEETCHECK_TIMEOUT env var
	tweetcheck_timeout := 30
	tweetcheck_timeout_env := os.Getenv("TWEETCHECK_TIMEOUT")
	if len(tweetcheck_timeout_env) != 0 {
		i, err := strconv.Atoi(tweetcheck_timeout_env)
		tweetcheck_timeout = i
		if err != nil {
			log.Fatal(err)
		}
	}

	tasks := getTasks(
		`https://twitter.com/login`,
		`//*[@id="react-root"]/div/div/div[2]/main/div/div/form/div/div[1]/label/div/div[2]/div/input`,
		`//*[@id="react-root"]/div/div/div[2]/main/div/div/form/div/div[2]/label/div/div[2]/div/input`,
		os.Getenv("TWITTER_USERNAME"),
		os.Getenv("TWITTER_PASSWORD"),
		&notifications,
		&messages)

	err := chromedp.Run(
		ctx,
		RunWithTimeOut(&ctx, time.Duration(tweetcheck_timeout), tasks),
	)

	//Check error and quit if we failed.
	if err != nil {
		if err.Error() == `context deadline exceeded` {
			fmt.Printf("Unable to complete in %v seconds, set a higher timeout with env var TWEETCHECK_TIMEOUT and check user credentials.\n", tweetcheck_timeout)
		}
		log.Fatal(err)
	}

	// Display appropriate data based on notifications or messages numbers
	if messages == `` {
		messages = `0`
	}
	if notifications == `` {
		notifications = `0`
	}

	if notifications == `0` && messages == `0` {
		fmt.Println("You don't have any Notifications nor Messages.")
	} else {
		fmt.Printf("\nYou have %s Notifications and %s Messages.\n\n", strings.TrimSpace(notifications), strings.TrimSpace(messages))
	}
}

func RunWithTimeOut(ctx *context.Context, timeout time.Duration, tasks chromedp.Tasks) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		timeoutContext, cancel := context.WithTimeout(ctx, timeout*time.Second)
		defer cancel()
		return tasks.Do(timeoutContext)
	}
}

func getTasks(urlstr, usernameSel, passwordSel, username string, password string, notifications *string, messages *string) chromedp.Tasks {
	//Execute chromedp tasks
	/*
		1. Navigate to page
		2. Wait for load
		3. Input username
		4. Input password
		5. Click submit (using submit on the form fails for some reason)
		6. Wait for next page load
		7. Read notifications and messages based on xml paths
	*/
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := fmt.Println("Navigating to Twitter..")
			return err
		}),
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible(usernameSel),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := fmt.Println("Logging in to Twitter..")
			return err
		}),
		chromedp.SendKeys(usernameSel, username),
		chromedp.SendKeys(passwordSel, password),
		chromedp.Click(`//*[@id="react-root"]/div/div/div[2]/main/div/div/form/div/div[3]/div/div`, chromedp.NodeVisible),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := fmt.Println("Retrieving message and notification data..")
			return err
		}),
		chromedp.WaitReady(`//*[@id="react-root"]/div/div/div[2]/main/div/div/div/div/div/div[4]/div/div/section/div`),
		chromedp.Text(`//*[@id="react-root"]/div/div/div[2]/header/div/div/div/div[1]/div[2]/nav/a[3]`, notifications),
		chromedp.Text(`//*[@id="react-root"]/div/div/div[2]/header/div/div/div/div[1]/div[2]/nav/a[4]`, messages),
	}
}
