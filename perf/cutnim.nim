import strutils, tables

for line in stdin.lines:
  #var map: OrderedTable[string,string]
  var map = {"":""}.newOrderedTable
  #var map = initTable[string, string]
  #var map: OrderedTable[string, string]
  #var map: newOrderedTable[string, string](16)
  for word in line.split(","):
      var pair = word.split("=")
      #echo(pair[0])
      #echo(pair[1])
      #echo()
      #map[pair[0]] = pair[1]
      map.add(pair[0], pair[1])
