$ ls -l ~/tmp/huge
-rw-r--r--  1 kerl  staff  614295000 Aug 25  2020 /Users/kerl/tmp/huge

$ wc -l ~/tmp/huge
 10000000 /Users/kerl/tmp/huge

$ justtime read-string ~/tmp/huge  > /dev/null
TIME IN SECONDS 8.707 -- read-string /Users/kerl/tmp/huge
$ justtime read-string ~/tmp/huge  > /dev/null
TIME IN SECONDS 8.540 -- read-string /Users/kerl/tmp/huge
$ justtime read-string ~/tmp/huge  > /dev/null
TIME IN SECONDS 8.549 -- read-string /Users/kerl/tmp/huge

$ justtime scanner ~/tmp/huge  > /dev/null
TIME IN SECONDS 8.774 -- scanner /Users/kerl/tmp/huge
$ justtime scanner ~/tmp/huge  > /dev/null
TIME IN SECONDS 8.873 -- scanner /Users/kerl/tmp/huge
$ justtime scanner ~/tmp/huge  > /dev/null
TIME IN SECONDS 8.777 -- scanner /Users/kerl/tmp/huge
