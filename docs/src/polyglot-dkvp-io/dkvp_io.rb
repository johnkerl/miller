#!/usr/bin/env ruby

# ================================================================
# Example of DKVP I/O using Ruby.
#
# Key point: Use Miller for what it's good at; pass data into/out of tools in
# other languages to do what they're good at.
#
#   bash$ irb -I. -r dkvp_io.rb
#
#   # READ
#   irb(main):001:0> map = dkvpline2map('x=1,y=2', '=', ',')
#   => {"x"=>"1", "y"=>"2"}
#
#   # MODIFY
#   irb(main):001:0> map['z'] = map['x'] + map['y']
#   => 3
#
#   # WRITE
#   irb(main):002:0> line = map2dkvpline(map, '=', ',')
#   => "x=1,y=2,z=3"
#
# ================================================================

# ----------------------------------------------------------------
# ips and ifs (input pair separator and input field separator) are nominally '=' and ','.
def dkvpline2map(line, ips, ifs)
  map = {}
  line.split(ifs).each do |pair|
    (k, v) = pair.split(ips, 2)

    # Type inference:
    begin
      v = Integer(v)
    rescue ArgumentError
      begin
        v = Float(v)
      rescue ArgumentError
        # Leave as string
      end
    end

    map[k] = v
  end
  map
end

# ----------------------------------------------------------------
# ops and ofs (output pair separator and output field separator) are nominally '=' and ','.
def map2dkvpline(map, ops, ofs)
  map.collect{|k,v| k.to_s + ops + v.to_s}.join(ofs)
end
