# Visualize Channel
Goのソースコードにインジェクションすることでチャネル情報を可視化するCLIツール

# TODOs
- まずはmakeのインジェクションから
  - makeを見つけ、変数名とchannelの型、長さが知りたい
  - 正規表現でもできそう
    - Goの文法を網羅できるかわからない
  - 静的解析でした方が確実？
    - https://pkg.go.dev/go/ast
- send, recvのインジェクション
  - いったんmain関数内のみを考える
- closeのインジェクション
- syncパッケージを使うときに正しく動作するか
