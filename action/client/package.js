const fs = require('fs')

// Read arguments
const args = {
    r: process.argv[2],
    w: process.argv[3],
    pkg: process.argv[4],
}

// Read client file
const client = fs.readFileSync(args.r, 'utf-8')

// Pack client into Go package file
const package = `package ${args.pkg}

var client = \`${client.replaceAll('`', '` + "`" + `')}\`
`

// Write result
fs.writeFileSync(args.w, package)
