#/usr/bin/ruby

rcorn   = -2.0
icorn   = -2.0
side    =  4.0
iheight =  500
iwidth  = 1000
maxits  =  100

for ii in 0..iheight
  ci = icorn + ((ii*1.0)/iheight) * side
	for ir in 0..iwidth
    cr = rcorn + ((ir*1.0)/iwidth) * side

		zr = 0.0
		zi = 0.0

		# z := z^2 + c
		iti = 0
		escaped = false
    for iti in 0..maxits
			mag = zr*zr + zi+zi
			if mag > 4.0
					escaped = true
					break
      end
			zt = zr*zr - zi*zi + cr
			zi = 2.0*zr*zi + ci
			zr = zt
    end
		if escaped
			print("o")
		else
			print(".")
    end

  end
	puts
end
