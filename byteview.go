package go_cache

// 存当前的缓存值
type ByteView struct {
	b []byte
}

// 返回当前的所占内存大小
func (v ByteView) Len() int {
	return len(v.b)
}

// b是只读的，使用 ByteSlice()方法返回一个拷贝,防止缓存值被外部程序修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

// 返回一个b的拷贝
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
