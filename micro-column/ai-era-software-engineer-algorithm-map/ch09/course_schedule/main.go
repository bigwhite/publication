package main

import "fmt"

func canFinish(numCourses int, prerequisites [][]int) bool {
	// 1. 构建图和入度表
	graph := make(map[int][]int) // key: 前置课, val: 后续课列表
	inDegree := make([]int, numCourses)

	for _, rel := range prerequisites {
		cur, pre := rel[0], rel[1] // 想修 cur，先修 pre
		graph[pre] = append(graph[pre], cur)
		inDegree[cur]++
	}

	// 2. 初始化队列：将所有入度为 0 的课（没门槛的课）入队
	queue := []int{}
	for i := 0; i < numCourses; i++ {
		if inDegree[i] == 0 {
			queue = append(queue, i)
		}
	}

	// 3. BFS 拓扑排序
	count := 0
	for len(queue) > 0 {
		course := queue[0]
		queue = queue[1:]
		count++

		// 解锁后续课程
		for _, nextCourse := range graph[course] {
			inDegree[nextCourse]--
			if inDegree[nextCourse] == 0 {
				queue = append(queue, nextCourse)
			}
		}
	}

	// 如果能修完所有课，说明没有环
	return count == numCourses
}

func main() {
	// 2 -> 1 -> 0
	reqs := [][]int{{1, 0}, {0, 1}}                             // 环: 1依赖0, 0依赖1
	fmt.Printf("Can finish (Circle): %v\n", canFinish(2, reqs)) // false

	reqs2 := [][]int{{1, 0}}                                     // 1依赖0
	fmt.Printf("Can finish (Linear): %v\n", canFinish(2, reqs2)) // true
}
