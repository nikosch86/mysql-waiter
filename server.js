var mysql = require("mysql")
var http = require("http")

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
          console.log("Unknown Database")
          break
        case "ECONNREFUSED":
          console.log("Server not running yet")
          break
        default:
          console.log(err)
      }

      con.end();
      setTimeout(tryDb, 500)
    } else {
      con.query('SHOW TABLES LIKE \''+process.env.MYSQL_TABLE+'\'', function(err, results, fields) {
        if (err) {
          console.log(err)
          con.end()
          setTimeout(tryDb, 500)
        } else {
          console.log('Connection established')
          var server = http.createServer(function(req, res) {
            res.write("NOOP").end()
          })
          server.listen(3000)
        }
      })
    }
  });
}
tryDb()
