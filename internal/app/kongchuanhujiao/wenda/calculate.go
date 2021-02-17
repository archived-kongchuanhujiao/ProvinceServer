package wenda

import "coding.net/kongchuanhujiao/server/internal/app/datahub/public/wenda"

// CalculateQuestion 计算问答结果
func CalculateQuestion(w *wenda.Detail) (calc *wenda.CalculationsTab) {

	var (
		rightStus []uint64
		wrongStus []wenda.CalculationsWrong
	)

	rightStus = []uint64{}
	wrongStus = []wenda.CalculationsWrong{}

	correctAnswer := w.Questions.Key

	for _, ans := range w.Answers {
		if ans.Answer == correctAnswer {
			rightStus = append(rightStus, ans.QQ)
		} else {
			for i, option := range w.Questions.Options {
				if option != ans.Answer {
					continue
				}

				op := wrongStus[i]
				if &op == nil {
					op = wenda.CalculationsWrong{
						Type:   option,
						Member: []uint64{ans.QQ},
					}
				}

				op.Member = append(op.Member, ans.QQ)
				wrongStus[i] = op
			}
		}
	}

	calc = &wenda.CalculationsTab{
		Question: w.Questions.ID,
		Count:    uint8(len(w.Answers)),
		Right:    rightStus,
		Wrong:    wrongStus,
	}

	return
}
