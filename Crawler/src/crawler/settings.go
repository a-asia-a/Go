package main

import (
	"time"
	"regexp"
	"crypto/tls"
	"net/http"
	"github.com/temoto/robotstxt"
	"sync"
)

var (
	DefaultClient 			    = http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	DefaultRobots			    = map[string](*robotstxt.Group){}
	DefaultUserAgent      string        = "Mateusz"
	DefaultRobotUserAgent string        = "MateuszRob"
	DefaultInitialWeb     string        = "https://www.ox.ac.uk"
	DefaultDomainSelector 		    = regexp.MustCompile("https?://www\\.ox\\.ac\\.(uk)")
	DefaultNumWorkers     int           = 5
	DefaultNumBuffers     int           = 4
	DefaultCrawlDelay     time.Duration = 30 * time.Second
	DefaultIdleTTL        time.Duration = 30 * time.Second
	DefaultDone			    = make(chan struct{})
	input			            = make(chan Graph, 1)
	RobotsLock sync.RWMutex
)