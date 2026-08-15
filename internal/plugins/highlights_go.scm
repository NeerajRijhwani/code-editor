; Keywords
[
  "break"
  "default"
  "func"
  "interface"
  "select"
  "case"
  "defer"
  "go"
  "map"
  "struct"
  "chan"
  "else"
  "goto"
  "package"
  "switch"
  "const"
  "fallthrough"
  "if"
  "range"
  "type"
  "continue"
  "for"
  "import"
  "return"
  "var"
] @keyword

; Functions & Methods
(function_declaration name: (identifier) @function)
(method_declaration name: (field_identifier) @function)
(call_expression function: (identifier) @function)
(call_expression function: (selector_expression field: (field_identifier) @function))

; Types
(type_identifier) @type

; Literals
(interpreted_string_literal) @string
(raw_string_literal) @string
(rune_literal) @string
(int_literal) @number
(float_literal) @number

; Comments & Identifiers
(comment) @comment

