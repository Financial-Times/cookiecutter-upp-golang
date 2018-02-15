var hooks = require('hooks');
var http = require('http');
var fs = require('fs');

const defaultFixtures = './_ft/ersatz-fixtures.yml';

hooks.beforeAll(function(t, done) {
   if(!fs.existsSync(defaultFixtures)){
      console.log('No fixtures found, skipping hook.');
      done();
      return;
   }

   var contents = fs.readFileSync(defaultFixtures, 'utf8');

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
