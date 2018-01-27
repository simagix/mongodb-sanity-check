package main

import (
	"bytes"
	"crypto/rand"
    "crypto/sha1"
	"flag"
	"fmt"
	"log"
	mrand "math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const head = "_MCHECK_"

var mcheck = head

type Statistics struct {
	Tasked  int
	Battery int
	Maint   bool
}

type Robot struct {
	Name     string
	Nickname string
	Descr    string
	Stats    Statistics
	Updated  time.Time
}

type Brand struct {
	Name     string
    Sku      string   
}

func connectMongo(mongoURI string, batch int, size int, once bool, thread int) {
	session, err := mgo.Dial(mongoURI)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(mcheck).C("robots")

	var buffer bytes.Buffer
	for i := 0; i < size/len("Simagix."); i++ {
		buffer.WriteString("Simagix.")
	}

	var newbuf bytes.Buffer
	for i := 0; i < size/len("MongoDB."); i++ {
		newbuf.WriteString("MongoDB.")
	}

	bnum := thread * 100000
	for {
		start := time.Now()
		for i := bnum; i < (bnum + batch); i++ {
			robot := "Robot-" + strconv.Itoa(i)
			num := mrand.Intn(20)
			batt := 100 - num*5
			err = c.Insert(&Robot{robot, robot, buffer.String(), Statistics{Tasked: num, Battery: batt, Maint: (batt > 25)}, time.Now()})
			if err != nil {
				log.Fatal(err)
			}
		}
		elapsed := time.Since(start)
		avg := time.Duration(elapsed.Nanoseconds() / int64(batch))
		log.Printf("INSERT %d %s %s size %d", batch, avg, elapsed, size)

		if once == true {
	        b := session.DB(mcheck).C("brands")
		    for i := bnum; i < (bnum + batch); i++ {
			    robot := "Robot-" + strconv.Itoa(i)
                h := sha1.New()
			    err = b.Insert(&Brand{robot, fmt.Sprintf("%X", h.Sum(nil))})
			if err != nil {
				log.Fatal(err)
			}
            }
			os.Exit(0)
		}

		result := Robot{}
		start = time.Now()
		for i := bnum; i < (bnum + batch); i++ {
			robot := "Robot-" + strconv.Itoa(i)
			err = c.Find(bson.M{"name": robot}).One(&result)
			if err != nil {
				log.Fatal(err)
			}
		}
		elapsed = time.Since(start)
		avg = time.Duration(elapsed.Nanoseconds() / int64(batch))
		log.Printf("FIND   %d %s %s with index {name: 1}", batch, avg, elapsed)

		result = Robot{}
		start = time.Now()
		for i := bnum; i < (bnum + batch); i++ {
			robot := "Robot-" + strconv.Itoa(i)
			err = c.Find(bson.M{"nickname": robot}).One(&result)
			if err != nil {
				log.Fatal(err)
			}
		}
		elapsed = time.Since(start)
		avg = time.Duration(elapsed.Nanoseconds() / int64(batch))
		log.Printf("FIND   %d %s %s without index", batch, avg, elapsed)

		start = time.Now()
		for i := bnum; i < (bnum + batch); i++ {
			robot := "Robot-" + strconv.Itoa(i)
			err = c.Find(bson.M{"name": robot}).One(&result)
			if err != nil {
				log.Fatal(err)
			}
			change := bson.M{"$inc": bson.M{"stats.tasked": 1}}
			err = c.Update(bson.M{"name": robot}, change)
			if err != nil {
				log.Fatal(err)
			}
		}
		elapsed = time.Since(start)
		avg = time.Duration(elapsed.Nanoseconds() / int64(batch))
		log.Printf("UPDATE %d %s %s $inc stats.tasked by 1", batch, avg, elapsed)

		start = time.Now()
		for i := bnum; i < (bnum + batch); i++ {
			robot := "Robot-" + strconv.Itoa(i)
			err = c.Find(bson.M{"name": robot}).One(&result)
			if err != nil {
				log.Fatal(err)
			}
			change := bson.M{"$set": bson.M{"descr": newbuf.String()}}
			err = c.Update(bson.M{"name": robot}, change)
			if err != nil {
				log.Fatal(err)
			}
		}
		elapsed = time.Since(start)
		avg = time.Duration(elapsed.Nanoseconds() / int64(batch))
		log.Printf("UPDATE %d %s %s $set descr string size of %d", batch, avg, elapsed, size)

		fmt.Println("")

		bnum = bnum + batch
		time.Sleep(time.Millisecond * 100)
	}
}

func cleanup(mongoURI string) {
	fmt.Println("cleanup", mongoURI)
	session, _ := mgo.Dial(mongoURI)
	defer session.Close()
	session.DB(mcheck).DropDatabase()
}

func createIndex(mongoURI string) {
	fmt.Println("createIndex", mongoURI)
	session, _ := mgo.Dial(mongoURI)
	defer session.Close()
	c := session.DB(mcheck).C("robots")
	index := mgo.Index{
		Key: []string{"name"},
	}

	c.EnsureIndex(index)
}

func main() {
	batch := flag.Int("batch", 512, "ops per batch")
	threads := flag.Int("t", 1, "number of threads")
	mongoURI := flag.String("mongoURI", "mongodb://localhost", "MongoDB URI")
	size := flag.Int("size", 1024, "document size")
	seed := flag.Bool("seed", false, "seed a database for demo")
	flag.Parse()
	fmt.Println("threads:", *threads)
	fmt.Println("MongoDB URI:", *mongoURI)
	fmt.Println("seed:", *seed)

	buf := make([]byte, 4)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}

    if *seed == false {
	    mcheck = fmt.Sprintf("%s%X", head, buf)
    }
	fmt.Println("Populate data under database", mcheck)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup(*mongoURI)
		os.Exit(1)
	}()

	createIndex(*mongoURI)
	for i := 0; i < *threads; i++ {
		go connectMongo(*mongoURI, *batch, *size, *seed, i)
	}

	var input string
	fmt.Println("Ctrl-C to quit...")
	fmt.Scanln(&input)
}
