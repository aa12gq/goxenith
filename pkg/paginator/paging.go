package paginator

func GetPageOffset(pageNum, pageSize uint32) uint32 {
	return (pageNum - 1) * pageSize
}
