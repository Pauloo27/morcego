const mysql = require('mysql');
const express = require('express');

const app = express()

var con = mysql.createConnection({
  host: "localhost",
  user: process.env.DB_USERNAME,
  password: process.env.DB_PASSWORD,
  database: process.env.DB_NAME,
});

con.connect(function(err) {
  if (err) throw err;
  console.log("Connected!");
});

app.get("/", (req, res) => {
  res.send(`
  <pre>
      /| ________________
O|===|* >________________> 1337 - EXPRESS SQLI
      \\|
  </pre>
  `);
});

app.get("/login", (req, res) => {
  res.send(`
  <form method="POST" action="/login">
    <input placeholder="Username" name="username"/>
    <input type="password" placeholder="Password" name="password"/>
    <input type="submit"/>
  </form>
  `)
})

app.get("/u/:name", (req, res) => {
  console.log(req.params.name)
  con.query(`SELECT id, username FROM users WHERE username = '${req.params.name}';`, (err, result) => {
    res.send(result)
  })
})

app.get("/users/:id", (req, res) => {
  con.query(`SELECT id, username FROM users WHERE id = ${req.params.id};`, (err, result) => {
    res.send(result)
  })
})

app.use(express.urlencoded())

app.get("/user", (req, res) => {
  if (req.query.test !== "2") {
    res.send("Invalid test")
    return
  }
  con.query(`SELECT username FROM users WHERE id = ${req.query.id};`, (err, result) => {
    res.send(result)
  })
});

app.post("/login", (req, res) => {
  console.log(req.body)
  con.query(`SELECT password from users WHERE username = '${req.body.username}'`, (err, result) => {
    if(result == undefined || result.length == 0) { 
      res.send("Invalid user");
    } else if(result[0].password == req.body.password) {
      res.send("Ok")
    } else {
      res.send("Invalid password")
    }
  })
})

console.log("Starting server");
app.listen(3032);
