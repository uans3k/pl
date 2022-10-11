white  : (' '|'\t')
       ; {{ Hidden }}

line   : '\n'
       ; {{ Hidden,Row }}

string : ('a'|'b'|'c')('0'|'1'|'a'|'b'|'c')*
       ;

number : '0' | '1' ('0'|'1')*
       ;