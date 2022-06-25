package util

// MinI64 ...
func MinI64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// Min ...
//
// Author : fuhaixu@ke.com<付海旭>
//
// Date : 下午9:02 2021/4/11
func Max(nums ...int64) int64 {
	if len(nums) == 0 {
		return 0
	}

	if len(nums) == 1 {
		return nums[0]
	}

	max := nums[0]

	for i := 0; i < len(nums); i++ {
		if max < nums[i] {
			max = nums[i]
		}
	}

	return max
}

// Min ...
//
// Author : fuhaixu@ke.com<付海旭>
//
// Date : 下午9:07 2021/4/11
func Min(nums ...int64) int64 {
	if len(nums) == 0 {
		return 0
	}

	if len(nums) == 1 {
		return nums[0]
	}

	min := nums[0]

	for i := 0; i < len(nums); i++ {
		if min > nums[i] {
			min = nums[i]
		}
	}

	return min
}

// CalculateI64Interval Calculate two intervals have intersection
func CalculateI64Interval(b1, e1, b2, e2 int64) int64 {
	// Each interval must greater than 0
	if e1 <= b1 || e2 <= b2 {
		return 0
	}

	// Keep smaller at front
	if b1 > b2 {
		tmp := b2
		b2 = b1
		b1 = tmp

		tmp = e2
		e2 = e1
		e1 = tmp
	}

	if b2 > b1 && b2 < e1 {
		return MinI64(e1, e2) - b2
	}

	return 0
}
