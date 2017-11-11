<h3>MongoDB explain("executionStats") Summary</h3>
Use *explain_summary.py* to display summary of JSON output from *explain("executionStats") command.

[TOC]

### 1. Usage

```
Usage: python explain_summary.py <json_file>
```
#### 1.1. Examples

```
python explain_summary.py explain-single.json
python explain_summary.py explain-shard.json
```

#### 1.2. Sample Outputs
```
python explain_summary.py explain-single.json

version : 3.4.9
-- FILTER --
{
   "$or": [
      {
         "stats": {
            "$elemMatch": {
               "$and": [
                  {
                     "src": {
                        "$eq": "cpu"
                     }
                  },
                  {
                     "pct": {
                        "$gte": 80
                     }
                  }
               ]
            }
         }
      },
      {
         "stats": {
            "$elemMatch": {
               "$and": [
                  {
                     "src": {
                        "$eq": "mem"
                     }
                  },
                  {
                     "pct": {
                        "$gte": 90
                     }
                  }
               ]
            }
         }
      }
   ]
}
-- SORT BY --
{
   "dt": -1
}

-- SUMMARY --
executionTimeMillis : 2
nReturned : 396
totalKeysExamined : 420
totalDocsExamined : 420

-- STAGES --
 stage SUBPLAN
|  - executionTimeMillisEstimate : 0
|  - nReturned : 396
+--stage SORT
.  |  - executionTimeMillisEstimate : 0
.  |  - nReturned : 396
.  +--stage SORT_KEY_GENERATOR
.  .  |  - executionTimeMillisEstimate : 0
.  .  |  - nReturned : 396
.  .  +--stage OR
.  .  .  |  - executionTimeMillisEstimate : 0
.  .  .  |  - nReturned : 396
.  .  .  +--stage FETCH
.  .  .  .  |  - executionTimeMillisEstimate : 0
.  .  .  .  |  - nReturned : 271
.  .  .  .  |  - docsExamined : 271
.  .  .  .  +--stage IXSCAN
.  .  .  .  .  |  - executionTimeMillisEstimate : 0
.  .  .  .  .  |  - nReturned : 271
.  .  .  .  .  |  - keysExamined : 271
.  .  .  .  .  |  - index used : stats.src_1_stats.pct_1_dt_-1
.  .  .  +--stage FETCH
.  .  .  .  |  - executionTimeMillisEstimate : 0
.  .  .  .  |  - nReturned : 149
.  .  .  .  |  - docsExamined : 149
.  .  .  .  +--stage IXSCAN
.  .  .  .  .  |  - executionTimeMillisEstimate : 0
.  .  .  .  .  |  - nReturned : 149
.  .  .  .  .  |  - keysExamined : 149
.  .  .  .  .  |  - index used : stats.src_1_stats.pct_1_dt_-1
```

### 2. References
- [Explain Results](https://docs.mongodb.com/manual/reference/explain-results/)
- [Analyze Query Performance](https://docs.mongodb.com/manual/tutorial/analyze-query-plan/)