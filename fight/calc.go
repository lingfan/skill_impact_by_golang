package fight

import "math"

//属性值计算公式
func CALCATTR1(A float64, B float64, C float64, D float64, E float64, F float64, G float64) float64 {
	return (A+B*C*0.01)*(1+D*0.01+E*0.01) + F + G
}
func CALCATTR2(A float64, B float64, C float64) float64 {
	return A + B + C
}
func CALCATTR3(A float64, B float64, C float64) float64 {
	return (1 - (1-A*0.01)*(1-B*0.01)*(1-C*0.01)) * 100
}
func CALCDAMAGE(A float64, B float64, C float64) float64 {
	return (A - B) * (1 - C*0.01)
}
func CHKMIN(A float64, B float64) float64 {
	A = math.Max(A, B)
	//A=(A<B)?B:A;
	return A
}
func GETSKILLGROUP(A float64) float64 {
	return A / 100
}
func GETSKILLLEVEL(A int) int {
	return A % 100
}
