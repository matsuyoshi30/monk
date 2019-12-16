# Writing An Interpreter in GO

### Chapter1: LEXING

##### Lexical Analysis

Source Code -<lexical analysis>-> Tokens -<parse>-> Abstract Syntax Tree  

##### Defining our tokens

##### The Lexer

string and rune  
string ... byte array  

[reference](http://text.baldanders.info/golang/string-and-rune/)

doesn't support float, hex notation, octal notation  
only integers  

##### Extending our token set and lexer

The lexer’s job is not to tell us whether code makes sense, works or contains errors.  

peekChar() is for reading two-characters operators like "==" or "!="  
peek ... スタック (LIFO) と呼ばれるコンピュータのデータ構造において、トップにあるデータを取り出さずに読み取る操作。  

##### Starting REPL

Read, Eval, Print and Loop  


### Chapter2: PARSING

##### Parsers

Source Code -<lexical analysis>-> Tokens -<parse>-> AST  

##### Why not a parser generator?

Parser generator is like yacc, bison or ANTLR.  

FORMAL DESCRIPTION LANGUAGE -<parser generator>-> PARSER  

The majority of them use a context-free grammer (CFG) as their input.  
The most common notational formats of CFGs are the Backus-Naur Form (BNF) or the Extended Backus-Naur Form (EBNF).  

The best way to do that is by getting our hands dirty and writing a parser ourselves.  

##### Writing a parser for monk programming language

top-down parsing or bottom-up parsing  

We use "recuisive descent parsing" which is one of top-down parsing.  

##### Parser's first step: parsing let statement

take a look at monk source code and understand how it is structured.  

```
let x = 5;
let y = 10;

let add = fn(a, b) {
    a + b;
}
```

it goes,

```
let <identifier> = <expression>;
```

The distinction: expression produces value, but statement doesn't produce value.  


##### Implementing the Pratt Parser

2018-12-26  

与えられた文字列を解析してトークン郡を生成し、それをもとに AST を構築する準備までできた  
注）let と return だけ  

Pratt Parser というのをもとにパーサを書いていく  
prefix (-x とか) と infix (x + y) のどちらについてもパースできるようにする  


##### Identifiers

どんな Identifier であっても、 produce a value する  

parseStatement() ---> parseExpressionStatement() ---> parseExpression()  

##### Integer

2018-12-27  

let と return のAST化やったので、次は数値（int64）のパーサを作成  
strconvつかって、文字列として入力されたもの（"5"）を数値（5）に変換  

##### Prefix Operator

2018-12-28  

prefix 前に ! か - がついているもの  

```usage
-5;
!foobar;
5 + -10;
```

<prefix operator><expression>;  

先にテストを書く  
"-5;" がインプットのとき、 prefix operator が - 、expression が 5 とパースされるかのテストを書く  

prefixExpression のパーサを書く(parsePrefixExpression)  
parsePrefixExpression は parseExpression から呼ばれる  
parsePrefixExpression は token, operator が curToken の prefixExpression を構築して、トークンを一個進めて数値を取得する  

##### Infix Operator

```usage
5 + 5;
5 - 5;
5 * 5;
5 / 5;
5 > 5;
5 < 5;
5 == 5;
5 != 5;
```

みたいなもの  

<expression><infix operator><expression>;  

ここも先にテスト  
上に挙げたケースで、left value, operator, right value が想定通りにパースされるかどうかのテスト  

パースするときは precedence（優先順位）の定数で構築する ast.Node の深さを判断する  
ここまでは LOWEST しか用いていないが、ここからは precedences table を作って対応する  
なぜかというと、ここで +, - と *, / が絡む AST の構築になるから  

### How Pratt Parsing works

2019-01-02  

新年明けましておめでとうございます  

Vaughan Pratt の “Top Down Operator Precedence という論文に基づいてパーサを書いている  

1 + 2 + 3; がどうやって我々の求めるASTの形である ((1 + 2) + 3) になるかを見る  

parser_tracing.go を作ってトレース用のメソッドを定義して各パーサにdeferする  

### Extending our parser

Boolean Literal, Grouped Expression を実装  

##### if expression

if-else 構文は、以下の通り考える  

if (<condition>) <consequence> else <alternative>  

ぼくたちのレシピは、  
AST ノードを定義する -> テストを書く -> テストが通るようなパーサを書く  
なので、 if-else のAST ノードを定義する  

2019-01-03  

if のあとには、 ( がくるはずで、その中には条件(condition)があるはずで、そのあとには ) がくるはずで、  
そのあとには { がきて結果(consequence)がくるはず  
という parseIfExpression を実装  
consequence は {} のなかのブロック文(BlockStatement)のパース結果が入るので、parseBlockStatement を別で実装  

parseBlockStatementの実装は以下の流れ  
BlockStatement は Statement の集合体。これまで同様、BlockStatement のASTノードを生成して、  
} がでてくるか、もしくは EOF にあたるまで、BlockStatement に stmt を追加していって、最後に AST を返す  

##### Function Literals

こんなやつのパーサを書く  

```
fn(x, y) {
  return x + y;
}
```

abstract structure はこんなかんじ  

fn <parameters> <block statements>  

block statements のパース方法は前で見ているので無問題。parameters については、identifierが  
(<parameter one>, <parameter two>, <parameter three>, ...) と考える   

Function Literal のASTノードは、、、Expression!  
Expression が許される場所であれば、どこでも fn できる  

いつものとおりASTノードを定義  

テストメソッド作成  
function の parameter と body をテスト  

で、テストを通るようにパーサをかく  
FuntionLiteral 自体のパーサの中で、パラメータ群をパースする parseFunctionParameters も実装する  
このとき、パラメータの数が0 のものは parseFunctionParameters のエッジケースになるので、別途このテストメソッドも作る  

##### Call Expression

function 呼び出しの式のパース  
<expression>(<comma separated expression>)  

```add(2, 3)```は、  
add が identifier, 2と3がexpression になる  
expressionなので、```add(2+2, 3*3*3)``` もあり  

Call Expression のASTノードを定義する  

テストメソッドを書く  
prefixparseのエラーがでてこない -> call expression の構成に着目する  

<expression>(<comma separated expression>)  

一番はじめの<expression> -> identifier としてトークナイズされる  
arguments をともなうカッコの中身は、infixとして考える  
5 * 5 が infix なのとおなじように、 ( 5 ) も infix として考えていく  

##### REMOVE TODOs

LetStatement と ReturnStatement のパーサ実装時に飛ばしていたTODOを消化する  

セミコロンがくるまでトークンを進めていたが、  
LetStatement は assign される value を正しくパースできるか、  
ReturnStatement は returnvalue を正しくパースできるか、  
のテストを書いて、TODOを消しておしまい  

### not REPL, but RPPL (read-parse-print-loop)

repl.go に parser をたして完了  


# Chapter3: Evaluation

### Giving meaning to symbols

トークナイズ、パースときて、あとは評価（eval）  
どの形式ならばどう評価されるかをここの monk の文法として実装していく  

### Strategies of evaluation

Eval は実装の方法が多岐にわたる  
AST をトラバースして評価していくやり方を tree-walking interpreters という  

JITの説明  

tree-walking は遅いけどわかりやすい  
JIT は早いけど複雑  

### A Tree-walking intepreter

簡単だしとっつきやすいから、 tree-walking interpreter 方式を採用する  

### Representing objects

2019-01-04  

eval 関数の戻り値について考える  
ASTノードを eval したあとの戻り値は何が適切？  

あとは変数の評価  
let a = 5; としたあとに、 a + a; を評価するときは、a にバインドされている 5 で評価する必要がある  

パフォーマンスとかも考えよう  
いろんなインタープリタのソース見て勉強するのもいいね  

##### foundation of our object system

token.go みたいに object を定義していく  

Token と TokenType のような、 Object と ObjectType の定義  
Object をインタフェースで定義しているのは、 integer と boolean を同じ構造体定義で定義していきたいから  

ここでは、  
integer, boolean, null を定義する  

### Evaluating Expressions

Eval は、 ASTノードを受け取って、オブジェクトを返す関数  

とりあえず、 type in 5, get back 5 をつくる  

##### integer literal

input: *ast.IntegerLiteral  
output: *object.Integer  
となる  

eval 関数は当然再帰的（ast.Program を評価して、Statement を評価して、Expression を評価して、、、）  
ASTノードの type に応じて eval 関数の戻り値を実装する  

##### Completing the REPL

parse したものを eval するように repl.go を修正  
eval の戻り値が nil 出ない場合に戻り値を Inspect() で出力する  

##### Boolean Literal

Boolean はふたつしか取りうる値がないのに、いちいちオブジェクトを生成しておくのは無駄なのでは？  
eval 側でそれを吸収して無駄を拝する  

具体的には、  
先に true false のオブジェクトを定義しておいて、都度そのオブジェクトを呼び出す（使う）イメージ  

##### NULL

Boolean と同じように NULL を定義しておしまい  

##### Prefix Expression

! と - のみのサポートなので、それらを考える  

##### Infix Expression

Infix Expression は、結果が数値を返すものと、結果が真偽を返すもののふたつに分かれる  

まずは、結果が数値を返すものを考える  
テストはすでに書いた TestEvalIntegerExpression で実施する  

evalInfixExpression で Infix を評価し、  
左右が数値の場合はオペレータに応じて左右の数値による計算結果のオブジェクトを生成する  

次に真偽を返すものを考える  
でもこれは簡単か。今は数値と真偽とNULLしか考えていないので、先程考えた評価のメソッドについて、  
オペレータの条件分岐を増やす  

真偽については、 == != を使用することができるので、対応する

### Conditionals

2019-01-05  

if-else 文は、まず if 文だけ評価する。それでマッチしなかったら、初めて else 文を評価する  
けど if 文の評価の結果はどう分岐させるのがただしいのか？  

ここでは、if 文の条件が null でも false でもないときに、マッチしていると判断する  

### Return Statement

return 文は top-level  

```
5 * 5 * 5; 
return 10; 
9 * 9 * 9;
```

の場合は、 10 が返って来なければならないし、最後の statement は絶対に評価されない  
return 文が評価された時点でおしまい  

ネスト構成のときに正しく動作するよう、 blockstatement と program の評価時の関数を分けて実装  

2019-01-06  

ゼロ除算がアウトだったので、エラーメッセージを出力するよう修正  

### Error Handling

エラーとする事象の分類  

return statement の実装と同じように、エラーハンドリングの実装を進める  
return statement と同様、エラーが発生したら statement の評価は止める  

### Binding & the environment

let statement による変数のバインドを実装する  

テストケースは、let a = 5; で宣言した変数 a を評価した場合に、正しくバインドした値になるか  
と、バインドしていない変数を評価したときに正しくエラーとなるか  

変数のバインドと、変数の評価は以下の流れで考える  
変数のバインドは、key - value の map に保存する  
変数の評価は、map の key に対して存在確認をする  

変数のバインド先は、 Environment という新しいオブジェクトを定義する  

### Function & function calls

関数の定義と関数の呼び出しについて  

やることは２つ  
- 関数の定義について、object としてどう保持するか
- 関数の呼び出しについて、eval でどう評価するか
を考える  

関数の定義について考える  
ASTノードとしての関数（FunctionLiteral）は、引数Parameter と処理BlockStatement を持つ  

2019-01-07  

関数呼び出し考える  

ASTノードの CallExpression を評価するときの処理を考える  
まず関数名の評価（これは変数の評価と同じ）  
次に引数の評価（引数は最低一つの Expression の評価となる）

ここで、関数定義の際に引数に使用した変数の適用範囲（グローバル変数との競合解決）を考える必要が出てくる  

```
let i = 5;
let printNum = fn(i) {
  puts(i); 
};
printNum(10);
puts(i);
```

上記は、ひとつめは 10 、ふたつめは 5 と出力されなければならないが、  
一番始めに宣言した変数名と、関数の宣言内で使用している変数名が同じであることを内部的にうまくデザインする必要がある  
ここでは、変数のマップが Environment に依存していることに着目し、 Environment の拡張というアプローチを採用する  

上記で問題となるのは、関数の宣言内で使用する変数をどうマップに保持するかということ  
（すでにその前で i という変数はマップにセットされているから）  

ここでは、 EncloseEnvironment というものを考える（外に別の環境が存在している構造）  


2019-01-09  

# Chapter4: Extending the interpreter

2019-01-10  

### Data types & functions

add new data type is  
add new token type, modify the lexer, extend the parser and add support to evaluater and object system!  

add built-in function  

### new data type, string

```
   "<sequence of characters>"
```

トークンタイプ「String」を定義  
String トークンのレキサーを定義  
ASTノードに StringLiteral ノードを定義  
パーサを拡張  
eval 関数を追加  

2019-01-11  

### Built-in Function

2019-01-13  

あるひとつの string を引数にとり、引数の文字数を返す len を builtin 関数として定義する  

builtin のオブジェクトを定義  
builtin で定義した関数を格納する builtin の map を定義  
evalIdentifier で、引数が builtin だった場合の処理を追加  
evaluator の applyFuntion で、引数が builtin だった場合の処理を追加  

### Array

2019-01-14  

in monk, array is an ordered list of elements of possibly different types.  

要素の順序立てられたリスト（異なる型の要素も許容する）  

##### supporting arrays in our lexer

token に新しく [, ] のトークンを定義  
lexer に新しいトークンの解析処理を追加  

##### parsing array literals

要素のパース方法  
→関数の引数をパースしたときの方法でできる  

ASTノードに arrayliteral を追加  
パーサのテストを書いてパーサを書く  

複数の Expression のリストをパースする関数を別途追加（引数は、見つけたら終了する tokentype）  
→関数の引数のパースもこっちに変える  

##### parsing index operator expressions

配列名[インデックス] での配列の要素へのアクセスもパースできるようにする  

<expression>[<expression>]  
が基本的な構造  

上記の基本的な構造を満たすASTノードを定義  

<expression>[<expression>] が配列の要素へのアクセスの文法となったので、  
オペレータの優先順位も再度考える


例  
```
a + [1, 2, 3, 4][b * c] * d
```

parser.go について  

const ブロックの最後に INDEX 追加して、 precedences の最後に token.LBRACKET: INDEX とする  
-> why?  

##### evaluating array literals

array の object を定義  
object.Array の評価を定義  

##### evaluating index operator expressions

範囲外へのアクセスはNULL応答とする  

##### adding built-in functions for arrays

first ... 配列の最初の要素を取り出す  
last  ... 配列の最後の要素を取り出す  

どちらも、  
引数が配列のオブジェクトであること、配列の要素数が0より大きいこと  

rest ... scheme でいう cdr. 配列の最初の要素以外の配列を返す  

push ... 引数に配列とオブジェクト１つを取り、オブジェクトを配列の末尾に追加する  

##### test-driving arrays

array の rest とか push とか使えば、以下の通り map 作れる  

```
let map = fn(arr, f) {
  let iter = fn(arr, accumulated) {
    if (len(arr) == 0) {
      accumulated
    } else {
      iter(rest(arr), push(accumulated, f(first(arr))));
    }
  };

  iter(arr, []);
};
```

reduce とか  

```
let reduce = fn(arr, initial, f) {
  let iter = fn(arr, result) {
    if (len(arr) == 0) {
      result
    } else {
      iter(rest(arr), f(result, first(arr)));
    }
  };

  iter(arr, initial);
};
```

2019-01-15  

### hash

キーと値の組み合わせのマップ  

monk では、name["string"] や name[boolean], name[integer] を許可することとする  
なので、通常みたいな name["string"] の他に name[5 > 1], name[100 - 1] とかもいける  
もちろん、与えたインデックスの expression を処理して、その結果のキーがあればの話！  

##### lexing hash literals

ハッシュの定義のうちまだ追加していないトークンはコロン : だけ  
-> コロンを追加（lexer のテストも追加）  

##### parsing hash literals

{<expression>: <expression>, <expression>: <expression> ... }  

という文法  

expression はこの段階では何でも許可  
エラー出すのは評価 eval のところでやる  

以下の例参考  

let key = "name";  
let hash = {key: "Monkey"};  

これは文法としては正しいのに、expression として key は string, boolean, integer どれでもない  
（上の key を評価していないと、 hash の key が "name" であると解決できない）  

まずは ast ノードの定義  

で、パーサ  
パーサのテストは、期待値のマップを作っておいて、  
まずはパースした結果のキーを取り出すところでチェック  
つぎにキー使って期待値のマップにアクセスして、値の突き合わせ  
というかたち  

コーナーケース（ハッシュが空）は、べつのテストメソッドつくる  

あとは、ハッシュのキーや値が単項目ではなくて、何か expression 5-1とか5>1とかになっているパターン  

##### hashing object

```
type Hash struct {
  Pairs map[object]object
}
```

としてはいけない理由  

```
let hash = {"name": "Monkey"};
hash["name"]
```

が正しい評価値を返さない  
なぜかというと、Go のポインタの問題  

```
name1 := &object.String{Value: "name"}
monkey := &object.String{Value: "Monkey"}

pairs := map[object.Object]object.Object{}
pairs[name1] = monkey

fmt.Printf("pairs[name1]=%+v\n", pairs[name1])
// => pairs[name1]=&{Value:Monkey}

name2 := &object.String{Value: "name"}
fmt.Printf("pairs[name2]=%+v\n", pairs[name2])
// => pairs[name2]=<nil>

fmt.Printf("(name1 == name2)=%t\n", name1 == name2)
// => (name1 == name2)=false
```

上記の例から分かる通り、 name1 とあとで定義した name2 は、value は同じ "name" で定義したはず  
なのにアロケートされるメモリが違うので、pairs の引数に name2 を渡すと nil が返ってくる  

ここは、、、  recheck  

2019-01-16  

##### Evaluating hash literals

評価  

ここは地味な作業  

ASTノードのHashLiteralを受け取ったら、キー、値という順番でEvalして、  
その結果を作成しておいたPairのマップにそれぞれ格納してオブジェクトを返す  
という流れ  

##### Evaluating index expressions with hashes

hashname[keyname] が正しく評価されるようにする  

