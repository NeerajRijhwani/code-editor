; Keywords
[
  "and"
  "as"
  "assert"
  "async"
  "await"
  "break"
  "class"
  "continue"
  "def"
  "del"
  "elif"
  "else"
  "except"
  "exec"
  "finally"
  "for"
  "from"
  "global"
  "if"
  "import"
  "in"
  "is"
  "lambda"
  "nonlocal"
  "not"
  "or"
  "pass"
  "print"
  "raise"
  "return"
  "try"
  "while"
  "with"
  "yield"
] @keyword

; Functions & Classes
(function_definition name: (identifier) @function)
(class_definition name: (identifier) @type)
(decorator "@" name: (identifier) @function)
(call_expression function: (identifier) @function)
(call_expression function: (attribute attribute: (identifier) @function))

; Literals
(string) @string
(integer) @number
(float) @number

; Comments & Constants
(comment) @comment
(none) @constant
[ "True" "False" ] @boolean
