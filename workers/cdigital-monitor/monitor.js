'use strict';


/**
 * Load files
 */
var fetch = require('node-fetch');
var extend = require('util')._extend;

/**
 * Module exports.
 */
 module.exports = monitor;

 /**
  * Main object
  */

 function cvsToJson(string,divisor,headers){
   //Primero por linea
   var lineas = string.split('\n');
   lineas = lineas.filter(function(x){ return x[0]!=null})
   //Despues por campo
   var arrJsonObject = lineas.map(function(linea){
     linea = linea.split(divisor);
     var result = {};
     for (var i = 0; i < linea.length; i++) {
       result[headers[i]] = linea[i];
     }
     return result;
   })
   return arrJsonObject;
 }

 function formatDate(urlDate){
   if(urlDate.fecha == ''){
     var m = new Date();
     urlDate.fecha = m.getUTCFullYear() +''+ (m.getUTCMonth()+1) +''+ m.getUTCDate();
   }
   return urlDate;
 }

 function urlComposer(options){
   if (options.sandbox == true){
     var url = options.urlSandbox
   } else {
      var url = options.urlProduction;
   }
   if(options.time){
     options.urlDate = formatDate(options.urlDate);
     return url + '?control=' + options.hash + '&fecha=' + options.urlDate.fecha +
            '&hour1=' + options.urlDate.hour1 + '&min1=' + options.urlDate.min1 +
            '&hour2=' + options.urlDate.hour2 + '&min2=' + options.urlDate.min2;
   }
   else {
     return url + '?control=' + options.hash;
   }
 }

 var monitor = {
   /* Customizable options */
  options: {
    sandbox: true,
    time: true,
    urlProduction: "https://www.cuentadigital.com/exportacion.php",
    urlSandbox:"https://www.cuentadigital.com/exportacionsandbox.php",
    urlDate: {
      fecha:'',
      hour1:'00',
      min1:'00',
      hour2:'23',
      min2:'59'
    },
    hash:"",
    headers: [
      'fechaDeCobro',
      'horarioDeOperacion',
      'monto',
      'montoNetoRecibido',
      'comision',
      'codigoDeBarraCobrado',
      'suReferencia',
      'medioDePagoUsado',
      'codigoInternoDeLaOperacion',
      'primerCobroDelArchivo'],
      onData : function(data){
        data.time = new Date();
        console.log(data);
        return;
      }
  },
  pullPay: function(req,res,next){
      fetch(urlComposer(monitor.options)).then(function(result) {
        return result.text()
       })
      .then(function(body) {
        if (body.indexOf("<center>") < 0){
          var data = {data: cvsToJson(body,'|',monitor.options.headers)};
        } else {
          var data = {error: 'Error de configuración en la herramienta'}
        }

        monitor.options.onData(data)
        res.send(data)
      });
  },
  init: function(opts) {
    monitor.options = extend(monitor.options, opts);
  }
 }
 module.exports = monitor;