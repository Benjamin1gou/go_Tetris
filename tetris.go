package main

import (
	"fmt"
	"image/color" // 色を扱うパッケージ
	"math/rand"   // 乱数を扱うパッケージ
	"time"        // 時間を扱うパッケージ

	"github.com/hajimehoshi/ebiten/v2"            // ゲーム開発ライブラリ
	"github.com/hajimehoshi/ebiten/v2/ebitenutil" // Ebitenのユーティリティ
)

// ボードの設定定数
const (
	boardWidth  = 10 // ボードの幅
	boardHeight = 20 // ボードの高さ
	cellSize    = 20 // セルのサイズ（ピクセル）
)

// セルの色を表す型
type Color int

// セルの色を表す定数
const (
	Empty  Color = iota // 空のセル
	Red                 // 赤色
	Blue                // 青色
	Green               // 緑色
	Yellow              // 黄色
	Purple              // 紫色
	Orange              // オレンジ色
	Cyan                // シアン色
)

// テトリスミノの形状を定義
var tetrominoShapes = [][][][]Color{
	// I形状
	{
		{
			{Cyan, Cyan, Cyan, Cyan},
		},
		{
			{Cyan},
			{Cyan},
			{Cyan},
			{Cyan},
		},
	},
	// O形状
	{
		{
			{Yellow, Yellow},
			{Yellow, Yellow},
		},
	},
	// T形状
	{
		{
			{0, Purple, 0},
			{Purple, Purple, Purple},
		},
		{
			{0, Purple},
			{Purple, Purple},
			{0, Purple},
		},
		{
			{Purple, Purple, Purple},
			{0, Purple, 0},
		},
		{
			{Purple, 0},
			{Purple, Purple},
			{Purple, 0},
		},
	},
	// S形状
	{
		{
			{0, Green, Green},
			{Green, Green, 0},
		},
		{
			{Green, 0},
			{Green, Green},
			{0, Green},
		},
	},
	// Z形状
	{
		{
			{Red, Red, 0},
			{0, Red, Red},
		},
		{
			{0, Red},
			{Red, Red},
			{Red, 0},
		},
	},
	// J形状
	{
		{
			{Blue, 0, 0},
			{Blue, Blue, Blue},
		},
		{
			{Blue, Blue},
			{Blue, 0},
			{Blue, 0},
		},
		{
			{Blue, Blue, Blue},
			{0, 0, Blue},
		},
		{
			{0, Blue},
			{0, Blue},
			{Blue, Blue},
		},
	},
	// L形状
	{
		{
			{0, 0, Orange},
			{Orange, Orange, Orange},
		},
		{
			{Orange, 0},
			{Orange, 0},
			{Orange, Orange},
		},
		{
			{Orange, Orange, Orange},
			{Orange, 0, 0},
		},
		{
			{Orange, Orange},
			{0, Orange},
			{0, Orange},
		},
	},
}

// テトリスミノを表す構造体
type Tetromino struct {
	x, y         int       // テトリスミノの座標
	shapeIndex   int       // 形状のインデックス
	rotation     int       // 回転の状態
	currentShape [][]Color // 現在の形状
}

// 新しいテトリスミノを生成する関数
func NewTetromino() *Tetromino {
	t := &Tetromino{
		x:          boardWidth/2 - 2,                // 初期位置をボードの中央に設定
		y:          0,                               // 初期位置のy座標は0
		shapeIndex: rand.Intn(len(tetrominoShapes)), // ランダムな形状を選ぶ
	}
	t.currentShape = tetrominoShapes[t.shapeIndex][t.rotation] // 現在の形状を設定
	return t
}

// ゲーム全体を管理する構造体
type Game struct {
	board            [boardHeight][boardWidth]Color // ゲームボード
	currentTetromino *Tetromino                     // 現在操作しているテトリスミノ
	nextTetromino    *Tetromino                     // 次に生成されるテトリスミノ
	currentFrame     int                            // フレーム数
	dropInterval     int                            // インターバル
}

// テトリスミノが他のブロックや壁と衝突するか判定する関数
func (g *Game) collision(t *Tetromino, x, y int) bool {
	for rowIdx, row := range t.currentShape {
		for colIdx, cell := range row {
			if cell != Empty {
				if x+colIdx < 0 || x+colIdx >= boardWidth || y+rowIdx >= boardHeight {
					return true
				}
				if g.board[y+rowIdx][x+colIdx] != Empty {
					return true
				}
			}
		}
	}
	return false
}

// テトリスミノをボード上に配置する関数
func (g *Game) placeTetromino() {
	for rowIdx, row := range g.currentTetromino.currentShape {
		for colIdx, cell := range row {
			if cell != Empty {
				g.board[g.currentTetromino.y+rowIdx][g.currentTetromino.x+colIdx] = cell
			}
		}
	}
	g.currentTetromino = g.nextTetromino // 現在の「次のテトロミノ」を「現在のテトロミノ」にセット
	g.nextTetromino = NewTetromino()     // 新しい「次のテトロミノ」を生成
	if g.collision(g.currentTetromino, g.currentTetromino.x, g.currentTetromino.y) {
		// 新しいテトリスミノが配置できない場合はゲームをリセット
		g.board = [boardHeight][boardWidth]Color{}
	}
}

