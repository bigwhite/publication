from numba import cuda
import numpy as np

# --- 1. 定义 Kernel (作战手册) ---
# @cuda.jit 装饰器告诉 Numba：这不是普通函数，请把它编译成 GPU 机器码！
@cuda.jit
def vector_add_kernel(A, B, C):
    # 这里的代码，会被成千上万个线程同时执行！

    # 关键点：我是谁？
    # Numba 提供了 cuda.grid(1) 这个语法糖
    # 它等价于：i = blockIdx.x * blockDim.x + threadIdx.x
    # 也就是算出当前线程在整个大军中的“全局唯一编号”
    i = cuda.grid(1)

    # 边界检查：防止兵力(线程数)多于任务(数据量)时越界
    if i < C.size:
        # 干活！每个线程只算这一个数
        C[i] = A[i] + B[i]

def main():
    # --- 2. 准备数据 (Host 端) ---
    N = 10_000_000  # 向量长度：1千万
    print(f"正在准备 {N} 个数据...")

    # 在 CPU 内存中创建数组 (float32 是 GPU 最常用的类型)
    A = np.ones(N, dtype=np.float32) * 2  # 全是 2
    B = np.ones(N, dtype=np.float32) * 3  # 全是 3
    C = np.zeros(N, dtype=np.float32)     # 结果容器

    # --- 3. 搬运数据 (Host -> Device) ---
    # 这一步通过 PCIe 总线，比较慢
    print("将数据从 CPU 搬运到 GPU...")
    d_A = cuda.to_device(A)
    d_B = cuda.to_device(B)
    d_C = cuda.device_array(N, dtype=np.float32) # 只分配空间，不拷贝

    # --- 4. 兵力规划 (Grid & Block) ---
    # Block: 设置每个班 256 人 (通常是 32 的倍数)
    threads_per_block = 256

    # Grid: 计算需要多少个班？
    # 公式：(总任务数 + 班级人数 - 1) // 班级人数，即向上取整
    blocks_per_grid = (N + (threads_per_block - 1)) // threads_per_block

    print(f"启动 Kernel: Grid={blocks_per_grid}, Block={threads_per_block}")
    print(f"总共派出线程数: {blocks_per_grid * threads_per_block}")

    # --- 5. 发射 Kernel (Launch) ---
    # 预热一次 (第一次运行会有编译开销，不计入时间)
    vector_add_kernel[blocks_per_grid, threads_per_block](d_A, d_B, d_C)
    cuda.synchronize() # 强制等待 GPU 跑完

    # 正式计时
    start = time.time()
    vector_add_kernel[blocks_per_grid, threads_per_block](d_A, d_B, d_C)
    cuda.synchronize() # 必须同步！因为 Kernel Launch 是异步的
    cost = time.time() - start

    print(f"GPU 计算耗时: {cost * 1000:.4f} ms")

    # --- 6. 搬回结果 (Device -> Host) ---
    C = d_C.copy_to_host()

    # 验证结果
    if np.allclose(C, A + B):
        print("验证成功！结果正确 (2 + 3 = 5)")
    else:
        print("验证失败！")

import time
if __name__ == "__main__":
    main()
