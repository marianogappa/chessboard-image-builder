# chessboard-image-builder

## Usage

```
go run main.go -fen 8/b7/8/8/7K/1R6/2n5/N2k4
```

You can change the piece and background images, but then you'll need to reconfigure these variables:

```
		// Set these parameters manually while running with -debug flag to match the grid
		boardMinX  = 40
		boardMinY  = 15
		boardMaxX  = 783
		boardMaxY  = 760
		cellWidth  = 93
		cellHeight = 93
```

This is not exactly a general-purpose tool; I wrote it to quickly regenerate all starting board images for a book I'm co-authoring.
But if it's useful for you, let me know.
