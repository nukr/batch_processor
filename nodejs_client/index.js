var grpc = require('grpc')
var path = require('path')

var protoDescriptor = grpc.load(path.join(__dirname, '../pb/batch.proto'))

var client = new protoDescriptor.BatchService("35.185.146.34:33333", grpc.credentials.createInsecure())

var query = {
  selector: Buffer.from(`{id: 1}`),
  document: Buffer.from(`
  {
    "$push": {
      "field1": {
        "$each": [4, 5, 6, 7, 8, 9]
      },
      "field2": {
        "$each": [{"aa": 11, "bb": 22}, {"aa": 33, "bb": 44}],
        "$slice": 2,
        "$sort": {"aa": -1}
      },
      "field3": "123123",
      "aa.bb.cc": {
        "$each": [100000]
      }
    }
  }
  `)
}

client.update(query, function(err, result) {
  if (err) {
    return console.error(err)
  }
  console.log(result)
})
