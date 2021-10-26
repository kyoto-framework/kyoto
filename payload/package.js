const fs = require("fs");

var data;

try {
  data = fs.readFileSync("dist/dynamics.js", "utf8");
  data =
    'package kyoto\n\nvar ssaclient = "<script>' +
    data.toString().replaceAll('"', '\\"').split("\n").join("") +
    '</script>"\n';
} catch (err) {
  console.error(err);
}
try {
  fs.writeFileSync("../ext.ssa.client.go", data);
} catch (err) {
  console.error(err);
}
