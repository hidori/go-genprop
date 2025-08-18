# テスト規約

## 作業時の必須チェック

- **テストコード変更前後に必ず行う**
  - `make format`を実行する
  - `make test`を実行する

## 基本方針

- 原則、テーブルドリブンテストを使用する
- テスト関数名：`Test` + 機能名とする
- `t.Parallel()`を積極使用する（有害な副作用がない限り）
- `testify/assert`と`testify/require`を使用する
- **テストコードには意図を明確化するためのコメント記述を許容する**
- コメントは英語で記述する

## エラー処理

- **エラー返却時に必ず行う**
  - `errors.WithStack()`または`errors.Wrap()`/`errors.Wrapf()`でスタックトレースを付与する
- コンテキスト情報が必要な場合：`errors.Wrap()`/`errors.Wrapf()`を使用する
- コンテキスト情報が不要な場合：`errors.WithStack()`を使用する
- **`errors.WithStack(err)`してエラーを返すだけの行に対するテストケースは不要**
- エラーをラップしてそのまま返すような処理は、テストの対象外とする

## テスト記述

### 一時ファイル管理

- **一時ファイル作成時に必ず行う**
  - プロジェクト直下の`tmp/`ディレクトリに作成する
- **作業完了後に必ず行う**
  - `make clean`で一時ファイルを削除する
- 例: `go test -coverprofile=tmp/coverage.out`

### テストケース名

| 種類 | 条件 | 命名規則 |
|------|------|----------|
| 正常系 | `err=nil`を返す | `success:` で開始 |
| 異常系 | `err!=nil`を返す | `failure:` で開始 |

### テストコメント

- **テストコメントの用途**
  - テストの意図を説明する
  - 複雑なアサーションを説明する
  - 特殊な条件設定を説明する
- 例: `// Test error handling when user is not found`

### テストケース固有のアサーション

- **テストケース固有のアサーションが必要な場合、テストケースの匿名構造体に`assert func(t *testing.T)`メンバを追加**
- **アサーションの実行**
  - 共通のアサーション: `t.Run`内で直接実行する
  - 固有のアサーション: `assert`関数内で実行する
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

- テストが必要なエラーハンドリング：
  - エラー時の状態変更
  - エラーメッセージの生成・加工
  - エラー種別による分岐処理
  - リソースのクリーンアップ処理

## テストデータ管理

- 配置：プロジェクト直下の`testdata/`配下
- パッケージ構造に対応したサブディレクトリ
- ファイル参照：可能な限り`go:embed`を使用する
- 命名規則：
  - 入力データ：`_input`で終わる
  - 出力検証：`_output`で終わる
