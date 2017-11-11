#! /usr/bin/env python

import json
import re
import sys

def convert(str):
    data = json.loads(str)
    stats = data["executionStats"]
    stages = stats["executionStages"]
    if data["queryPlanner"]["winningPlan"].get("shards", None):
        d = data["queryPlanner"]["winningPlan"]["shards"][0]
        print "version : %s" % d["serverInfo"]["version"]
        print "-- FILTER --"
        print json.dumps(d["parsedQuery"], indent=3, sort_keys=False)
        if d["winningPlan"]["inputStage"].get("sortPattern", None):
            print "-- SORT BY --"
            print json.dumps(d["winningPlan"]["inputStage"]["sortPattern"], indent=3, sort_keys=True)
    else:
        print "version : %s" % data["serverInfo"]["version"]
        print "-- FILTER --"
        print json.dumps(data["queryPlanner"]["parsedQuery"], indent=3, sort_keys=False)
        wplan = data["queryPlanner"]["winningPlan"]
        if wplan.get("inputStage", None):
            d = wplan["inputStage"]
            if d.get("sortPattern", None):
                print "-- SORT BY --"
                print json.dumps(d["sortPattern"], indent=3, sort_keys=True)
    print "\n-- SUMMARY --"
    print "executionTimeMillis : %d" % stats["executionTimeMillis"]
    print "nReturned : %d" % stats["nReturned"]
    print "totalKeysExamined : %d" % stats["totalKeysExamined"]
    print "totalDocsExamined : %d" % stats["totalDocsExamined"]
    print "\n-- STAGES --"
    printInputStage(stages, 0)   

def printInputStage(stats, level):
    while stats.get("stage", None):
        tabs = getTabs(level)
        print getExtTabs(level) + "%s" % stats["stage"]
        if stats.get("executionTimeMillisEstimate", None):
            print tabs + "executionTimeMillisEstimate : %d" % stats["executionTimeMillisEstimate"]
        else:
            print tabs + "executionTimeMillisEstimate : 0"
        print tabs + "nReturned : %d" % stats["nReturned"]
        if stats.get("keysExamined", None):
            print tabs + "keysExamined : %d" % stats["keysExamined"]
        if stats.get("docsExamined", None):
            print tabs + "docsExamined : %d" % stats["docsExamined"]
        if stats.get("indexName", None):
            print tabs + "index used : %s" % stats["indexName"]

        if stats.get("inputStage", None):
            level = level + 1
            stats = stats["inputStage"]
        elif stats.get("shards", None):
            for sh in stats["shards"]:
                print "(%s)" % sh["shardName"]
                printInputStage(sh["executionStages"], level+1)
            return
        elif stats.get("inputStages", None):
            for wp in stats["inputStages"]:
                printInputStage(wp, level+1)
            return
        else:
            return

def getTabs(level):
    tabs = ""
    if level > 0:
        for i in range(0, level):
            tabs += ".  "
    tabs += "|  - "
    return tabs

def getExtTabs(level):
    tabs = ""
    if level > 0:
        for i in range(0, level-1):
            tabs += ".  "
        tabs = tabs + "+--"
    return tabs

def getSep(level):
    tabs = ""
    for i in range(0, level):
        tabs += ".  "
    tabs += "|"
    return tabs

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print "Usage: %s json_file" % (sys.argv[0])
        exit(1)
    with open(sys.argv[1]) as file:
        str = file.read()
        str = re.sub('ISODate\((.*)\)', r'{ "$date": \1 }', str)
        str = re.sub('Number.*\((.*)\)', '1', str)
        convert(str)
