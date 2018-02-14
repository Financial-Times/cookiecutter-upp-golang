var hooks = require('hooks');
var http = require('http');
var fs = require('fs');

hooks.beforeAll(function(t, done) {
   var contents = fs.readFileSync('./_ft/ersatz-fixtures.yml', 'utf8');

   var options = {
      host: 'localhost',
      port: '9000',
      path: '/__configure',
      method: 'POST',
      headers: {
         'Content-Type': 'application/x-yaml'
      }
   };

   var req = http.request(options, function(res) {
      res.setEncoding('utf8');
   });

   req.write(contents);
   req.end();
   done();
});
