# テスト規約

## 基本方針

- 原則、テーブルドリブンテストを使用
- テスト関数名：`Test` + 機能名
- `t.Parallel()`を積極使用（有害な副作用がない限り）
- `testify/assert`と`testify/require`を使用

## テストコメント規約

- **テストコードには意図を明確化するためのコメント記述を許容する**
- ただし、コメントは英語で記述すること
- テストの意図、複雑なアサーション、特殊な条件設定の説明に使用
- 例: `// Test error handling when user is not found`

## テストケース命名

| 種類 | 条件 | 命名規則 |
|------|------|----------|
| 正常系 | `err=nil`を返す | `success:` で開始 |
| 異常系 | `err!=nil`を返す | `failure:` で開始 |

## テストケース固有のアサーション規約

- **テストケース固有のアサーションが必要な場合、テストケースの匿名構造体に`assert func(t *testing.T)`メンバを追加**
- 共通のアサーションは`t.Run`内で直接実行し、固有のアサーションのみ`assert`関数内で実行
- `assert`関数が`nil`の場合は、共通のアサーションのみ実行される

### 使用例

```go
tests := []struct {
    name   string
    input  string
    want   string
    assert func(t *testing.T) // テストケース固有のアサーション
}{
    {
        name:  "success: standard case",
        input: "test",
        want:  "expected",
        // assert: nil (共通アサーションのみ)
    },
    {
        name:  "success: special validation case",
        input: "special",
        want:  "result",
        assert: func(t *testing.T) {
            // テストケース固有の検証ロジック
            assert.Contains(t, result, "special")
        },
    },
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        result := targetFunction(tt.input)

        // 共通のアサーション
        assert.Equal(t, tt.want, result)

        // テストケース固有のアサーション
        if tt.assert != nil {
            tt.assert(t)
        }
    })
}
```

## エラーハンドリングテスト規約

- **`errors.WithStack(err)`してエラーを返すだけの行に対するテストケースは不要**
- エラーをラップしてそのまま返すような処理は、テストの対象外とする
- テストが必要なエラーハンドリング：
  - エラー時の状態変更
  - エラーメッセージの生成・加工
  - エラー種別による分岐処理
  - リソースのクリーンアップ処理

## テストデータ管理

- 配置：レポジトリ直下の`testdata/`配下
- パッケージ構造に対応したサブディレクトリ
- ファイル参照：可能な限り`go:embed`を使用
- 命名規則：
  - 入力データ：`_input`で終わる
  - 出力検証：`_output`で終わる
