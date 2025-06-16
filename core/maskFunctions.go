package core

func Mask0(row, col int) bool { return (row+col)%2 == 0 }

func Mask1(row, col int) bool { return row%2 == 0 }

func Mask2(row, col int) bool { return col%3 == 0 }

func Mask3(row, col int) bool { return (row+col)%3 == 0 }

func Mask4(row, col int) bool { return (row/2+col/3)%2 == 0 }

func Mask5(row, col int) bool { return ((row*col)%2)+((row*col)%3) == 0 }

func Mask6(row, col int) bool { return (((row*col)%2)+((row*col)%3))%2 == 0 }

func Mask7(row, col int) bool { return (((row+col)%2)+((row*col)%3))%2 == 0 }

var maskPatterns = []func(row, col int) bool{
	Mask0, Mask1, Mask2, Mask3, Mask4, Mask5, Mask6, Mask7,
}
