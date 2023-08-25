package swift

import (
	"context"
	"os"
	"path"

	"github.com/khulnasoft-lab/go-dep-parser/pkg/swift/swift"
	"github.com/khulnasoft-lab/vul/pkg/fanal/analyzer"
	"github.com/khulnasoft-lab/vul/pkg/fanal/analyzer/language"
	"github.com/khulnasoft-lab/vul/pkg/fanal/types"

	"golang.org/x/xerrors"
)

func init() {
	analyzer.RegisterAnalyzer(&swiftLockAnalyzer{})
}

const (
	version = 1
)

// swiftLockAnalyzer analyzes Package.resolved files
type swiftLockAnalyzer struct{}

func (a swiftLockAnalyzer) Analyze(_ context.Context, input analyzer.AnalysisInput) (*analyzer.AnalysisResult, error) {
	p := swift.NewParser()
	res, err := language.Analyze(types.Swift, input.FilePath, input.Content, p)
	if err != nil {
		return nil, xerrors.Errorf("%s parse error: %w", input.FilePath, err)
	}
	return res, nil
}

func (a swiftLockAnalyzer) Required(filePath string, _ os.FileInfo) bool {
	return path.Base(filePath) == types.SwiftResolved
}

func (a swiftLockAnalyzer) Type() analyzer.Type {
	return analyzer.TypeSwift
}

func (a swiftLockAnalyzer) Version() int {
	return version
}
