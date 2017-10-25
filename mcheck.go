package main

import (
    "flag"
    "fmt"
    "log"
    "strconv"
    "time"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

type Robot struct {
    Name string
    Visited int
    Updated time.Time
}

func connectMongo(mongoURI string, total int, thread int) {
    fmt.Println(mongoURI)
    session, err := mgo.Dial("localhost")
    if err != nil {
            panic(err)
    }
    defer session.Close()

    // Optional. Switch the session to a monotonic behavior.
    session.SetMode(mgo.Monotonic, true)

    c := session.DB("mcheck").C("robots")
    err = c.Insert(&Robot{"Robot", 0, time.Now()})
    if err != nil {
        log.Fatal(err)
    }

    bnum := thread * 100000
    for {
	    start := time.Now()
	    for i := bnum; i < (bnum+total); i++ {
            robot := "Robot-" + strconv.Itoa(i)
	        err = c.Insert(&Robot{robot, 0, time.Now()})
	        if err != nil {
	            log.Fatal(err)
	        }
	    }
	    elapsed := time.Since(start)
	    log.Printf("%d inserts took %s", total, elapsed)

        result := Robot{}
	    start = time.Now()
	    for i := bnum; i < (bnum+total); i++ {
            robot := "Robot-" + strconv.Itoa(i)
            err = c.Find(bson.M{"name": robot}).One(&result)
	        if err != nil {
	            log.Fatal(err)
	        }
	    }
	    elapsed = time.Since(start)
	    log.Printf("%d finds took %s", total, elapsed)

	    start = time.Now()
	    for i := bnum; i < (bnum+total); i++ {
            robot := "Robot-" + strconv.Itoa(i)
            err = c.Find(bson.M{"name": robot}).One(&result)
	        if err != nil {
	            log.Fatal(err)
	        }
            change := bson.M{"$inc": bson.M{"visited": 1}}
            err = c.Update(bson.M{"name": robot}, change)
	        if err != nil {
	            log.Fatal(err)
	        }
	    }
	    elapsed = time.Since(start)
	    log.Printf("%d updates took %s", total, elapsed)
        fmt.Println("")

        bnum = bnum + total
        time.Sleep(time.Millisecond * 100)
    }

    result := Robot{}
    err = c.Find(bson.M{"name": "Robot"}).One(&result)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Updated:", result.Updated)
    fmt.Println("=== END ===")
}

func main() {
    total := flag.Int("total", 1000, "total ops in a batch")
    threads := flag.Int("t", 1, "number of threads")
    mongoURI := flag.String("mongoURI", "mongodb://localhost", "MongoDB URI")
    flag.Parse()

    fmt.Println("threads:", *threads)
    fmt.Println("MongoDB URI:", *mongoURI)

    for i := 0; i < *threads; i++ {
        go connectMongo(*mongoURI, *total, i)
    }

    var input string
    fmt.Scanln(&input)
}