// ラインが完成したら消去する関数
func (g *Game) clearLines() {
	for y := 0; y < boardHeight; y++ {
		full := true // 行が全て埋まっているかのフラグ
		for x := 0; x < boardWidth; x++ {
			if g.board[y][x] == Empty {
				full = false
				break
			}
		}
		if full { // 行が全て埋まっていた場合
			for yy := y; yy > 0; yy-- {
				for xx := 0; xx < boardWidth; xx++ {
					g.board[yy][xx] = g.board[yy-1][xx]
				}
			}
			for xx := 0; xx < boardWidth; xx++ {
				g.board[0][xx] = Empty
			}
		}
	}
}

// ゲームの状態を更新する関数（毎フレーム呼ばれる）
func (g *Game) Update() error {

	g.currentFrame++
	if g.currentFrame >= g.dropInterval {
		if !g.collision(g.currentTetromino, g.currentTetromino.x, g.currentTetromino.y+1) {
			g.currentTetromino.y++
		} else {
			g.placeTetromino()
			g.clearLines()
		}
		g.currentFrame = 0
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) && !g.collision(g.currentTetromino, g.currentTetromino.x+1, g.currentTetromino.y) {
		g.currentTetromino.x++ // 右キーが押されたら、テトリスミノを右に1移動
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && !g.collision(g.currentTetromino, g.currentTetromino.x-1, g.currentTetromino.y) {
		g.currentTetromino.x-- // 左キーが押されたら、テトリスミノを左に1移動
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) && !g.collision(g.currentTetromino, g.currentTetromino.x, g.currentTetromino.y+1) {
		g.currentTetromino.y++ // 下キーが押されたら、テトリスミノを下に1移動
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.currentTetromino.rotation = (g.currentTetromino.rotation + 1) % len(tetrominoShapes[g.currentTetromino.shapeIndex]) // 上キーが押されたら、テトリスミノを回転
		g.currentTetromino.currentShape = tetrominoShapes[g.currentTetromino.shapeIndex][g.currentTetromino.rotation]         // 回転後の形状を現在の形状に設定
		if g.collision(g.currentTetromino, g.currentTetromino.x, g.currentTetromino.y) {
			g.currentTetromino.rotation-- // 回転後に衝突する場合、回転を1つ前に戻す
			if g.currentTetromino.rotation < 0 {
				g.currentTetromino.rotation = len(tetrominoShapes[g.currentTetromino.shapeIndex]) - 1 // 回転インデックスが負になった場合、最後の回転状態に設定
			}
			g.currentTetromino.currentShape = tetrominoShapes[g.currentTetromino.shapeIndex][g.currentTetromino.rotation] // 回転を戻した後の形状を現在の形状に設定
		}
	}
	if g.collision(g.currentTetromino, g.currentTetromino.x, g.currentTetromino.y+1) {
		g.placeTetromino() // 下に移動すると衝突する場合、テトリスミノをボードに固定
		g.clearLines()     // 完成したラインを消去
	}
	return nil
}

// ゲームの描画を行う関数
func (g *Game) Draw(screen *ebiten.Image) {
	for y, row := range g.board {
		for x, cell := range row {
			drawCell(screen, x, y, cell) // ボード上の各セルを描画
		}
	}
	for y, row := range g.currentTetromino.currentShape {
		for x, cell := range row {
			if cell != Empty {
				drawCell(screen, g.currentTetromino.x+x, g.currentTetromino.y+y, cell) // 現在動かしているテトリスミノの各セルを描画
			}
		}
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Next Tetromino: %+v", g.nextTetromino.currentShape))
	// 次のテトロミノを描画
	offsetX := (boardWidth + 2) * cellSize // 例: ボードの右側に表示
	offsetY := 0                           // 上端に表示
	for y, row := range g.nextTetromino.currentShape {
		for x, cell := range row {
			if cell != Empty {
				drawCell(screen, offsetX/cellSize+x, offsetY/cellSize+y, cell)
			}
		}
	}
}

// セルを描画する関数
func drawCell(screen *ebiten.Image, x, y int, cell Color) {
	col := color.RGBA{255, 255, 255, 255} // デフォルトは白色
	switch cell {                         // セルの色に応じてRGBA値を設定
	case Red:
		col = color.RGBA{255, 0, 0, 255}
	case Blue:
		col = color.RGBA{0, 0, 255, 255}
	case Green:
		col = color.RGBA{0, 255, 0, 255}
	case Yellow:
		col = color.RGBA{255, 255, 0, 255}
	case Purple:
		col = color.RGBA{128, 0, 128, 255}
	case Orange:
		col = color.RGBA{255, 165, 0, 255}
	case Cyan:
		col = color.RGBA{0, 255, 255, 255}
	}
	x0 := x * cellSize                                                                         // 描画するセルの左上のx座標
	y0 := y * cellSize                                                                         // 描画するセルの左上のy座標
	x1 := (x + 1) * cellSize                                                                   // 描画するセルの右下のx座標
	y1 := (y + 1) * cellSize                                                                   // 描画するセルの右下のy座標
	ebitenutil.DrawRect(screen, float64(x0), float64(y0), float64(x1-x0), float64(y1-y0), col) // セルを描画
}

// ゲーム画面のサイズを設定する関数
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return boardWidth * cellSize, boardHeight * cellSize // ゲーム画面のサイズをセルサイズとボードサイズから計算
}

// メイン関数（プログラムのエントリーポイント）
func main() {
	rand.Seed(time.Now().UnixNano()) // 乱数のシードを設定
	game := &Game{                   // Gameの新しいインスタンスを作成
		currentTetromino: NewTetromino(), // 初期のテトリスミノを設定
		nextTetromino:    NewTetromino(), // 次に生成されるミノを設定
		dropInterval:     60,             // 落下インターバルのフレーム数を設定
	}
	ebiten.RunGame(game) // ゲームループを開始
}
