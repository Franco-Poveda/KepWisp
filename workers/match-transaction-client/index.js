var amqp = require('amqplib');
var pool = require('./db');
var mysql = require('mysql');
var moment = require('moment');


amqp.connect('amqp://localhost').then(function (conn) {
    process.once('SIGINT', function () { conn.close(); });
    return conn.createChannel().then(function (ch) {

        var ok = ch.assertQueue('q1', { durable: true });


        ok = ok.then(function () {

            return ch.consume('q1', match, { noAck: false });
        });
        return ok.then(function () {
            console.log(' [*] Waiting for logs. To exit press CTRL+C.');
        });
        function match(msg) {
            var tarr = JSON.parse(msg.content.toString());
            tarr.pop();
            pool.getConnection(function (err, connection) {
                if (err) throw err; // not connected!
                connection.beginTransaction(function (err) {
                    if (err) { throw err; }
                    var qvalues = tarr.map(t => {
                        var set = t.slice(0, -2);
                        var time = set.splice(2, 1);
                        set[1] = moment(set[1] + time[0], 'DDMMYYYYkkmmss').format('YYYY-MM-DD kk:mm:ss');
                        return set;
                    });
                    var query = "INSERT INTO `Transactions` (`type`,`tdate`,`amount`,`barcode`,`ref`,`method`,`cduid`) VALUES ";
                    var pvalues = mysql.escape(qvalues)
                    query = query + pvalues;//[1,"2018-06-20 12:00:00",100, "01909124354188", "ejemplo referencia 1 cliente dni 11111111", "PagoFacilPRUEBA", "66a3bfdfff68f605ce460b7f71ef79a642f10d9ae112077e0dd2dda895b3f5cc"]);
                    console.log(query);
                    connection.query(query, function (error, results, fields) {
                        if (error) {
                            return connection.rollback(function () {
                                throw error;
                            });
                        }
                        connection.commit(function (err) {
                            if (err) {
                                return connection.rollback(function () {
                                    throw err;
                                });
                            }
                            console.log('success!');
                            ch.sendToQueue('q2', Buffer.from(JSON.stringify(qvalues)));
                            ch.ack(msg);
                        });
                    });
                });
            });
        }
    });
}).catch(console.warn);
