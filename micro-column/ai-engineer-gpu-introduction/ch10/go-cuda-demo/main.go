package main

import (
	"fmt"
	"math"
	"time"
	"unsafe"
	"runtime"

	"go-cuda-demo/cuda" // 引入我们刚才写的包
)

func main() {
	// 锁定 OS 线程
    // CUDA Context 是线程绑定的。如果不锁，Go runtime 可能会把当前的 Goroutine
    // 调度到另一个没有 Context 的 OS 线程上，导致 "CUDA Error: 201" (Invalid Context)。
    runtime.LockOSThread()
    defer runtime.UnlockOSThread()

	// --- 1. 初始化环境 ---
	err := cuda.LoadCudaLibrary()
	if err != nil {
		panic(err)
	}
	cuda.Init()
	
	fmt.Println("CUDA 初始化成功")
	dev := cuda.GetDevice(0) // 获取 0 号卡
	_ = cuda.CreateContext(dev) // 创建上下文
    // 注意：真实场景需要处理 Context 的 Destroy，但在 main 退出时系统会自动回收

	// --- 2. 加载 Kernel ---
	fmt.Println("加载 Kernel...")
	mod := cuda.LoadModule("./kernel/vector_add.ptx")
	kernel := cuda.GetFunction(mod, "vectorAdd")

	// --- 3. 准备数据 (Host) ---
	const N = 10_000_000 // 1千万
	const bytes = N * 4  // float32 占 4 字节
	
	h_A := make([]float32, N)
	h_B := make([]float32, N)
	h_C := make([]float32, N)

	for i := 0; i < N; i++ {
		h_A[i] = 2.0
		h_B[i] = 3.0
	}

	// --- 4. 显存分配与搬运 (H2D) ---
	fmt.Printf("搬运 %d MB 数据到 GPU...\n", bytes/1024/1024*2)
	d_A := cuda.Malloc(bytes)
	d_B := cuda.Malloc(bytes)
	d_C := cuda.Malloc(bytes)
	defer cuda.Free(d_A)
	defer cuda.Free(d_B)
	defer cuda.Free(d_C)

    startH2D := time.Now()
	cuda.MemcpyHtoD(d_A, unsafe.Pointer(&h_A[0]), bytes)
	cuda.MemcpyHtoD(d_B, unsafe.Pointer(&h_B[0]), bytes)
    fmt.Printf("H2D 耗时: %v\n", time.Since(startH2D))

	// --- 5. 兵力规划 (Grid & Block) ---
	blockSize := 256
	// 向上取整公式: (N + block - 1) / block
	gridSize := (N + blockSize - 1) / blockSize
	
	fmt.Printf("启动 Kernel: Grid=%d, Block=%d\n", gridSize, blockSize)

	// --- 6. 发射 Kernel (Launch) ---
    // 预热
    cuda.Launch(kernel, uint32(gridSize), uint32(blockSize), d_A, d_B, d_C, N)
    cuda.Synchronize()

    // 计时
	startKernel := time.Now()
    // 注意参数顺序要和 .cu 文件里的一致：A, B, C, N
	cuda.Launch(kernel, uint32(gridSize), uint32(blockSize), d_A, d_B, d_C, N)
	cuda.Synchronize() // 必须同步！
	fmt.Printf("GPU 计算耗时: %v\n", time.Since(startKernel))

	// --- 7. 搬回结果 (D2H) ---
    startD2H := time.Now()
	cuda.MemcpyDtoH(unsafe.Pointer(&h_C[0]), d_C, bytes)
    fmt.Printf("D2H 耗时: %v\n", time.Since(startD2H))

	// --- 8. 验证 ---
	success := true
    // 抽样验证前5个和后5个
	for i := 0; i < 5; i++ {
		if math.Abs(float64(h_C[i]-5.0)) > 1e-5 {
			success = false
            break
		}
	}
    // 检查末尾
    if math.Abs(float64(h_C[N-1]-5.0)) > 1e-5 {
        success = false
    }

	if success {
		fmt.Println("✅ 验证成功！2.0 + 3.0 = 5.0")
	} else {
		fmt.Println("❌ 验证失败！")
	}
}
