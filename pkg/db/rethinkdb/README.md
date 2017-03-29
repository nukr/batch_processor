
## When slice is negative
mongodb

```JavaScript
db.getCollection('test').update(
  {_id: ObjectId("58d40bb3169306ac2b32d65d")},
  {$push: {"aa.bb.cc": {$each: [11, 22, 33, 44], $slice: -2}}}
)
```

rethinkdb

```JavaScript
r.db('aa').table('aa').getAll("6cd77585-92b4-4895-ad1c-435326986c6d").update(function(row){
  return {
    aa: {
      bb: {
        cc: row('aa')('bb')('cc').add([1, 2, 3, 4]).slice(-3, -1, {rightBound: 'closed'})
      }
    }
  }
})
```

## When slice is positive

```JavaScript
db.getCollection('test').update(
  {_id: ObjectId("58d40bb3169306ac2b32d65d")},
  {$push: {"aa.bb.cc": {$each: [11, 22, 33, 44], $slice: 3}}}
)
```

rethinkdb

```JavaScript
r.db('aa').table('aa').getAll("6cd77585-92b4-4895-ad1c-435326986c6d").update(function(row){
  return {
    aa: {
      bb: {
        cc: row('aa')('bb')('cc').add([1, 2, 3, 4]).slice(0, 3)
      }
    }
  }
})
```
