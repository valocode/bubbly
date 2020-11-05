var fs = require('fs')
var Converter = require('openapi-to-postmanv2')
var swagger = fs.readFileSync('../docs/swagger.json', {encoding: 'UTF8'});
swagger = JSON.parse(swagger)
swagger.openapi = "3.0"

Converter.convert({type: "json", data: swagger}, {} ,(err, result) => {
    var collection = JSON.stringify(result.output[0].data, null, 4)
    fs.writeFileSync("postman_collection.json", collection)
})