package main

import (
	"database/sql"
	"fmt"
	"github.com/yookoala/buflog"
	"github.com/yookoala/crawler"
	"github.com/yookoala/crawler/sqlcache"
	"time"
)

// gets all cached result and display
func example5(host string, db *sql.DB, log *buflog.Logger) (resp *crawler.Response, err error) {

	log.Print("# Get old cache by context time")

	url := host + "/example/5"
	c := sqlcache.New(*dbdriver, db)
	f := crawler.NewFetcher(c)

	// render context time
	d, err := time.ParseDuration("24h")
	if err != nil {
		return
	}
	t, err := time.Parse(time.RFC822Z, "01 Apr 10 00:00 +0800")
	if err != nil {
		return
	}
	st := t // beginning time
	l := 10 // test scale limit

	for i := 1; i <= l; i++ {
		ctx := crawler.Context{
			Str:  "example/5",
			Time: t,
		}
		_, err = f.Get(url, ctx)
		if err != nil {
			return
		}
		t = t.Add(d)
	}

	t = st // reset to beginning time
	var rs crawler.ResponseColl
	for i := 0; i < l; i++ {

		// search the existing url
		rs, err = c.
			FindAt(t).
			In("example/5").
			GetAll()
		if err != nil {
			return
		}

		// load response into response slice
		resps := make([]crawler.Response, 0)
		for rs.Next() {
			resp, err := rs.Get()
			if err != nil {
				log.Fatal("Error getting next response")
			}
			resps = append(resps, *resp)
		}

		// check number of records
		if len(resps) == 0 {
			err = fmt.Errorf("No cache found at %s",
				t.Format("2006-01-02"))
		} else if len(resps) > 1 {
			err = fmt.Errorf("Too many cache found at %s. "+
				"%d found while expecting 1",
				t.Format("2006-01-02"),
				len(resps))
		}

		// log record found
		log.Printf("[#%d] (%s) Body: \"%s\"", i,
			resps[0].ContextTime.Format("2006-01-02"),
			string(resps[0].Body))

		t = t.Add(d)
	}
	return
}
