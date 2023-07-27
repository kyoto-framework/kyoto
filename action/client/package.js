const fs = require("fs");

var data;

try {
  data = fs.readFileSync("dist/client.js", "utf8");
  data =
    'package action\n\nvar client = "<script>' +
    data.toString().replaceAll('"', '\\"').split("\n").join("") +
    '</script>"\n';
} catch (err) {
  console.error(err);
}
try {
  fs.writeFileSync("../client.go", data);
} catch (err) {
  console.error(err);
}