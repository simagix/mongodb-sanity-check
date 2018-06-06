#! /usr/bin/env python

from requests.auth import HTTPDigestAuth

import getopt
import json
import requests
import sys

user = ''
key = ''
org = ''
search = ''

def processOrg(_link):
    link = _link + "?pretty=true"
    response = requests.get(link, auth=HTTPDigestAuth(user, key))
    doc = json.loads(response.content)
    if search == "invoices":
        link = _link + "/invoices" + "?pretty=true"
        response = requests.get(link, auth=HTTPDigestAuth(user, key))
        doc = json.loads(response.content)
        for result in doc['results']:
            #print("From %s to %s, Status: %s, Billed: %d, Paid: %d" %
            #    (result['created'], result['endDate'], result['status'], result['amountBilledCents'], result['amountPaidCents']))
            for ilink in result['links']:
                response = requests.get(ilink['href'], auth=HTTPDigestAuth(user, key))
                print response.text
        return

    if None == doc.get("links"):
        if doc.get("errorCode") != None:
            print doc["errorCode"]
            print response.text
        return

    for link in doc['links']:
        if search == "groups" and link['rel'] == "http://mms.mongodb.com/groups":
            processLinks(link['href'])
        elif search == "users" and link['rel'] == "http://mms.mongodb.com/users":
            processLinks(link['href'])
        elif search == "teams" and link['rel'] == "http://mms.mongodb.com/teams":
            processLinks(link['href'])
        #else:
        #    print link['href']

def processLink(result, name):
    print 'Retrieving ' + result[name]
    for doc in result['links']:
        link = doc['href']
        link = link + "?pretty=true"
        response = requests.get(link, auth=HTTPDigestAuth(user, key))
        str = response.text
        doc = json.loads(response.content)
        try:
            for link in doc['links']:
                response = requests.get(link['href'] + "?pretty=true", auth=HTTPDigestAuth(user, key))
                print(response.text)
        except:
            print str
            pass
    
def processLinks(link):
    link = link + "?pretty=true"
    response = requests.get(link, auth=HTTPDigestAuth(user, key))
    #print(response.text)
    doc = json.loads(response.content)
    for result in doc['results']:
        if search == "groups":
            processLink(result, 'name')
        elif search == "teams":
            processLink(result, 'name')
        elif search == "users":
            processLink(result, 'username')

def printResults(link):
    link = link + "?pretty=true"
    response = requests.get(link, auth=HTTPDigestAuth(user, key))
    print(response.text)

if __name__ == "__main__":
    options, remainder = getopt.getopt(sys.argv[1:], 'u:k:s:o:', ['user=', 'key=', 'search=', 'org='])
    for opt, arg in options:
        if opt in ('-u', '--user'):
            user = arg
        elif opt in ('-k', '--key'):
            key = arg
        elif opt in ('-s', '--search'):
            search = arg
        elif opt in ('-o', '--org'):
            org = arg

    if search not in ["invoices", "groups", "teams", "users"]:
        print "Invalid search: " + search
        print "valid searches are invoices, groups, teams, and users"
        exit()

    url = "https://cloud.mongodb.com/api/atlas/v1.0/orgs/"
    if org != "":
        url =  url + org
        processOrg(url)
        exit()

    response = requests.get(url + "?pretty=true", auth=HTTPDigestAuth(user, key))

    if not response.ok:
        sys.exit(response.text)

    doc = json.loads(response.content)
    for result in doc['results']:
        for link in result['links']:
            processOrg(link['href'])

