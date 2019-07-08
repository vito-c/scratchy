package scratch

import (
	. "github.com/alecthomas/chroma" // nolint
	// "github.com/alecthomas/chroma/lexers/internal"
)

// Java lexer.
var JqLex = MustNewLexer(
	&Config{
		Name:      "jq",
		Aliases:   []string{"jq"},
		Filenames: []string{"*.jq"},
		MimeTypes: []string{"text/x-jq"},
		DotAll:    true,
	},

	//
	// " jqConditional
	// syntax keyword jqConditional if then elif else end
	//
	// " jqSpecials
	// syntax keyword jqType type
	// syntax match jqType /[\|;]/ " not really a type I did this for coloring reasons though :help group-name
	// syntax region jqParentheses start=+(+ end=+)+ fold transparent
	//
	// " TODO: $__loc__ is going to be a pain
	//
	//
	Rules{
		"root": {
			{`[^\S\n]+`, Text, nil},
			{`#.*?\n`, CommentSingle, nil},
			// jq Math Functions
			{`(acos|acosh|asin|asinh|atan|atanh|cbrt|ceil|cos|cosh|erf|erfc|exp|exp10|exp2|expm1|fabs|floor|gamma|j0|j1|lgamma|lgamma_r|log|log10|log1p|log2|logb|nearbyint|pow10|rint|round|significand|sin|sinh|sqrt|tan|tanh|tgamma|trunc|y0|y1|atan2|copysign|drem|fdim|fmax|fmin|fmod|frexp|hypot|jn|ldexp|modf|nextafter|nexttoward|pow|remainder|scalb|scalbln|yn|fma)\b`, Keyword, nil},
			// Conditioanls
			{`(if|then|elif|else|end)\b`, Keyword, nil},
			// " jq Functions
			{`(and|or|not|empty|try|catch|reduce|as|label|break|foreach|import|include|module|modulemeta|env|nth|has|in|while|error|stderr|debug)\b`, Keyword, nil},
			{`(add|all|any|arrays|ascii_downcase|floor|ascii_upcase|booleans|bsearch|builtins|capture|combinations|contains|del|delpaths|endswith|explode|finites|first|flatten|format|from_entries|fromdate|fromdateiso8601|fromjson|fromstream|get_jq_origin|get_prog_origin|get_search_list|getpath|gmtime|group_by|gsub|halt|halt_error|implode|index|indices|infinite|input|input_filename|input_line_number|inputs|inside|isempty|isfinite|isinfinite|isnan|isnormal|iterables|join|keys|keys_unsorted|last|leaf_paths|length|limit|localtime|ltrimstr|map|map_values)\b`, Keyword, nil},
			{`(match|max|max_by|min|min_by|mktime|nan|normals|now|nulls|numbers|objects|path|paths|range|recurse|recurse_down|repeat|reverse|rindex|rtrimstr|scalars|scalars_or_empty|scan|select|setpath|sort|sort_by|split|splits|with_entries|startswith|strflocaltime|strftime|strings|strptime|sub|test|to_entries|todate|todateiso8601|tojson|__loc__|tonumber|tostream|tostring|transpose|truncate_stream|unique|unique_by|until|utf8bytelength|values|walk)\b`, Keyword, nil},
			// " jq SQL-style Operators
			{`(syntax|keyword|jqFunction|INDEX|JOIN|IN)`, Keyword, nil},
			// syntax keyword jqCondtions true false null
			{`(true|false|null)\b`, KeywordConstant, nil},
			// No idea if this is right  I just copied it from java
			{`((?:(?:[^\W\d]|\$)[\w.\[\]$<>]*\s+)+?)((?:[^\W\d]|\$)[\w$]*)(\s*)(\()`, ByGroups(UsingSelf("root"), NameFunction, Text, Operator), nil},
			// {`@[^\W\d][\w.]*`, NameDecorator, nil},
			// {`(abstract|const|enum|extends|final|implements|native|private|protected|public|static|strictfp|super|synchronized|throws|transient|volatile)\b`, KeywordDeclaration, nil},
			// {`(boolean|byte|char|double|float|int|long|short|void)\b`, KeywordType, nil},
			// {`(package)(\s+)`, ByGroups(KeywordNamespace, Text), Push("import")},
			// {`(class|interface)(\s+)`, ByGroups(KeywordDeclaration, Text), Push("class")},
			// {`(import(?:\s+static)?)(\s+)`, ByGroups(KeywordNamespace, Text), Push("import")},
			// {`"(\\\\|\\"|[^"])*"`, LiteralString, nil},
			// {`'\\.'|'[^\\]'|'\\u[0-9a-fA-F]{4}'`, LiteralStringChar, nil},
			// {`(\.)((?:[^\W\d]|\$)[\w$]*)`, ByGroups(Operator, NameAttribute), nil},
			// {`^\s*([^\W\d]|\$)[\w$]*:`, NameLabel, nil},
			// {`([^\W\d]|\$)[\w$]*`, Name, nil},
			{`([0-9][0-9_]*\.([0-9][0-9_]*)?|\.[0-9][0-9_]*)([eE][+\-]?[0-9][0-9_]*)?[fFdD]?|[0-9][eE][+\-]?[0-9][0-9_]*[fFdD]?|[0-9]([eE][+\-]?[0-9][0-9_]*)?[fFdD]|0[xX]([0-9a-fA-F][0-9a-fA-F_]*\.?|([0-9a-fA-F][0-9a-fA-F_]*)?\.[0-9a-fA-F][0-9a-fA-F_]*)[pP][+\-]?[0-9][0-9_]*[fFdD]?`, LiteralNumberFloat, nil},
			{`0[xX][0-9a-fA-F][0-9a-fA-F_]*[lL]?`, LiteralNumberHex, nil},
			{`0[bB][01][01_]*[lL]?`, LiteralNumberBin, nil},
			{`0[0-7_]+[lL]?`, LiteralNumberOct, nil},
			{`0|[1-9][0-9_]*[lL]?`, LiteralNumberInteger, nil},
			{`[~^*!%&\[\](){}<>|+=:;,./?-]`, Operator, nil},
			{`\n`, Text, nil},
		},
		"class": {
			{`([^\W\d]|\$)[\w$]*`, NameClass, Pop(1)},
		},
		"import": {
			{`[\w.]+\*?`, NameNamespace, Pop(1)},
		},
	},
)
