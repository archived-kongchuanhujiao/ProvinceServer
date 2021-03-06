package wenda

import public "github.com/kongchuanhujiao/server/internal/app/datahub/public/wenda"

// CalculateResult 计算结果。
// 参数 i：问题标识号
func CalculateResult(i uint32) (r *public.Result, err error) {

	a, err := SelectAnswers(i)
	if err != nil {
		return
	}

	r = &public.Result{Count: uint8(len(a)), Right: []uint64{}, Wrong: []public.ResultWrongField{}}
	wrong := map[string][]uint64{}

	for _, i := range a {
		if m := i.Mark; m != "" {
			wrong[m] = append(wrong[m], i.QQ)
			continue
		}
		r.Right = append(r.Right, i.QQ)
	}

	for k, v := range wrong {
		r.Wrong = append(r.Wrong, public.ResultWrongField{
			Type:  k,
			Value: v,
		})
	}

	return
}
