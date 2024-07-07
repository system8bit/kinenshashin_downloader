# 概要
pan.kinenshashin.netの画像を取得するスクリプト  
部分的なサムネイルを取得し、それらを結合して全体の画像を作成する

# 使い方
`setting.yaml`に設定を記述する必要がある  
`session_id`にはpan.kinenshashin.netにログインした際のセッションIDを指定する  
`photo_id`には写真のIDを指定する  
その他のパラメーターは写真サイズによって変更する必要があるかもしない

```yaml
session_id: ""
photo_id: ""
max_x: 144
max_y: 213
step_x: 10.0
step_y: 23.0
```
