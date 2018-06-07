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
Monthly invoices are written to individual CSV file.

```
mongo_atlas_api.py --user $ATLAS_USER --key $ATLAS_KEY \
    --search invoices [--org $ATLAS_ORG]
```
