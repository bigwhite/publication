// kernel/vector_add.cu
extern "C" { // 必须加上 extern "C"，防止 C++ 名字修饰(Mangling)，方便 Go 调用

// 这是一个标准的 CUDA Kernel
__global__ void vectorAdd(const float *A, const float *B, float *C, int numElements) {
    // 计算全局 ID
    int i = blockDim.x * blockIdx.x + threadIdx.x;

    // 边界检查 + 干活
    if (i < numElements) {
        C[i] = A[i] + B[i];
    }
}

}
