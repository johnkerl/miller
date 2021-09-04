#!/usr/bin/env ruby

# ================================================================
# Markdown autogen for built-in functions: table, or full details page.
# Invoked by files in *.md.in.
# ================================================================

# ----------------------------------------------------------------
# Displayname is for HTML rendered in the document body. Linkname is for
# #-style links which don't do so well with special characters.

LOOKUP = {
  # name     display_name link_name
'!'    => [ '\!',      'exclamation-point'        ],
'!='   => [ '\!=',     'exclamation-point-equals' ],
'!=~'  => [ '!=~',     'regnotmatch'              ],
'%'    => [ '%',       'percent'                  ],
'&&'   => [ '&&',      'logical-and'              ],
'&'    => [ '&',       'bitwise-and'              ],
'*'    => [ '\*',      'times'                    ],
'**'   => [ '\**',     'exponentiation'           ],
'+'    => [ '\+',      'plus'                     ],
'-'    => [ '\-',      'minus'                    ],
'.'    => [ '\.',      'dot'                      ],
'.*'   => [ '\.\*',    'dot-times'                ],
'.+'   => [ '\.\+',    'dot-plus'                 ],
'.-'   => [ '\.\-',    'dot-minus'                ],
'./'   => [ '\./',     'dot-slash'                ],
'/'    => [ '/',       'slash'                    ],
'//'   => [ '//',      'slash-slash'              ],
':'    => [ '\:',      'colon'                    ],
'<'    => [ '<',       'less-than'                ],
'<<'   => [ '<<',      'lsh'                      ],
'<='   => [ '<=',      'less-than-or-equals'      ],
'='    => [ '=',       'equals'                   ],
'=='   => [ '==',      'double-equals'            ],
'=~'   => [ '=~',      'regmatch'                 ],
'>'    => [ '>',       'greater-than'             ],
'>='   => [ '>=',      'greater-than-or-equals'   ],
'>>'   => [ '>>',      'srsh'                     ],
'>>>'  => [ '>>>',     'ursh'                     ],
'>>>=' => [ '>>>=',    'ursheq'                   ],
'?'    => [ '?',       'question-mark'            ],
'?.:'  => [ '?:',      'question-mark-colon'      ],
'?:'   => [ '?:',      'question-mark-colon'      ],
'??'   => [ '??',      'absent-coalesce'          ],
'???'  => [ '???',     'absent-empty-coalesce'    ],
'^'    => [ '^',       'bitwise-xor'              ],
'^^'   => [ '^^',      'logical-xor'              ],
'|'    => [ '\|',      'bitwise-or'               ],
'||'   => [ '\|\|',    'logical-or'               ],
'~'    => [ '~',       'bitwise-not'              ],
}

def name_to_display_name_and_link_name(name)
  entry = LOOKUP[name]
  if entry.nil?
    return [name, name]
  else
    display_name, link_name = entry
    return [display_name, link_name]
  end
end

# ----------------------------------------------------------------
def make_func_details
  function_classes = `mlr help list-function-classes`.split

  puts
  puts "## Summary"
  puts
  for function_class in function_classes
    class_link_name = "#"+"#{function_class}-functions"
    class_display_name = "#{function_class.capitalize} functions"
    #puts "* [#{display_name}](#{link_name})"

    functions_in_class = `mlr help list-functions-in-class #{function_class}`.split
    foo = []
    for function_name in functions_in_class
      display_name, link_name = name_to_display_name_and_link_name(function_name)
      foo.append " [#{display_name}]("+"#"+"#{link_name})"
    end
    #puts "    * #{foo.join(',')}"
    puts "* [**#{class_display_name}**](#{class_link_name}): #{foo.join(', ')}."

  end
  puts

  for function_class in function_classes
    display_name = "#{function_class.capitalize} functions"
    puts "## #{display_name}"
    puts

    functions_in_class = `mlr help list-functions-in-class #{function_class}`.split

    for function_name in functions_in_class
      display_name, link_name = name_to_display_name_and_link_name(function_name)

      puts
      if display_name != link_name
        puts "<a id=#{link_name} />"
      end
      puts "### #{display_name}"

      puts '<pre class="pre-non-highlight-non-pair">'
      system("mlr help function '#{function_name}'")
      puts '</pre>'
      puts
    end
  end
end

# ----------------------------------------------------------------
make_func_details
