
# Tetris in ebiten

このプログラムは、Go言語とEbitenライブラリを使用して簡単なテトリスゲームを作成します。

## インポート

```go
import (
	"image/color"
	"math/rand"
	"time"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)
```

- `image/color`: 色を扱うためのパッケージ
- `math/rand`: 乱数を生成するためのパッケージ
- `time`: 時間を扱うためのパッケージ
- `github.com/hajimehoshi/ebiten/v2`: ゲーム開発のためのライブラリ
- `github.com/hajimehoshi/ebiten/v2/ebitenutil`: Ebitenのユーティリティ関数

## 定数と型の定義

```go
const (
	boardWidth  = 10
	boardHeight = 20
	cellSize    = 20
)
```

- `boardWidth`: ボードの幅（列数）
- `boardHeight`: ボードの高さ（行数）
- `cellSize`: セル（ブロック）のサイズ（ピクセル）

```go
type Color int

const (
	Empty  Color = iota
	Red
	Blue
	Green
	Yellow
	Purple
	Orange
	Cyan
)
```

- `Color`: セルの色を表す列挙型
- `Empty, Red, Blue, ...`: 色を表す定数

## テトリスミノ（テトロミノ）の形状

```go
var tetrominoShapes = [][][][]Color{ ... }
```

- `tetrominoShapes`: テトリスミノの形状を定義する多次元配列

## テトリスミノとゲーム管理の構造体

```go
type Tetromino struct { ... }
type Game struct { ... }
```

- `Tetromino`: テトリスミノの情報を格納する構造体
- `Game`: ゲーム全体の状態を管理する構造体

## 主要な関数

- `NewTetromino()`: 新しいテトリスミノを生成する関数
- `collision(t *Tetromino, x, y int) bool`: 衝突判定を行う関数
- `placeTetromino()`: テトリスミノをボード上に配置する関数
- `clearLines()`: 完成したラインを消去する関数
- `Update() error`: ゲームの状態を更新する関数
- `Draw(screen *ebiten.Image)`: ゲームの描画を行う関数
- `drawCell(screen *ebiten.Image, x, y int, cell Color)`: セルを描画する関数

## メイン関数

```go
func main() {
	rand.Seed(time.Now().UnixNano())
	game := &Game{
		currentTetromino: NewTetromino(),
		dropInterval:     60,
	}
	ebiten.RunGame(game)
}
```

プログラムのエントリーポイントです。ゲームループを開始します。
