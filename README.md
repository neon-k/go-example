# go-example

業務でGoを使って開発をすることがあったので、実際に手を動かしながら勉強した時のリポジトリです。

## Purpose
- gin + gormを使ってrest APIを構成する
- DDDアーキテクチャーの使い方確認
- AWSのs3に画像（base64）のアップロード方法
- JWTを使ってトークン認証実装

## Library

使用しているライブラリは下記のとおりです。

| Library                                             | 説明                                         |
| --------------------------------------------------- | -------------------------------------------- |
| [gin](https://github.com/gin-gonic/gin)             | goのフレームワーク                         |
| [gorm](https://gorm.io/index.html)                  | ORMマッパー                          |
| [goose](https://github.com/pressly/goose)           | マイグレーション                                   |
| [mysql](https://www.mysql.com/jp/)                  | DB                                 |
| [swagger](https://swagger.io/)                      | API仕様書                                   |
| [gin-jwt](https://github.com/appleboy/gin-jwt)      | JWT認証      |

## Impressions

- gormがActive Recordと使い方が似ていたので使いやすかった。ただ、更新する処理のsaveの使い方を気おつけないと事故る危険性が大いにあると感じた。
- 構造体をそのままswaggerのparamsに渡せるのはとてもよかった。
- DDDの構造体で作ることによって『細かく作る』が実現できるのと、*モデル* の考え方を上手くプロジェクトチームと共有できれば企画チームと開発チームの会話がスムーズになると感じた。
- gormのマイグレーション機能はあまり実践的ではないのでgooseを使ったが、TypeORMと比較した時にモデルの構造体と別々で管理しないといけないのが少し辛いと感じた。

## Reference

- [https://gorm.io/ja_JP/docs/conventions.html](https://gorm.io/ja_JP/docs/conventions.html)
- [https://qiita.com/Syoitu/items/8e7e3215fb7ac9dabc3a](https://qiita.com/Syoitu/items/8e7e3215fb7ac9dabc3a)
- [https://github.com/appleboy/gin-jwt](https://github.com/appleboy/gin-jwt)
- [https://qiita.com/Bmouthf/items/d2d6ef55e131cb5a9e82](https://qiita.com/Bmouthf/items/d2d6ef55e131cb5a9e82)
- [http://golang-jp.org/](http://golang-jp.org/)
