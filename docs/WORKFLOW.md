# 開発ワークフロー

## 重要な基本ルール

- **一時ファイルは必ず`tmp/`ディレクトリに作成し、作業完了後は`make clean`で削除**
- **コード変更前に必ず`make lint`と`make test`を実行**
- **コード変更後に必ず`make lint`、`make test`、`make example/run`を実行**
- **許可なく新規ファイルを作成しない**

## 基本ルール

- 明確な指示が無い限り、修正点のコミットを自発的に行う必要はない

## 開発フロー

以下に該当する変更では、必ずこのフローを厳守してください：

- Goソースコード（`.go`ファイル）の変更
- 設定ファイル（`Makefile`, `.golangci.yml`等）の変更
- CI/CD設定（`.github/workflows/`等）の変更
- コード生成に影響するファイル（テンプレート、スキーマ等）の変更

純粋なドキュメント（`README.md`, `DEVELOPMENT.md`等）やライセンスファイルの変更では不要です。

### 必須チェックリスト

#### 変更前（必須）

```bash
make lint   # コード品質チェック（golangci-lint）
make test   # 全テスト実行（レース条件検出付き）
```

#### コード変更

- 新ファイル作成時は事前承認が必要
- 適切なディレクトリ構造を遵守
- ファイル命名規則に従う

#### 変更後（必須）

```bash
make lint        # コード品質再チェック
make test        # 全テスト再実行
make example/run # サンプル動作確認
```

#### コミット前の最終確認

- [ ] すべてのテストが通過
- [ ] lintエラーがゼロ
- [ ] サンプルが正常実行

## プロジェクト構造とファイル配置

### ディレクトリ構造

- **docs/**: ドキュメント類
- **example/basic/**: 基本例
- **example/private-setter/**: private setter例
- **cmd/**: コマンド
- **public/**: 公開パッケージ
- **internal/app/**: 内部パッケージ
- **tmp/**: 一時ファイル（`.gitignore`対象、`make clean`で削除）

### 一時ファイル管理規約

- **一時ファイルは必ず`tmp/`ディレクトリに作成**
- **作業完了後は`make clean`で適切に削除**
- カバレッジファイル、ビルド成果物等の作業用ファイルが対象
- 例: `go test -coverprofile=tmp/coverage.out`

### ファイル命名規則

- パッケージ名と一致する適切なファイル名を使用
- テストファイル: `_test.go`
- モックファイル: `_mock.go`
- 例: `example/basic/user.go`, `example/private-setter/user.go`

## バージョン管理

- [Semantic Versioning 2.0.0](https://semver.org/)を使用
- 形式：`MAJOR.MINOR.PATCH`（例：`0.0.21`）
- プレリリース版：`MAJOR.MINOR.PATCH-prerelease`（例：`1.0.0-beta`）
- ビルドメタデータ：`MAJOR.MINOR.PATCH+build`（例：`1.0.0+build.1`）
- バージョン情報は`public/meta/version.txt`に記録
- `public/meta/version.go`の`GetVersion()`関数で"v"プレフィックス付きで提供
- 更新には`make version/patch`を使用（自動インクリメント、Git操作含む）

```bash
make version/patch
```
