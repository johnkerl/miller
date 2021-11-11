#!/usr/bin/env ruby

# ================================================================
# This is a no-longer-needed one-off script -- but it's an example of
# batch-editing for files in test/cases.
# ================================================================

require 'fileutils'

$parentdir = "test/cases/dsl-first-class-functions"

def make_case(name, mlr)
  casedir = $parentdir + '/' + name
  FileUtils.mkdir_p(casedir)
  puts casedir

  mlr = wrap(mlr)

  File.open(casedir + '/mlr', 'w') do |handle|
    handle.write(mlr)
  end
  File.open(casedir + '/cmd', 'w') do |handle|
    handle.puts('mlr -n put -f ${CASEDIR}/mlr')
  end
end

def wrap(mlr)
  return "end {\n  #{mlr}\n}\n"
end

done_entries = [
  [ "select-errors-01", 'print select("not array or map", func (k,v) { return v % 10 >= 5});' ],
  [ "select-errors-02", 'print select([])' ],
  [ "select-errors-03", 'print select({})' ],
  [ "select-errors-04", 'print select([], 2, 3, 4)' ],
  [ "select-errors-05", 'print select({}, 2, 3, 4)' ],
  [ "select-errors-06", 'print select([], "not a function")' ],
  [ "select-errors-07", 'print select({}, "not a function")' ],
  [ "select-errors-08", 'print select([], func () { return true});' ],
  [ "select-errors-09", 'print select({}, func () { return true});' ],
  [ "select-errors-10", 'print select([], func (a,b) { return true});' ],
  [ "select-errors-11", 'print select({}, func (a,b,c) { return true});' ],
  [ "select-errors-12", 'print select([1,2,3], func (e) { });' ],
  [ "select-errors-13", 'print select({"a":1,"b":2,"c":3}, func (k,v) { });' ],
  [ "select-errors-14", 'print select([1,2,3], func (e) { return "not a boolean"});' ],
  [ "select-errors-15", 'print select({"a":1,"b":2,"c":3}, func (k,v) { return "not a boolean"});' ],

  [ "apply-errors-01", 'print apply("not array or map", func (k,v) { return v % 10 >= 5});' ],
  [ "apply-errors-02", 'print apply([])' ],
  [ "apply-errors-03", 'print apply({})' ],
  [ "apply-errors-04", 'print apply([], 2, 3, 4)' ],
  [ "apply-errors-05", 'print apply({}, 2, 3, 4)' ],
  [ "apply-errors-06", 'print apply([], "not a function")' ],
  [ "apply-errors-07", 'print apply({}, "not a function")' ],
  [ "apply-errors-08", 'print apply([], func () { return true});' ],
  [ "apply-errors-09", 'print apply({}, func () { return true});' ],
  [ "apply-errors-10", 'print apply([], func (a,b) { return true});' ],
  [ "apply-errors-11", 'print apply({}, func (a,b,c) { return true});' ],
  [ "apply-errors-12", 'print apply([1,2,3], func (e) { });' ],
  [ "apply-errors-13", 'print apply({"a":1,"b":2,"c":3}, func (k,v) { });' ],
  [ "apply-errors-14", 'print apply({"a":1,"b":2,"c":3}, func (k,v) { return 999 });' ],
  [ "apply-errors-15", 'print apply({"a":1,"b":2,"c":3}, func (k,v) { return {} });' ],
  [ "apply-errors-16", 'print apply({"a":1,"b":2,"c":3}, func (k,v) { return {"x":7,"y":8} });' ],

  [ "reduce-errors-01", 'print reduce("not array or map", func (k,v) { return v % 10 >= 5});' ],
  [ "reduce-errors-02", 'print reduce([])' ],
  [ "reduce-errors-03", 'print reduce({})' ],
  [ "reduce-errors-04", 'print reduce([], 2, 3, 4)' ],
  [ "reduce-errors-05", 'print reduce({}, 2, 3, 4)' ],
  [ "reduce-errors-06", 'print reduce([], "not a function")' ],
  [ "reduce-errors-07", 'print reduce({}, "not a function")' ],
  [ "reduce-errors-08", 'print reduce([], func () { return true});' ],
  [ "reduce-errors-09", 'print reduce({}, func () { return true});' ],
  [ "reduce-errors-10", 'print reduce([], func (a,b,c) { return true});' ],
  [ "reduce-errors-11", 'print reduce({}, func (a,b,c,d,e) { return true});' ],
  [ "reduce-errors-12", 'print reduce([1,2,3], func (acc,e) { });' ],
  [ "reduce-errors-13", 'print reduce({"a":1,"b":2,"c":3}, func (acck,accv,ek,ev) { });' ],
  [ "reduce-errors-14", 'print reduce({"a":1,"b":2,"c":3}, func (acck,accv,ek,ev) { return 999 });' ],
  [ "reduce-errors-15", 'print reduce({"a":1,"b":2,"c":3}, func (acck,accv,ek,ev) { return {} });' ],
  [ "reduce-errors-16", 'print reduce({"a":1,"b":2,"c":3}, func (acck,accv,ek,ev) { return {"x":7,"y":8} });' ],

  [ "fold-errors-01", 'print fold("not array or map", func (k,v) { return v % 10 >= 5});' ],
  [ "fold-errors-02", 'print fold([])' ],
  [ "fold-errors-03", 'print fold({})' ],
  [ "fold-errors-04", 'print fold([], 2, 3, 4)' ],
  [ "fold-errors-05", 'print fold({}, 2, 3, 4)' ],
  [ "fold-errors-06", 'print fold([], "not a function", 9)' ],
  [ "fold-errors-07", 'print fold({}, "not a function", {"x":7})' ],
  [ "fold-errors-08", 'print fold([], func () { return true}, 9);' ],
  [ "fold-errors-09", 'print fold({}, func () { return true}, {"x":7});' ],
  [ "fold-errors-10", 'print fold([], func (a,b,c) { return true}, 9);' ],
  [ "fold-errors-11", 'print fold({}, func (a,b,c,d,e) { return true}, {"x":7});' ],
  [ "fold-errors-12", 'print fold([1,2,3], func (acc,e) { }, 9);' ],
  [ "fold-errors-13", 'print fold({"a":1,"b":2,"c":3}, func (acck,accv,ek,ev) { }, {"x":7});' ],
  [ "fold-errors-14", 'print fold({"a":1,"b":2,"c":3}, func (acck,accv,ek,ev) { return 999 }, {"x":7});' ],
  [ "fold-errors-15", 'print fold({"a":1,"b":2,"c":3}, func (acck,accv,ek,ev) { return {} }, {"x":7});' ],
  [ "fold-errors-16", 'print fold({"a":1,"b":2,"c":3}, func (acck,accv,ek,ev) { return {"x":7,"y":8} }, {"x":7});' ],
  [ "fold-errors-17", 'print fold({"a":1,"b":2,"c":3}, func (acck,accv,ek,ev) { return {"x":7} }, {"x":7, "y":8});' ],

  [ "sort-errors-01", 'print sort("not array or map", func (k,v) { return v % 10 >= 5});' ],
  [ "sort-errors-02", 'print sort([])' ],
  [ "sort-errors-03", 'print sort({})' ],
  [ "sort-errors-04", 'print sort([], 2, 3, 4)' ],
  [ "sort-errors-05", 'print sort({}, 2, 3, 4)' ],
  [ "sort-errors-06", 'print sort([], func () { return true});' ],
  [ "sort-errors-07", 'print sort({}, func () { return true});' ],
  [ "sort-errors-08", 'print sort([], func (a,b,c) { return true});' ],
  [ "sort-errors-09", 'print sort({}, func (a,b,c,d,e) { return true});' ],
  [ "sort-errors-10", 'print sort([1,2,3], func (a,b) { });' ],
  [ "sort-errors-11", 'print sort({"a":1,"b":2,"c":3}, func (ak,av,bk,bv) { });' ],
  [ "sort-errors-12", 'print sort({"a":1,"b":2,"c":3}, func (ak,av,bk,bv) { return {} });' ],

]

new_entries = [

]

new_entries.each do |entry|
  casedir = entry[0]
  mlr = entry[1]
  make_case(casedir, mlr)
end
