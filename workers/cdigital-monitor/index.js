'use strict';
var monitor = require("./monitor");
var express = require('express');
var app = express();

//CONFIGURACION DE CUENTA DIGITAL Y DEL CALLBACK
monitor.init({
  hash: process.env.CDIGITAL || "43414e465afce7510a5c506644c9c35a",
  sandbox: false,
  time:true
});

app.use('/monitor',monitor.pullPay);

//INICIO EL SERVER
var server = app.listen(process.env.PORT || 3000, function () {
  var host = server.address().address;
  var port = server.address().port;
  console.log('Example app listening at http://%s:%s', host, port);
});