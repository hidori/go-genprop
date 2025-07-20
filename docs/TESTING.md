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

### アサーション方針

- **testify/assert.*で表現可能なアサーションは、期待値をテストケース構造体のフィールドとして定義し、`assert.Equal()`等で検証する**
- **testify/assert.*で表現不可能な複雑なアサーションのみ、`assert func(t *testing.T)`を使用する**

### testify/assert と testify/require の使い分け

- **require.***: テストの継続に必要な前提条件の検証に使用。失敗時にテストを即座に停止
  - nil チェック: `require.NotNil(t, obj)` - 後続でobjを使用する場合
  - エラーチェック: `require.NoError(t, err)` - 後続処理がエラーに依存する場合
  - 必須の戻り値チェック: `require.True(t, condition)` - 後続処理の前提条件
- **assert.***: テストの継続に影響しない検証に使用。失敗してもテストを継続
  - 値の比較: `assert.Equal(t, want, got)`
  - 補助的なチェック: `assert.Contains(t, result, substring)`
  - 複数の独立した検証

### testify/assert.*で表現可能な例

- 値の等価性比較: `assert.Equal(t, want, got)`
- 型の検証: `assert.IsType(t, expectedType, got)`
- nil検証: `assert.NotNil(t, got)`
- エラー検証: `assert.NoError(t, err)` / `assert.Error(t, err)`
- 文字列の包含: `assert.Contains(t, got, substring)`
- 真偽値: `assert.True(t, condition)` / `assert.False(t, condition)`

### assert func(t *testing.T)が必要な例

- 複数の条件を組み合わせた複雑な検証
- 副作用の検証（ファイル作成、ログ出力等）
- 動的な期待値の計算が必要な場合
- モックの呼び出し回数・引数の検証

### 使用例

```go
tests := []struct {
    name string
    input string
    want string
    wantType interface{}
    assert func(t *testing.T, result string) // 複雑なアサーションのみ
}{
    {
        name: "success: simple case",
        input: "test",
        want: "expected",
        wantType: "",
        assert: nil, // testify/assert.*で十分
    },
    {
        name: "success: complex validation case",
        input: "special",
        want: "result",
        wantType: "",
        assert: func(t *testing.T, result string) {
            // 複雑な検証ロジック（testify/assert.*では表現困難）
            parts := strings.Split(result, "-")
            assert.Len(t, parts, 3)
            assert.True(t, len(parts[0]) > 0)
            assert.Regexp(t, `^\d+$`, parts[1])
        },
    },
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        result := targetFunction(tt.input)

        // 前提条件の検証（失敗時はテスト停止）
        require.NotNil(t, result)

        // 主要な検証（testify/assert.*で表現可能）
        assert.Equal(t, tt.want, result)
        assert.IsType(t, tt.wantType, result)

        // 複雑なアサーション（testify/assert.*では困難）
        if tt.assert != nil {
            tt.assert(t, result)
        }
    })
}
```

## エラーハンドリングテスト規約

- **`errors.WithStack(err)`してエラーを返すだけの行に対するテストケースは不要**
- エラーをラップしてそのまま返すような処理は、テストの対象外とする
- **テスト対象の関数内で生成・返されるエラーメッセージが特定の文字列を含むことが確実な場合、`assert.Contains(t, err.Error(), expectedString)`でエラーメッセージを検証する**
- テストが必要なエラーハンドリング：
  - エラー時の状態変更
  - エラーメッセージの生成・加工
  - エラー種別による分岐処理
  - リソースのクリーンアップ処理

### エラーメッセージ検証の指針

- **検証対象**: テスト対象の関数内で生成・返されるエラーメッセージ
  - 関数内で`errors.New()`や`fmt.Errorf()`で作成されるエラー
  - 関数内でカスタムエラー型を生成する場合
  - 関数内で独自のメッセージを追加してエラーをラップする場合

- **テーブル構造体のフィールド規約**:
  - **`wantErr bool`**: 全てのテストケースで必須設定（省略不可）
  - **`wantErrContains string`**: テスト対象の関数内でエラーメッセージを生成する場合のみ追加
  - **外部からのエラーをそのまま返すだけの関数では`wantErrContains`フィールドは不要**
  - エラーが発生しない場合：`wantErr: false`を明示的に設定
  - 内部でエラーメッセージを生成する場合：`wantErr: true, wantErrContains: "期待するエラーメッセージ"`を設定

- **判定基準**:
  - テスト対象の関数内に`errors.New()`、`fmt.Errorf()`、カスタムエラー型の生成がある → `wantErrContains`必要
  - 外部の関数呼び出しで発生したエラーをそのまま返すだけ → `wantErrContains`不要

```go
// エラーメッセージを内部生成する関数のテスト例
tests := []struct {
    name            string
    input           string
    wantErr         bool
    wantErrContains string // この関数が内部でエラーメッセージを生成するため必要
}{
    {
        name:            "success: valid input",
        input:           "valid",
        wantErr:         false,
        wantErrContains: "", // 正常系でも明示的に空文字列を設定
    },
    {
        name:            "failure: invalid input",
        input:           "invalid",
        wantErr:         true,
        wantErrContains: "file argument is required", // 内部生成エラーメッセージを検証
    },
}

// 外部エラーをそのまま返す関数のテスト例
tests := []struct {
    name    string
    input   string
    wantErr bool
    // wantErrContains は不要（外部エラーをそのまま返すのみ）
}{
    {
        name:    "success: valid input",
        input:   "valid",
        wantErr: false,
    },
    {
        name:    "failure: external error",
        input:   "invalid",
        wantErr: true, // エラーの有無のみ検証
    },
}

if tt.wantErr {
    require.Error(t, err)
    if tt.wantErrContains != "" {
        assert.Contains(t, err.Error(), tt.wantErrContains)
    }
} else {
    require.NoError(t, err)
}
```

## テストデータ管理

- 配置：レポジトリ直下の`testdata/`配下
- パッケージ構造に対応したサブディレクトリ
- ファイル参照：可能な限り`go:embed`を使用
- 命名規則：
  - 入力データ：`_input`で終わる
  - 出力検証：`_output`で終わる
