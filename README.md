# Hellog

指定した間隔で「Hello」を出力するサービス用Goアプリです。

- Unix系(Linux等): systemd サービスとして起動し、syslog へ出力
- Windows: サービスとして起動し、イベントログへ出力（イベントソース: `hellog` 既定）

## 使い方

実行時の設定は JSON ファイルから読み込みます（CLI フラグではなく設定駆動）。

## ビルド

```
go build -o dist/hellog ./cmd/hellog
GOOS=windows GOARCH=amd64 go build -o dist/hellog.exe ./cmd/hellog
```

## テスト

コアロジックはモックを用いたユニットテストで検証します。

```
go test ./...
```

## パッケージング (deb/rpm/msi)

標準的な手法として以下の設定ファイルを同梱しています。

- `packaging/nfpm.yaml`: nfpm を用いた `deb` / `rpm` 生成（systemdユニット/設定ファイル同梱）
- `packaging/systemd/hellog.service`: systemd ユニット
- `packaging/config/config.json`: Unix系のサンプル設定（`/etc/hellog/config.json` に配置）
- `packaging/windows/hellog.wxs`: WiX Toolset による `msi` 生成テンプレート（サービス登録含む）
- `packaging/windows/config.json`: Windows 用サンプル設定（`%ProgramData%\Hellog\config.json`）

### deb / rpm (nfpm)

1. バイナリを `dist/hellog` に配置するか、`nfpm.yaml` のパスを書き換えます。
2. nfpm を用いてパッケージを作成します。

```
nfpm pkg --config packaging/nfpm.yaml --packager deb
nfpm pkg --config packaging/nfpm.yaml --packager rpm
```

### msi (WiX)

1. WiX Toolset(heat/candle/light) をインストール
2. `packaging/windows/hellog.wxs` を使ってMSIを作成

```
candle -dVersion=0.1.0 -dBinPath=dist\\hellog.exe -out build\\ packaging\\windows\\hellog.wxs
light -ext WixUIExtension -o dist\\hellog.msi build\\hellog.wixobj
```

> 注意: Windowsのイベントソースの登録には管理者権限が必要な場合があります。インストーラや初回起動時に権限昇格が必要になることがあります。

## 標準準拠/構成

- `cmd/hellog`: サービス実行エントリ（Windowsは svc / Unixは通常プロセス）
- `internal/app/runner`: 間隔制御と実行ループ
- `internal/platform`: プラットフォーム依存のロガー実装（`!windows` は syslog、`windows` はイベントログ）
- `internal/config`: 設定ファイル(JSON)読み込み

## 配置と起動

### Linux (systemd)
- バイナリ: `/usr/bin/hellog`
- 設定: `/etc/hellog/config.json`（nfpm の `config_files` で導入）
- ユニット: `/lib/systemd/system/hellog.service`

コマンド:
```
sudo systemctl daemon-reload
sudo systemctl enable hellog
sudo systemctl start hellog
sudo systemctl status hellog
```

### Windows (サービス)
- バイナリ: `C:\Program Files\hellog\hellog.exe`
- 設定: `C:\ProgramData\Hellog\config.json`
- MSI インストールでサービス `hellog` が自動登録/起動

サービス管理例:
```
sc stop hellog
sc start hellog
sc query hellog
```
