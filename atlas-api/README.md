# MongoDB Atlas API
A simple Python script to retrieve info from MongoDB Atlas.

## Groups
### List Groups (Projects) Summary

```
mongo_atlas_api.py --user $ATLAS_USER --key $ATLAS_KEY \
    --search groups [--org $ATLAS_ORG]
```
### List Groups (Projects) Details

```
mongo_atlas_api.py --user $ATLAS_USER --key $ATLAS_KEY \
    --search groups [--org $ATLAS_ORG] --verbose
```

## Teams

```
mongo_atlas_api.py --user $ATLAS_USER --key $ATLAS_KEY \
    --search teams [--org $ATLAS_ORG]
```

## Users

### List Users Summary

```
mongo_atlas_api.py --user $ATLAS_USER --key $ATLAS_KEY \
    --search users [--org $ATLAS_ORG]
```

### List Users Details

```
mongo_atlas_api.py --user $ATLAS_USER --key $ATLAS_KEY \
    --search users [--org $ATLAS_ORG] --verbose
```

## Invoices
### Download Invoices
Monthly invoices are written to individual CSV file.

```
mongo_atlas_api.py --user $ATLAS_USER --key $ATLAS_KEY \
    --search invoices [--org $ATLAS_ORG]
```

### Load Data to Mongo

```
rm -rf idata; mkdir -p idata/db
mongod --dbpath idata/db --logpath idata/mongod.log --port 30097 --fork

for file in $(ls *.csv)
do
    mongoimport --port 30097 -d example -c invoices --file=$file \
        --type=csv --headerline
done
```

### Examples
To get only desired fields and project, you can use `mongoexport` to retieve data from mongo server.  First create a file with all fields you need, for example, a file named fields,

```
Date
Project
Project ID
SKU
Region
Cluster
Unit
Unit Price
Quantity
Discount ID
Discount Percent
Discount Total
Amount
Tax
```

Next, use `mongoexport` command.  In the example below, we will query all docs with project ID **5bc2ccbd4f65817a1671f0e7**.

#### Get Data of a Project
```
mongoexport --port 30097 -d example -c invoices --type=csv --fieldFile=fields \
    --query='{"Project ID":"5bc2ccbd4f65817a1671f0e7"}' \
    > project-5bc2ccbd4f65817a1671f0e7.csv
```

#### Get Data of a Project for a Month
Add additional query parameter to only retrieve data in the month of April, 2018.

```
mongoexport --port 30097 -d example -c invoices --type=csv --fieldFile=fields \
    --query='{"Project ID":"5bc2ccbd4f65817a1671f0e7","Date":{"$gte":"04/01/2018","$lt":"05/01/2018"}}' \
    > project-5bc2ccbd4f65817a1671f0e7.csv
```
5ab2bbac4e65817f1671e0d7
5bc2ccbd4f65817a1671f0e7

#### Get Data of All Projects to Different Files

```
for pid in $(mongo --quiet --port 30097 example \
    --eval "db.invoices.find({}, {'Project ID': 1}).forEach(function(doc) { print(doc['Project ID']); })"|sort|uniq)
do
    mongoexport --port 30097 -d example -c invoices --type=csv --fieldFile=fields \
        --query="{'Project ID':'$pid','Date':{'\$gte':'04/01/2018','\$lt':'05/01/2018'}}" \
        > project-$pid-2018-04.csv
done
```






