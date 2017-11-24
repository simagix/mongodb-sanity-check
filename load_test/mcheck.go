package main

import (
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
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

func connectMongo(mongoURI string, tps int, size int, once bool, thread int) {
	fmt.Println(mongoURI)
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
		for i := bnum; i < (bnum + tps); i++ {
			robot := "Robot-" + strconv.Itoa(i)
			err = c.Insert(&Robot{robot, robot, buffer.String(), Statistics{Tasked: 0, Battery: 95, Maint: false}, time.Now()})
			if err != nil {
				log.Fatal(err)
			}
		}
    elapsed := time.Since(start)
		log.Printf("INSERT %d %s size %d", tps, elapsed, size)
    
		if once == true {
      os.Exit(0)
		}

		result := Robot{}
		start = time.Now()
		for i := bnum; i < (bnum + tps); i++ {
			robot := "Robot-" + strconv.Itoa(i)
			err = c.Find(bson.M{"name": robot}).One(&result)
			if err != nil {
				log.Fatal(err)
			}
		}
		elapsed = time.Since(start)
		log.Printf("FIND   %d %s with index {name: 1}", tps, elapsed)

		result = Robot{}
		start = time.Now()
		for i := bnum; i < (bnum + tps); i++ {
			robot := "Robot-" + strconv.Itoa(i)
			err = c.Find(bson.M{"nickname": robot}).One(&result)
			if err != nil {
				log.Fatal(err)
			}
		}
		elapsed = time.Since(start)
		log.Printf("FIND   %d %s without index", tps, elapsed)

		start = time.Now()
		for i := bnum; i < (bnum + tps); i++ {
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
		log.Printf("UPDATE %d %s $inc stats.tasked by 1", tps, elapsed)

		start = time.Now()
		for i := bnum; i < (bnum + tps); i++ {
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
		log.Printf("UPDATE %d %s $set descr string size of %d", tps, elapsed, size)

		fmt.Println("")

		bnum = bnum + tps
		time.Sleep(time.Millisecond * 100)
	}
}

func cleanup(mongoURI string) {
	fmt.Println("cleanup", mongoURI)
	session, err := mgo.Dial(mongoURI)
	defer session.Close()
	err = session.DB(mcheck).DropDatabase()
	if err != nil {
		panic(err)
	}
}

func createIndex(mongoURI string) {
	fmt.Println("createIndex", mongoURI)
	session, err := mgo.Dial(mongoURI)
	defer session.Close()
	c := session.DB(mcheck).C("robots")
	index := mgo.Index{
		Key: []string{"name"},
	}

	err = c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

func main() {
	tps := flag.Int("tps", 100, "transactions per second")
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
	mcheck = fmt.Sprintf("%s%X", head, buf)
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
		go connectMongo(*mongoURI, *tps, *size, *seed, i)
	}

	var input string
	fmt.Scanln(&input)
}
