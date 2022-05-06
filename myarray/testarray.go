package testArrayPkg

func numIslands(grid [][]byte) int {
    rows := len(grid)
    cols := len(grid[0])
    var count int
    
    for r :=0; r<rows; r++ {
        for c :=0; c<cols; c++{
            if (grid[r][c] == '1'){
                count ++
                markByDFS(grid, r, c, rows, cols)
            }
        }
    }
	return count
}

func markByDFS(grid [][]byte, r int, c int, rows int, cols int){
    if (r < 0 || c <0 || r >=rows || c >=cols){
        return
    }
    if (grid[r][c] == '1'){
        grid[r][c] = 2
        markByDFS(grid, r-1, c, rows, cols)
        markByDFS(grid, r, c-1, rows, cols)
        markByDFS(grid, r+1, c, rows, cols)
        markByDFS(grid, r, c+1, rows, cols)
    }
}