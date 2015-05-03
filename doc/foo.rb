output_handle=$stdout
pages_map = {'bar.html' => 'barbar'}
line = "PUT_LINK_FOR_PAGE(bar.html)HERE"
      if false
        puts 'no'
      elsif line =~ /PUT_LINK_FOR_PAGE\(([^)]+)\)HERE/
        other_page_name = $1
        other_page_title = pages_map[other_page_name]
        if other_page_title.nil?
          raise "Couldn't find page title for \"#{other_page_name}\" in $0."
        end
        href = "<a href=\"#{other_page_name}\">#{other_page_title}</a>"
        line.sub!(/PUT_LINK_FOR_PAGE\([^)]+\)HERE/, href)
        output_handle.puts(line)
      else
        puts 'no match'
      end
