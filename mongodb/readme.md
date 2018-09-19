## Setup:
      On Linux,
          service start mongodb
      Install NPM:
          npm install mongodb@2.2.5 --save
      Then go to Mongo Console:
          use users;
          db.users.insert({"username":"usertest","password":"usertest"});

## Test:
  node index.js usertest usertest
