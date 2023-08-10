package xcopy

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/jinzhu/copier"
)

// Copy 复制简单的数据类型和结构体， 注意含有枚举字段的结构体不要使用，需手动处理
func Copy(src, dst interface{}) error {
	return copyByJson(dst, src)
}

func copierCopy(src, dst interface{}) error {
	return copier.CopyWithOption(dst, src, copier.Option{DeepCopy: true, IgnoreEmpty: true})
}

func copyByJson(dst interface{}, src interface{}) error {
	if dst == nil {
		return fmt.Errorf("dst cannot be nil")
	}
	if src == nil {
		return fmt.Errorf("src cannot be nil")
	}
	bytes, err := sonic.Marshal(src)
	if err != nil {
		return fmt.Errorf("unable to marshal src: %s", err)
	}
	err = sonic.Unmarshal(bytes, dst)
	if err != nil {
		return fmt.Errorf("unable to unmarshal into dst: %s", err)
	}
	return nil
}
