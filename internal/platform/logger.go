package platform

// HelloLogger は「Hello」を出力するための最小インタフェースです。
type HelloLogger interface {
    LogHello() error
    Close() error
}

// NewLogger は現在のプラットフォーム向けの HelloLogger を生成します。
// tag は syslog のタグ名や Windows のイベントソース名として使われます。
func NewLogger(tag string) (HelloLogger, error) { // 実体はビルドタグ付きファイルに実装
    return newPlatformLogger(tag)
}

