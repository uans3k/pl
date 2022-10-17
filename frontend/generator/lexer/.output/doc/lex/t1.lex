White  : (' '|'\t')
       ; {{ Hidden }}

Line   : '\n'
       ; {{ Hidden,Row }}

Left_b : '{'
       ;

Right_b: '}'
       ;

String : [a-zA-Z_][a-zA-Z0-9_]*
       ;

Number : '0' | [1-9][0-9]*
       ;
