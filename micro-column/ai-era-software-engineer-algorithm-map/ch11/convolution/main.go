package main

import "fmt"

// applyConvolution 对图像进行 3x3 卷积操作
// image: 输入矩阵
// kernel: 3x3 卷积核
func applyConvolution(image [][]int, kernel [][]int) [][]int {
	rows := len(image)
	cols := len(image[0])
	kSize := len(kernel) // 假设是正方形 3
	pad := kSize / 2     // padding = 1

	// 1. Padding: 创建一个更大的矩阵，四周补 0
	// 这一步是为了让卷积核能扫描到边缘像素
	paddedRows := rows + 2*pad
	paddedCols := cols + 2*pad
	paddedImg := make([][]int, paddedRows)
	for i := range paddedImg {
		paddedImg[i] = make([]int, paddedCols)
	}

	// 填充原图数据到中心
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			paddedImg[i+pad][j+pad] = image[i][j]
		}
	}

	// 2. 卷积计算
	// 输出尺寸与原图一致（因为步长为1且加了padding）
	output := make([][]int, rows)
	for i := range output {
		output[i] = make([]int, cols)
	}

	// 滑动窗口
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			sum := 0
			// 卷积核运算
			for ki := 0; ki < kSize; ki++ {
				for kj := 0; kj < kSize; kj++ {
					// 图像坐标需要加上 padding 偏移
					// 卷积核坐标直接映射
					val := paddedImg[i+ki][j+kj]
					weight := kernel[ki][kj]
					sum += val * weight
				}
			}
			output[i][j] = sum
		}
	}

	return output
}

func main() {
	// 模拟一张 5x5 的图片（简单的像素值）
	image := [][]int{
		{10, 10, 10, 0, 0},
		{10, 10, 10, 0, 0},
		{10, 10, 10, 0, 0},
		{10, 10, 10, 0, 0},
		{10, 10, 10, 0, 0},
	}

	// 定义一个 "垂直边缘检测" 卷积核 (Vertical Edge Detector)
	// 左边是正，右边是负，中间是0。遇到左右差异大的地方，结果会很大。
	kernel := [][]int{
		{1, 0, -1},
		{1, 0, -1},
		{1, 0, -1},
	}

	featureMap := applyConvolution(image, kernel)

	fmt.Println("--- Feature Map (Edge Detection) ---")
	for _, row := range featureMap {
		fmt.Println(row)
	}
	// 预期输出：中间那一列数值会很大（30），检测到了边缘
}
