#!/usr/bin/env ruby 

alphabet = 'abcdefghijklmnopqrstuvwxyz0123456789~!#$%^&*()-_=+{}[];:"<>,./?'
nlines=50

alphabet = alphabet.split('')
nlines.times do
  length = rand(40)
  if rand < 0.1
    length = 0
  end
  output = (1..length).to_a.collect{alphabet.sample}.join('')
  puts output
end
