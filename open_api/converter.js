var fs = require('fs');
Swagger2Postman = require("swagger2-to-postman");

var collectionPath = "./open_api/postman_collection.json";

converter = new Swagger2Postman();

console.log("Converting OpenApi Docs to a Postman Collection");

// Load Swagger File
var json = fs.readFileSync("../docs/swagger.json")
json = JSON.parse(json)

// Convert Swagger to Postman Json
convertResult = converter.convert(json);

// Stringify the results for writing to a file
var str = JSON.stringify(convertResult, null, 4)

var count = str.split("::").length - 1;
var i;
for (i = 0; i < count; i++) {
    // The converter messes up variables in the GET path, this will replace ::
    // and include {{ variable }} syntax for use in Newman
    var first = str.split("::")[0];
    var second = str.split("::")[1];

    var variable = second.substring(0, second.indexOf("\""));
    variable = "{{" + variable + "}}";
    str = first + variable + second.substring(variable.length -4)
}


fs.writeFile(collectionPath, str, 'utf8', () => {})
console.log("Collection Created at: " + collectionPath);