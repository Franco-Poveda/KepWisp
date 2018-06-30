'use strict'

const morgan = require('morgan');
const express = require('express');
const bodyParser = require('body-parser')
const http = require('http');
const CSV = require('csv-string');
var amqp = require('amqplib');

const app = express()
app.use(morgan('combined'))
let outch;
// parse an HTML body into a string
app.use(bodyParser.text({ type: 'text/csv' }))
app.post('/', (req, res) => {
  const arr = CSV.parse(req.body);

  console.log(arr)
  res.send('');
  outch.publish('topic_logs', 'info', Buffer.from(JSON.stringify(arr)));
  console.log(" [x] Sent %s:'%s'", 'info', arr);
});
  app.get('/', (req, res) => res.send('Hello!'))

const server = http.createServer(app)

let ip = process.env.IP || "0.0.0.0";
let port = process.env.PORT || 3003;

amqp.connect('amqp://localhost').then(function(conn) {
  conn.on('error', err => {
    console.log('error', `[AMQP]: [${err.message}]`);
});
conn.on('close', () => {
    console.log('info', '[AMQP]: connection closed');
});
  return conn.createChannel().then(function(ch) {
    var ex = 'topic_logs';
    var ok = ch.assertExchange(ex, 'topic', {durable: true});
    outch = ch;
    return ok.then(function() {
      server.listen(port, ip, () => {
        console.log('[SERVER] WebHook API listening on http://%s:%d', ip, port)
      });
    });
  })
}).catch(console.log);

