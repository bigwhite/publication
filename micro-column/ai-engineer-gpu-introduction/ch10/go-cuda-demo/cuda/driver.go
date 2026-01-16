package cuda

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
)

// 定义 CUDA 的基本类型
type CUdevice int
type CUcontext uintptr
type CUmodule uintptr
type CUfunction uintptr
type CUdeviceptr uint64 // 显存地址通常是 64 位

// 定义 CUDA API 的函数签名变量
var (
	cuInit             func(flags uint32) int32
	cuDeviceGet        func(device *CUdevice, ordinal int) int32
	cuCtxCreate        func(pctx *CUcontext, flags uint32, dev CUdevice) int32
	cuModuleLoad       func(module *CUmodule, fname string) int32
	cuModuleGetFunction func(hfunc *CUfunction, hmod CUmodule, name string) int32
	cuMemAlloc         func(dptr *CUdeviceptr, bytesize uint64) int32
	cuMemFree          func(dptr CUdeviceptr) int32
	cuMemcpyHtoD       func(dstDevice CUdeviceptr, srcHost unsafe.Pointer, ByteCount uint64) int32
	cuMemcpyDtoH       func(dstHost unsafe.Pointer, srcDevice CUdeviceptr, ByteCount uint64) int32
	cuLaunchKernel     func(f CUfunction,
		gridDimX, gridDimY, gridDimZ uint32,
		blockDimX, blockDimY, blockDimZ uint32,
		sharedMemBytes uint32,
		hStream uintptr,
		kernelParams unsafe.Pointer,
		extra unsafe.Pointer) int32
	cuCtxSynchronize func() int32
)

// LoadCudaLibrary 加载 libcuda.so 并绑定函数
func LoadCudaLibrary() error {
	libName := "libcuda.so"
	if runtime.GOOS == "windows" {
		libName = "nvcuda.dll"
	}
    // 1. 动态加载动态库
	lib, err := purego.Dlopen(libName, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("failed to load cuda driver: %v. Make sure NVIDIA driver is installed", err)
	}

    // 2. 绑定符号 (Symbol)
	purego.RegisterLibFunc(&cuInit, lib, "cuInit")
	purego.RegisterLibFunc(&cuDeviceGet, lib, "cuDeviceGet")
	purego.RegisterLibFunc(&cuCtxCreate, lib, "cuCtxCreate_v2") // 注意：现代驱动通常用 v2 版本
	purego.RegisterLibFunc(&cuModuleLoad, lib, "cuModuleLoad")
	purego.RegisterLibFunc(&cuModuleGetFunction, lib, "cuModuleGetFunction")
	purego.RegisterLibFunc(&cuMemAlloc, lib, "cuMemAlloc_v2")
	purego.RegisterLibFunc(&cuMemFree, lib, "cuMemFree_v2")
	purego.RegisterLibFunc(&cuMemcpyHtoD, lib, "cuMemcpyHtoD_v2")
	purego.RegisterLibFunc(&cuMemcpyDtoH, lib, "cuMemcpyDtoH_v2")
	purego.RegisterLibFunc(&cuLaunchKernel, lib, "cuLaunchKernel")
	purego.RegisterLibFunc(&cuCtxSynchronize, lib, "cuCtxSynchronize")

	return nil
}

// 辅助函数：检查 CUDA 错误码
func Check(err int32) {
	if err != 0 {
		panic(fmt.Sprintf("CUDA Error: %d", err))
	}
}

// 封装一些易用的 API
func Init() {
	Check(cuInit(0))
}

func GetDevice(ordinal int) CUdevice {
	var dev CUdevice
	Check(cuDeviceGet(&dev, ordinal))
	return dev
}

func CreateContext(dev CUdevice) CUcontext {
	var ctx CUcontext
	Check(cuCtxCreate(&ctx, 0, dev))
	return ctx
}

func LoadModule(path string) CUmodule {
	var mod CUmodule
	Check(cuModuleLoad(&mod, path))
	return mod
}

func GetFunction(mod CUmodule, name string) CUfunction {
	var fn CUfunction
	Check(cuModuleGetFunction(&fn, mod, name))
	return fn
}

func Malloc(size uint64) CUdeviceptr {
	var ptr CUdeviceptr
	Check(cuMemAlloc(&ptr, size))
	return ptr
}

func Free(ptr CUdeviceptr) {
	Check(cuMemFree(ptr))
}

func MemcpyHtoD(dst CUdeviceptr, src unsafe.Pointer, size uint64) {
	Check(cuMemcpyHtoD(dst, src, size))
}

func MemcpyDtoH(dst unsafe.Pointer, src CUdeviceptr, size uint64) {
	Check(cuMemcpyDtoH(dst, src, size))
}

func Synchronize() {
	Check(cuCtxSynchronize())
}

// LaunchKernel 是最关键的封装
func Launch(fn CUfunction, gridDim, blockDim uint32, args ...interface{}) {
    // 构造参数指针数组
    // CUDA Driver API 要求参数必须以 (void*[]) 的形式传递
    // 即：一个指针数组，数组里的每个元素是指向实际参数的指针
	kernelParams := make([]unsafe.Pointer, len(args))
	for i, arg := range args {
		switch v := arg.(type) {
		case CUdeviceptr:
			kernelParams[i] = unsafe.Pointer(&v)
		case *int32:
			kernelParams[i] = unsafe.Pointer(v)
        case *float32:
            kernelParams[i] = unsafe.Pointer(v)
        // 注意：这里传递的是 int 的值（作为指针内容），如果是 int 标量参数，需要取地址
        case int:
             // 这种写法需要注意生命周期，但在当前栈帧是安全的
             val := v
             kernelParams[i] = unsafe.Pointer(&val)
		default:
			panic("unsupported argument type")
		}
	}

	Check(cuLaunchKernel(
		fn,
		gridDim, 1, 1,    // Grid 维度
		blockDim, 1, 1,   // Block 维度
		0,                // Shared Memory 大小
		0,                // Stream (0 表示默认流)
		unsafe.Pointer(&kernelParams[0]), // 参数数组的指针
		nil,
	))
}
