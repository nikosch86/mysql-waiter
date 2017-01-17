var mysql = require("mysql")
var http = require("http")

var server_port = 3000
if (typeof process.env.SERVER_PORT != "undefined") {
  server_port = process.env.SERVER_PORT
}

function tryDb() {
  var con = mysql.createConnection({
    host: process.env.MYSQL_HOST,
    user: process.env.MYSQL_USER,
    password: process.env.MYSQL_PASSWORD,
    database: process.env.MYSQL_DB
  });

  con.connect(function(err){
    if(err){
      switch (err.code) {
        case "ER_BAD_DB_ERROR":
          console.log(Date.now()+" Unknown Database: "+process.env.MYSQL_DB)
          break
        case "ECONNREFUSED":
          console.log(Date.now()+" Server not running yet: "+process.env.MYSQL_HOST)
          break
        default:
          console.log(err)
      }

      con.end();
      setTimeout(tryDb, 500)
    } else {
      console.log(Date.now()+" Connection established, database selected")
      con.query("SELECT 1 FROM `"+process.env.MYSQL_TABLE+"` LIMIT 1", function(err, results, fields) {
        if (err) {
          console.log(err)
          con.end()
          setTimeout(tryDb, 500)
        } else {
          console.log(Date.now()+' table found')
          var server = http.createServer(function(req, res) {
            res.write("NOOP").end()
          })
          server.listen(server_port)
          console.log(Date.now()+" server listening on port "+server_port)
        }
      })
    }
  });
}
tryDb()
