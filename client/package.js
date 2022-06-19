const fs = require("fs");

var data;

try {
  data = fs.readFileSync("dist/dynamics.js", "utf8");
  data =
    'package kyoto\n\nvar ActionClient = "<script>' +
    data.toString().replaceAll('"', '\\"').split("\n").join("") +
    '</script>"\n';
} catch (err) {
  console.error(err);
}
try {
  fs.writeFileSync("../actions.client.go", data);
} catch (err) {
  console.error(err);
}